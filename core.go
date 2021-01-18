package emu

import (
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"os/signal"
	"strings"
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
}

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

// Run a routine and return after the last RTS
func (c *Core) RunRoutine(address uint16) error {
	if c.DebugFile != nil {
		c.Debug = true
	}

	//start := time.Now()
	//defer func() {
	//	fmt.Printf("time: %s\nticks: %d\n",
	//		time.Now().Sub(start),
	//		c.ticks)
	//}()

	// Start value of stack pointer
	//sp := c.SP
	c.routineDepth = 0
	c.runRoutine = true
	c.PC = address

	var err error
	for c.routineDepth > -1 && !c.stop {
		err = c.tick()
		if err != nil {
			return err
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
	fmt.Println("CPU Halt()'d")
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
