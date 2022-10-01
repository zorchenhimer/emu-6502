package emu

import (
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sort"
	"time"

	"github.com/zorchenhimer/emu-6502/mappers"
)

const HistoryLength int = 100

const (
	VECTOR_NMI   uint16 = 0xFFFA
	VECTOR_RESET uint16 = 0xFFFC
	VECTOR_IRQ   uint16 = 0xFFFE
)

const (
	NTSC time.Duration = time.Nanosecond * 16666667 // close enough, lol
	PAL  time.Duration = time.Millisecond * 20
)

type Core struct {
	// Main registers
	A uint8
	X uint8
	Y uint8

	// Other registers
	PC     uint16 // Program counter
	Phlags uint8  // Status flags
	SP     uint8  // Stack pointer

	NmiFrequency time.Duration

	memory mappers.Mapper

	InstructionLimit uint64 // number of instructions to run
	testing          bool
	testDone         bool
	ticks            uint64

	lastPC       uint16
	lastSame     int
	lastReadAddr uint16
	checkStuck   bool

	// VERY verbose output
	Debug     bool
	DebugFile io.Writer

	history    [HistoryLength]string
	historyIdx int

	nmiTicker *time.Ticker
	nmiCount  uint

	Breakpoints *Breakpoints

	stop bool // set to true to end the Run loop

	// used for RunRoutine()
	runRoutine   bool
	routineDepth int

	EnableCDL bool
	//cdl *cdlData

	visited map[uint32]string
	destinations []uint32
	cdl map[uint32]byte
}

const (
	cdl_Unknown      byte = 0x00
	cdl_Code         byte = 0x01
	cdl_Data         byte = 0x02
	cdl_Label        byte = 0x10  // "Jump Target"
	cdl_IndirectData byte = 0x20
	cdl_PcmData      byte = 0x40
	cdl_SubEntry     byte = 0x80
)

func NewCore(rom mappers.Mapper) (*Core, error) {
	c := &Core{
		A:      0,
		X:      0,
		Y:      0,
		PC:     0,
		Phlags: 0,
		SP:     0,

		memory: rom,

		//InstructionLimit: instrLimit,

		history:   [HistoryLength]string{},
		//nmiTicker: time.NewTicker(nmiFrequency),
		Breakpoints: &Breakpoints{},

		visited:      make(map[uint32]string),
		destinations: []uint32{},
		cdl:          make(map[uint32]byte),
	}

	c.PC = c.ReadWord(VECTOR_RESET)
	return c, nil
}

// Read address.  This will read from API registers if needed.
func (c *Core) ReadByte(addr uint16) uint8 {
	c.lastReadAddr = addr
	val := c.memory.ReadByte(addr)
	c.Breakpoints.Read(c, addr, val)
	return val
}

func (c *Core) ReadWord(addr uint16) uint16 {
	defer func() { c.lastReadAddr = addr }() // will this fire off correctly? idk
	return uint16(c.ReadByte(addr)) | (uint16(c.ReadByte(addr+1)) << 8)
}

// Write to an address.  This will delegate to API if needed.
func (c *Core) WriteByte(addr uint16, value byte) {
	c.Breakpoints.Write(c, addr, value)
	c.memory.WriteByte(addr, value)
}

func (c *Core) Run() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			c.stop = true
		}
	}()

	if c.nmiTicker != nil {
		defer c.nmiTicker.Stop()
	}

	if c.DebugFile != nil {
		c.Debug = true
	}

	start := time.Now()
	defer func() { fmt.Printf("time: %s\n", time.Now().Sub(start)) }()

	limit := false
	if c.InstructionLimit > 0 {
		//fmt.Printf("Setting instruction limit to %d\n", c.InstructionLimit)
		limit = true
	}

	c.initCdl()

	done := false
	var err error
	for !(done || c.stop) {
		err = c.tick()
		if err != nil {
			return err
		}

		if limit {
			c.InstructionLimit -= 1
			if c.InstructionLimit <= 0 {
				if c.testing {
					return fmt.Errorf("Instruction limit hit")
				}
				done = true
			}
		}

		if c.testing {
			done = c.testDone
		}
	}

	if c.stop {
		c.dumpHistory()
		return fmt.Errorf("Halt received")
	}

	fmt.Printf("nmiCount: %d\n", c.nmiCount)

	return nil
}

func (c *Core) initCdl() {
	if !c.EnableCDL {
		return
	}

	if cbm, ok := c.memory.(mappers.CallbackMapper); ok {
		f := func(address uint16, data uint8){
			offset, _ := c.memory.Offset(address)
			c.cdl[offset] |= cdl_Data
		}

		cbm.CallbackRead(f)
		cbm.CallbackWrite(f)
	}
}

