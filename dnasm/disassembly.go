package dnasm

import (
	"fmt"
	"os"
	"sort"

	"github.com/zorchenhimer/emu-6502/mappers"
	"github.com/zorchenhimer/emu-6502"
)

//type OpAnalize func(d *Dissassembly, 

type Disassembly struct {
	m mappers.Mapper
	core *emu.Core

	chunks []*Chunk

	tokens map[uint32]emu.Token
	branches map[uint32]interface{} // element added if branch has been seen before
	jsrs map[uint32]bool // true if JSR has been returned from in a standard way

	labels map[uint32]*LabelMeta
	ramLabels map[uint16]*RamLabelMeta

	processed int
}

func New(rom []byte) (*Disassembly, error) {
	mapper, err := mappers.LoadFromBytes(rom)
	if err != nil {
		return nil, err
	}
	fmt.Println("Found mapper:", mapper.Name())

	core, err := emu.NewCore(mapper)
	core.Debug = true
	if err != nil {
		return nil, err
	}

	d := &Disassembly{
		m: mapper,
		core: core,

		chunks: []*Chunk{},
		tokens: make(map[uint32]emu.Token),
		labels: make(map[uint32]*LabelMeta),
		ramLabels: make(map[uint16]*RamLabelMeta),
		branches: make(map[uint32]interface{}),
		jsrs: make(map[uint32]bool),
	}

	return d, nil
}

func (d *Disassembly) AddVector(address uint16) {
	offset := d.m.Offset(address)
	d.labels[offset] = &LabelMeta{
		Address: address,
		Offset: offset,
		Type: LT_Vector}

	c := &Chunk{
		MapperState: d.m.GetState(),
		CpuState:    d.core.GetState(),
	}

	c.CpuState.PC = address
	d.chunks = append(d.chunks, c)
}

type offsetSort []uint32

func (o offsetSort) Len() int           { return len(o) }
func (o offsetSort) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o offsetSort) Less(i, j int) bool { return o[i] < o[j] }

// WriteToFile writes the full dissassembly to a single source
// file with the given filename.
func (d *Disassembly) WriteToFile(filename string) error {

	offsets := []uint32{}
	for offset, _ := range d.tokens {
		offsets = append(offsets, offset)
	}

	s := offsetSort(offsets)
	sort.Sort(s)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//for offset, token := range d.tokens {
	for _, offset := range s {
		if label, exists := d.labels[offset]; exists {
			fmt.Fprintf(file, "\n%s:\n", label.Label())
		}
		fmt.Fprintf(file, "[0x%08X] %s\n", offset, d.tokens[offset].String())
	}
	//return fmt.Errorf("not implemented")
	return nil
}

func (d *Disassembly) process() error {
	if len(d.chunks) == 0 {
		return fmt.Errorf("No chunks to process!")
	}

	//fmt.Printf("Starting at address %04X\n", address)
	d.core.HardReset()
	//d.core.PC = address

	for len(d.chunks) != 0 {
		fmt.Printf("chunk length: %d\n", len(d.chunks))

		chunk := d.chunks[0]
		d.chunks = d.chunks[1:]
		err := d.processChunk(chunk)
		if err != nil {
			return err
		}
	}

	fmt.Printf("No more chunks to process. Processed %d chunks.\n", d.processed)
	return nil
}