// Run a routine and return after the last RTS
func (c *Core) RunRoutine(address uint16) error {
	if c.DebugFile != nil {
		c.Debug = true
	}

	c.pushAddress(address)
	c.routineDepth = 0
	c.runRoutine = true
	c.PC = address

	limit := false
	if c.InstructionLimit > 0 {
		limit = true
	}

	c.initCdl()

	var err error
	for c.routineDepth > -1 && !c.stop {
		err = c.tick()
		if err != nil {
			return err
		}

		if limit {
			c.InstructionLimit -= 1
			if c.InstructionLimit <= 0 {
				if c.testing {
					return fmt.Errorf("Instruction limit hit")
				}
				c.stop = true
			}
		}
	}

	return nil
}

func (c *Core) dumpHistory() {
	if !c.Debug {
		return
	}

	for i := c.historyIdx; i < HistoryLength; i++ {
		if c.history[i] == "" {
			return
		}
		fmt.Println(c.history[i])
	}

	for i := 0; i < c.historyIdx; i++ {
		if c.history[i] == "" {
			return
		}
		fmt.Println(c.history[i])
	}

}

func (c *Core) Halt() {
	c.stop = true
	//fmt.Println("CPU Halt()'d")
}

func (c *Core) HardReset() {
	c.memory.ClearRam()
	c.A = 0
	c.X = 0
	c.Y = 0
	c.PC = 0
	c.Phlags = 0
	c.SP = 0
}

func (c *Core) Reset() {
	c.runInterrupt(VECTOR_RESET)
}

func (c *Core) IRQ() {
	c.runInterrupt(VECTOR_IRQ)
}

func (c *Core) NMI() {
	c.runInterrupt(VECTOR_IRQ)
}

func (c *Core) runInterrupt(interrupt uint16) {
	if vector, ok := interruptList[interrupt]; ok {
		vector.Execute(c)
	}
}

func (c *Core) tick() error {
	//c.PC += 1
	if c.checkStuck {
		if c.PC == c.lastPC {
			c.lastSame++
		} else {
			c.lastSame = 0
			c.lastPC = c.PC
		}

		if c.lastSame > 0 {
			c.dumpHistory()
			return fmt.Errorf("Stuck at $%04X", c.PC)
		}
	}

	if c.nmiTicker != nil {
		// If it's time to NMI, do it.
		// Note that this can never happen during the execution of another
		// instruction in this implementation.  That isn't the case for
		// real hardware.
		select {
		case <-c.nmiTicker.C:
			c.nmiCount++
			c.NMI()
		default:
		}
	}

	c.Breakpoints.Execute(c, c.PC, 0)

	opcode := c.ReadByte(c.PC)

	if opcode == 0xFF && c.testing {
		c.testDone = true
		return nil // 0xFF means end of test
	}

	instr, ok := instructionList[opcode]
	if !ok || instr == nil {
		c.dumpHistory()
		return fmt.Errorf("OP Code not implemented: [$%04X] $%02X", c.PC, opcode)
	}

	oppc := c.PC

	c.ticks++
	instr.Execute(c)

	if c.EnableCDL && c.PC >= 0x8000 {
		l := uint16(instr.InstrLength(c))
		bin := []string{}
		for i := uint16(0); i < l; i++ {
			bin = append(bin, fmt.Sprintf("%02X", c.memory.ReadByte(oppc+i)))
			offset, _ := c.memory.Offset(oppc+i)
			c.cdl[offset] |= cdl_Code
		}

		// If we've already visited this code, don't regenerate the assembly.
		offset, _ := c.memory.Offset(oppc)
		if _, ok := c.visited[offset]; !ok {
			c.visited[offset] = fmt.Sprintf("[$%06X:$%04X] %s", offset, oppc, instr.Name() +" "+ instr.AddressMeta().CleanAsm(c, oppc) +" ; "+ strings.Join(bin, " "))
		}
	}

	if c.Debug {
		dbgLine := c.HistoryString(oppc, instr)

		c.history[c.historyIdx] = dbgLine
		c.historyIdx += 1
		if c.historyIdx >= HistoryLength {
			c.historyIdx = 0
		}

		if c.DebugFile != nil {
			fmt.Fprintln(c.DebugFile, dbgLine)
		}
	}

	return nil
}

type branch struct {
	address uint16
	state any
}

func (c *Core) Disassemble(start uint16) error {
	branches := []branch{}	// list of branches in rom space
	c.initCdl()

	c.PC = start
	nextBranch := false
	limit := 10000

	Main:
	for {
		limit--
		if limit < 0 {
			fmt.Println("limit reached")
			return nil
		}
		offset, _ := c.memory.Offset(c.PC)
		if _, ok := c.visited[offset]; ok || nextBranch {
			if ok {
				fmt.Printf("already visited $%06X\n", offset)
			}
			if len(branches) == 0 {
				fmt.Println("no more branches")
				break Main
			}

			noff, _ := c.memory.Offset(branches[0].address)
			fmt.Printf("going to next branch at $%04X ($%06X)\n", branches[0].address, noff)

			nextBranch = false
			c.PC = branches[0].address
			c.memory.SetState(branches[0].state)
			if len(branches) > 1 {
				branches = branches[1:]
			} else {
				branches = []branch{}
			}
		}

		opcode := c.ReadByte(c.PC)
		if opcode == 0 {
			fmt.Printf("zero opcode @ $%04X ($%04X)\n", c.PC+0x10 -0x8000, c.PC)
		}
		instr, ok := instructionList[opcode]
		if !ok || instr == nil {
			nextBranch = true
			continue
		}

		c.cdl[offset] |= cdl_Code
		operandAddr, size := instr.AddressMeta().Address(c)

		var dataAddr uint16
		if size == 2 {
			dataAddr = uint16(c.memory.ReadByte(operandAddr))
			c.cdl[offset+1] |= cdl_Code
		} else if size == 3 {
			c.cdl[offset+1] |= cdl_Code
			c.cdl[offset+2] |= cdl_Code

			dataAddr = c.memory.ReadWord(operandAddr)
		}

		//fmt.Printf("%q\n", instr.Name())

		// If we've already visited this code, don't regenerate the assembly.
		if _, ok := c.visited[offset]; !ok {
			str := fmt.Sprintf("[$%06X:$%04X] %s", offset, c.PC, instr.Name() +" "+ instr.AddressMeta().CleanAsm(c, c.PC))
			c.visited[offset] = str
			fmt.Println(str)
		}

		switch instr.Name() {
		// standard instructions
		case "ADC", "AND", "ASL", "BIT", "CLC", "CLD",
			 "CLI", "CLV", "CMP", "CPX", "CPY", "DEC",
			 "DEX", "DEY", "EOR", "INC", "INX", "INY",
			 "LDA", "LDX", "LDY", "LSR", "NOP", "ORA",
			 "PHA", "PHP", "PLA", "PLP", "ROL", "ROR",
			 "SBC", "SEC", "SED", "SEI", "STA", "STX",
			 "STY", "TAX", "TAY", "TSX", "TXA", "TXS",
			 "TYA":

			 if size > 1 {
				 dataOffset, isRom := c.memory.Offset(dataAddr)
				 if isRom {
					 c.cdl[dataOffset] |= cdl_Data
				 }
			 }

		// branches
		case "BCC", "BCS", "BEQ", "BMI", "BNE", "BPL", "BVC", "BVS", "JSR":
			//sint := int8(c.memory.ReadByte(operandAddr))
			//baddr := uint16(int32(c.PC) + int32(sint))
			//baddr := c.addrRelative(c.PC, c.memory.ReadByte(operandAddr))
			////baddr := c.PC + sint
			boff, _ := c.memory.Offset(operandAddr)
			if instr.Name() == "JSR" {
				fmt.Printf("adding JSR at $%06X:$%04X\n", boff, operandAddr)
			} else {
				fmt.Printf("adding branch at $%06X:$%04X\n", boff, operandAddr)
			}
			branches = append(branches, branch{
				address: operandAddr,
				state: c.memory.GetState(),
			})

		case "JMP":
			if instr.AddressMeta() == ADDR_Absolute {
				c.PC = operandAddr
				continue
			} else {
				// TODO: relative jump
				return fmt.Errorf("relative jump not implemented")
			}

		//case "JSR":
		//	//jsrAddr := c.ReadWord(operandAddr)
		//	joff, _ := c.memory.Offset(operandAddr)
		//	fmt.Printf("adding JSR at $%06X:$%04X\n", joff, operandAddr)
		//	branches = append(branches, branch{
		//		address: operandAddr,
		//		state: c.memory.GetState(),
		//	})

		case "BRK":
			fmt.Printf("BRK: %02X\n", opcode)
			nextBranch = true

		case "RTI", "RTS":
			// BRK shouldn't really do this, but w/e
			nextBranch = true

		default:
			return fmt.Errorf("unknown instruction: %q", instr.Name())
		}

		if instr.Name() != "JMP" {
			c.PC += uint16(size)
		}
	}

	return nil
}