func (d *Disassembly) processChunk(c *Chunk) error {
	defer func() { d.processed++ }()

	d.core.SetState(c.CpuState)
	d.m.SetState(c.MapperState)

	fmt.Printf("\nProcessing chunk starting at 0x%08X\n  FromJsr: %t\n  Address: %04X\n",
		d.core.PC, c.FromJsr, c.Address)

	for i := 0; i < 10000; i++ {
		t, err := d.core.Peek()
		if err != nil {
			return err
		}

		//f, ok := opFunctions[t.OpCode]
		//if ok {
		//	doTick, err := f(t)
		//}

		//if doTick {
		//	d.core.Tick()
		//}

		//d.tokens[off] = t
		//fmt.Println(">", d.core.LastInstruction())


		////
		if t.Type() == emu.TT_Data {
			return fmt.Errorf("Found data instead of instruction")
		}

		if t.Type() == emu.TT_Unknown {
			return fmt.Errorf("Unknown token type")
		}

		ii := t.(emu.InstructionAny)
		off := d.m.Offset(d.core.PC)

		if d.core.PC < 0x8000 {
			fmt.Println("PC in RAM; aborting")
			return nil
		}

		tok, exists := d.tokens[off]
		if exists {
			fmt.Println("[token exists]")
			ta, ok := tok.(emu.InstructionAny)
			if !ok {
				return fmt.Errorf("Unable to cast to InstructionAny")
			}

			if isBranch(ta.OpCode()) {
				if _, hit := d.branches[off]; hit {
					fmt.Printf("Back in known code; offset: %08X; address: %04X\n", off, d.core.PC)
					return nil
				}
				fmt.Println("[branch not hit]")
				d.branches[off] = nil
			} else {
				fmt.Println("[not a branch]")
				fmt.Printf("Back in known code; offset: %08X; address: %04X\n", off, d.core.PC)
				return nil
			}
		}

		if ii.OpCode() == emu.OP_JSR {
			iw, ok := t.(*emu.InstructionWord)
			if !ok {
				// Shouldn't ever happen
				fmt.Printf("%v\n", t)
				panic("JSR isn't a word instruction?")
			}

			nc := d.newChildChunk(c)
			nc.CpuState.PC = iw.Arg() // Start new chunk after the jump
			nc.FromJsr = true
			nc.Address = d.core.PC+2
			d.chunks = append(d.chunks, nc)

			d.labels[d.m.Offset(iw.Arg())] = &LabelMeta{
				Address: iw.Arg(),
				Offset: d.m.Offset(iw.Arg()),
				Type: LT_Jsr,
			}

			// Skip jump
			d.core.PC += 3

			fmt.Printf("[JSR] routine at %08X:%04X\n", d.m.Offset(nc.CpuState.PC), nc.CpuState.PC)
			d.tokens[off] = t
			continue
		}

		if ii.OpCode() == emu.OP_JMP_ID {
			iw, ok := t.(*emu.InstructionWord)
			if !ok {
				// Shouldn't ever happen
				fmt.Printf("%v\n", t)
				panic("JMP (Implied) isn't a word instruction?")
			}

			rlm := &RamLabelMeta{
				Address: iw.Arg(),
				Type: RLT_Pointer,
				Used: make(map[uint32]interface{}),
			}

			rlm.Used[d.m.Offset(d.core.PC)] = nil
			d.ramLabels[iw.Arg()] = rlm
		}

		if ii.OpCode() == emu.OP_RTS {
			fmt.Println("[rts]")
			if !c.FromJsr {
				return fmt.Errorf("RTS outside of routine")
			}

			ret := d.core.PeekStackWord()
			if ret != c.Address {
				fmt.Println("RTS to somewhere unknown; continuing.")
			} else {
				fmt.Println("RTS to calling address; routine finished.")
				// Chunk done
				return nil
			}
		}

		if ii.OpCode() == emu.OP_RTI {
			fmt.Println("Found RTI; aborting")
			return nil
		}

		if  _, hit := d.branches[off]; !hit && isBranch(ii.OpCode()) {
			fmt.Println("[not hit && isBranch]")
			nc := d.newChildChunk(c)

			fmt.Printf("State %s -> ", flagsToString(nc.CpuState.Phlags))
			switch ii.OpCode() {
			case emu.OP_BCC, emu.OP_BCS:
				nc.CpuState.Phlags ^= emu.FLAG_CARRY
			case emu.OP_BEQ, emu.OP_BNE:
				nc.CpuState.Phlags ^= emu.FLAG_ZERO
			case emu.OP_BMI, emu.OP_BPL:
				nc.CpuState.Phlags ^= emu.FLAG_NEGATIVE
			case emu.OP_BVC, emu.OP_BVS:
				nc.CpuState.Phlags ^= emu.FLAG_OVERFLOW
			}
			fmt.Printf("%s\n", flagsToString(nc.CpuState.Phlags))

			d.chunks = append(d.chunks, nc)

			ib, ok := t.(*emu.InstructionBranch)
			if !ok {
				panic(fmt.Sprintf("Branch is not InstructionBranch?\n%T", t))
			}

			d.labels[d.m.Offset(ib.Destination())] = &LabelMeta{
				Address: ib.Destination(),
				Offset: d.m.Offset(ib.Destination()),
				Type: LT_Branch,
			}
		}

		d.tokens[off] = t

		d.core.Tick()
		fmt.Println(">", d.core.LastInstruction())
	}

	fmt.Println("outside of processChunk() for loop")
	return nil
}

func (d *Disassembly) newChildChunk(c *Chunk) *Chunk {
	nc := &Chunk{
		MapperState: d.m.GetState(),
		CpuState:    d.core.GetState(),
	}

	if c.FromJsr {
		nc.FromJsr = true
		nc.Address = d.core.PC+2 // Address expected to return to
	}

	return nc
}

func flagsToString(ph uint8) string {
	sc := "-"
	sz := "-"
	si := "-"
	sd := "-"
	sv := "-"
	sn := "-"

	if ph&emu.FLAG_CARRY != 0 {
		sc = "C"
	}

	if ph&emu.FLAG_ZERO != 0 {
		sz = "Z"
	}

	if ph&emu.FLAG_INTERRUPT != 0 {
		si = "I"
	}

	if ph&emu.FLAG_DECIMAL != 0 {
		sd = "D"
	}

	if ph&emu.FLAG_OVERFLOW != 0 {
		sv = "V"
	}

	if ph&emu.FLAG_NEGATIVE != 0 {
		sn = "N"
	}

	return fmt.Sprintf("%s%s--%s%s%s%s", sn, sv, sd, si, sz, sc)
}

func isBranch(op uint8) bool {
	switch op {
	case emu.OP_BCC, emu.OP_BCS, emu.OP_BEQ, emu.OP_BMI,
	     emu.OP_BNE, emu.OP_BPL, emu.OP_BVC, emu.OP_BVS:
		return true
	}
	return false
}