func (c *Core) WriteCdl(writer io.Writer) error {
	// Find the max address, make a slice that size, fill in all the data
	// in that slice, then write said slice to the io.Writer.
	//
	// This will sort the data correctly and pad it with 0x00 where we do not
	// have any data.

	max := uint32(0)
	for addr, _ := range c.cdl {
		if addr > max {
			max = addr
		}
	}

	vals := make([]byte, max+1)

	for addr, val := range c.cdl {
		vals[int(addr)] = val
	}

	_, err := writer.Write(vals)
	return err
}

func (c *Core) WriteVisited(writer io.Writer) error {
	addrs := []int{}

	for addr, _ := range c.visited {
		addrs = append(addrs, int(addr))
	}

	sort.Ints(addrs)

	for _, addr := range addrs {
		//_, err := fmt.Fprintln(writer, "%s ; %04X\n", c.visited[uint32(addr)], addr)
		line := strings.TrimSpace(c.visited[uint32(addr)])
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}

		parts := strings.Split(line, " ")
		if parts[1] == "RTS" || parts[1] == "RTI" || parts[1] == "BRK" || strings.HasPrefix(parts[1], "JMP") {
			fmt.Fprintln(writer, "")
		}
	}

	return nil
}

func (c *Core) HistoryString(oppc uint16, instr Instruction) string {
	l := instr.InstrLength(c)
	ops := []string{}
	for i := uint8(0); i < l; i++ {
		ops = append(ops, fmt.Sprintf("%02X", c.ReadByte(oppc+uint16(i))))
	}

	return fmt.Sprintf("[%06d] $%04X: %-9s %s %-17s %s %s",
		c.ticks,
		oppc,
		strings.Join(ops, " "),
		instr.Name(),
		instr.AddressMeta().Asm(c, oppc), // oppc == OP code PC
		c.Registers(),
		c.stackString(),
	)
}

func (c *Core) Instructions() []string {
	ret := []string{}
	for _, instr := range instructionList {
		var op byte
		switch instr.(type) {
		case StandardInstruction:
			si := instr.(StandardInstruction)
			op = si.OpCode
		case Branch:
			br := instr.(Branch)
			op = br.OpCode
		case Jump:
			j := instr.(Jump)
			op = j.OpCode
		case ReadModifyWrite:
			rmw := instr.(ReadModifyWrite)
			op = rmw.OpCode
		}
		ret = append(ret, fmt.Sprintf("$%02X %s %s", op, instr.Name(), instr.AddressMeta().Name))
	}
	return ret
}

func (c *Core) stackString() string {
	st := []string{}
	length := 0xFF - c.SP
	if length == 0 {
		return ""
	}

	for i := length; i > 0; i-- {
		st = append(st, fmt.Sprintf("$%02X", c.ReadByte(uint16(c.SP+i)|0x0100)))
	}

	return strings.Join(st, " ")
}

func (c *Core) DumpMemoryRange(filename string, start, end uint16) error {
	if end < start {
		return fmt.Errorf("Invalid dump range given")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "start: $%02X end: $%02X\n", start, end)

	vals := []byte{}
	current := start

	for current <= end {
		vals = append(vals, c.ReadByte(current))
		current++
	}

	for i, b := range vals {
		fmt.Fprintf(file, "$%02X: $%02X (%d)\n", i+int(start), b, b)
	}

	return nil
}

const (
	FLAG_CARRY     uint8 = 0x01
	FLAG_ZERO      uint8 = 0x02
	FLAG_INTERRUPT uint8 = 0x04
	FLAG_DECIMAL   uint8 = 0x08

	FLAG_BREAK    uint8 = 0x30
	FLAG_IRQ      uint8 = 0x20
	FLAG_OVERFLOW uint8 = 0x40
	FLAG_NEGATIVE uint8 = 0x80
)

func flagToString(ph uint8) string {
	switch ph {
	case FLAG_CARRY:
		return "FLAG_CARRY"
	case FLAG_ZERO:
		return "FLAG_ZERO"
	case FLAG_INTERRUPT:
		return "FLAG_INTERRUPT"
	case FLAG_DECIMAL:
		return "FLAG_DECIMAL"
	case FLAG_OVERFLOW:
		return "FLAG_OVERFLOW"
	case FLAG_NEGATIVE:
		return "FLAG_NEGATIVE"
	}
	return "FLAG_UNUSED"
}

func flagsToString(ph uint8) string {
	sc := "-"
	sz := "-"
	si := "-"
	sd := "-"
	sv := "-"
	sn := "-"

	if ph&FLAG_CARRY != 0 {
		sc = "C"
	}

	if ph&FLAG_ZERO != 0 {
		sz = "Z"
	}

	if ph&FLAG_INTERRUPT != 0 {
		si = "I"
	}

	if ph&FLAG_DECIMAL != 0 {
		sd = "D"
	}

	if ph&FLAG_OVERFLOW != 0 {
		sv = "V"
	}

	if ph&FLAG_NEGATIVE != 0 {
		sn = "N"
	}

	return fmt.Sprintf("%s%s--%s%s%s%s", sn, sv, sd, si, sz, sc)
}

func (c *Core) Registers() string {
	return fmt.Sprintf("A: %02X (%-3d) X: %02X (%-3d) Y: %02X (%-3d) SP: %02X (%-3d) [%02X] %s",
		c.A,
		c.A,
		c.X,
		c.X,
		c.Y,
		c.Y,
		c.SP,
		c.SP,
		c.Phlags,
		flagsToString(c.Phlags),
	)
}

func (c *Core) DumpPage(page uint8) {
	vals := []string{}
	base := uint16(page) << 8
	for i := uint16(0); i < 256; i++ {
		vals = append(vals, fmt.Sprintf("%02X", c.ReadByte(base+i)))
	}

	for i := 0; i < 256; i += 16 {
		fmt.Printf("%04X: %s\n", int(base)+i, strings.Join(vals[i:i+16], " "))
	}
}

func (c *Core) DumpMemoryToFile(filename string) error {
	vals := []string{}
	for i := uint(0); i < 0x10000; i++ {
		vals = append(vals, fmt.Sprintf("%02X", c.ReadByte(uint16(i))))
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < 0x10000; i += 16 {
		fmt.Fprintf(file, "%04X: %s\n", i, strings.Join(vals[i:i+16], " "))
	}
	return nil
}

func (c *Core) Ticks() uint64 {
	return c.ticks
}

// Set zero and negative flags based on the given value
func (c *Core) setZeroNegative(value uint8) uint8 {
	// zero
	if value == 0 {
		c.Phlags = c.Phlags | FLAG_ZERO
	} else {
		c.Phlags = c.Phlags & (FLAG_ZERO ^ 0xFF)
	}

	// negative
	if value&0x80 != 0 {
		c.Phlags = c.Phlags | FLAG_NEGATIVE
	} else {
		c.Phlags = c.Phlags & (FLAG_NEGATIVE ^ 0xFF)
	}

	return value
}

// addrRelative works differently than all other addressing functions.
// It takes the value for offset, and uses the PC of the instruction
// as the start point.  Call this before incrementing PC.
func (c *Core) addrRelative(pc uint16, offset uint8) uint16 {
	addr := pc + 2
	val, negative := TwosCompInv(offset)

	if negative {
		addr -= uint16(val)
	} else {
		addr += uint16(val)
	}

	return addr
}

func (c *Core) twosCompAdd(a, b uint8) uint8 {
	carry := uint8(0)
	if (c.Phlags & FLAG_CARRY) == FLAG_CARRY {
		carry = 1
	}
	val := a + b + carry

	if (val < a) || ((val == a) && (carry != 0)) {
		// set carry
		c.Phlags = c.Phlags | FLAG_CARRY
	} else {
		// reset carry
		c.Phlags = c.Phlags & (FLAG_CARRY ^ 0xFF)
	}

	if ((a ^ val) & (b ^ val) & 0x80) == 0x80 {
		c.Phlags |= FLAG_OVERFLOW
	} else {
		c.Phlags &^= FLAG_OVERFLOW
	}

	c.setZeroNegative(val)
	return val
}

func (c *Core) twosCompSubtract(a, b uint8) uint8 {
	b = (b - 1) ^ 0xFF
	return c.twosCompAdd(a, b)
}

func (c *Core) pushAddress(addr uint16) {
	c.pushByte(uint8(addr >> 8))
	c.pushByte(uint8(addr & 0xFF))
}

func (c *Core) pullAddress() uint16 {
	return uint16(c.pullByte()) | uint16(c.pullByte())<<8
}

func (c *Core) pushByte(val uint8) {
	c.WriteByte(uint16(c.SP)|0x0100, val)
	c.SP -= 1
}

func (c *Core) pullByte() uint8 {
	c.SP += 1
	return c.ReadByte(uint16(c.SP) | 0x0100)
}
