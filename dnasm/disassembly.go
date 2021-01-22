package dnasm

import (
	"fmt"
	"os"
	"sort"

	"github.com/zorchenhimer/emu-6502/mappers"
	"github.com/zorchenhimer/emu-6502"
)

type Disassembly struct {
	m mappers.Mapper
	core *emu.Core

	chunks []*Chunk

	tokens map[uint32]emu.Token
	branches map[uint32]interface{} // element added if branch has been seen before
	jsrs map[uint32]bool // true if JSR has been returned from in a standard way

	processed int
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
		fmt.Fprintf(file, "[0x%08X] %s\n", offset, d.tokens[offset].String())
	}
	//return fmt.Errorf("not implemented")
	return nil
}

// Listing returns a slice containing all of the instructions in
// the disassembly.
//func (d *Disassembly) Listing() []Token {
//	return nil
//}

func (d *Disassembly) process(address uint16) error {
	fmt.Printf("Starting at address %04X\n", address)
	d.core.HardReset()
	d.core.PC = address
	d.branches = make(map[uint32]interface{})
	d.jsrs = make(map[uint32]bool)

	startChunk := &Chunk{
		//Address: address,
		MapperState: d.m.GetState(),
		CpuState: d.core.GetState(),
	}

	first := true
	
	//for i := 0; i < 100000; i++{
	for len(d.chunks) != 0 || first {
		if first {
			first = false
			d.processChunk(startChunk)
			continue
		}

		fmt.Printf("chunk length: %d\n", len(d.chunks))

		chunk := d.chunks[0]
		d.chunks = d.chunks[1:]
		err := d.processChunk(chunk)
		if err != nil {
			return err
		}
	}

	fmt.Println("No more chunks to process")
	return nil
}

func (d *Disassembly) processChunk(c *Chunk) error {
	defer func() { d.processed++ }()

	d.core.SetState(c.CpuState)
	d.m.SetState(c.MapperState)

	// local to current process
	// index is JSR instruction, value is 
	jstack := &JsrStack{}

	fmt.Printf("\nProcessing chunk starting at 0x%08X\n", d.core.PC)

	for i := 0; i < 10000; i++ {
		t, err := d.core.Peek()
		if err != nil {
			return err
		}

		if t.Type() == emu.TT_Data {
			return fmt.Errorf("Found data instead of instruction")
		}

		if t.Type() == emu.TT_Unknown {
			return fmt.Errorf("Unknown token type")
		}

		ii := t.(emu.InstructionAny)
		//instr := emu.InstructionList[ii.OpCode()]
		//raw := d.m.ReadByte(d.core.PC)
		off := d.m.Offset(d.core.PC)

		//fmt.Printf("[%08X:%04X] %s %s\n",
		//	off, d.core.PC,
		//	instr.Name(), instr.AddressMeta().Asm(d.core, d.core.PC))

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
				fmt.Printf("%v\n", t)
				panic("JSR isn't a word instruction?")
			}

			// followed routine before and it returned back to its calling JSR
			if ret, hit := d.jsrs[d.m.Offset(iw.Arg())]; hit && ret {
				// don't follow JSR
				d.core.PC += 3
				continue

			// haven't followed routine before
			} else if !hit {
				// push return address
				jstack.Push(&JsrElement{d.m.Offset(d.core.PC+2), d.m.Offset(iw.Arg())})
			}
		}

		if ii.OpCode() == emu.OP_RTS {
			fmt.Println("[rts]")
			val, err := jstack.Pop()
			if err != nil {
				fmt.Println("[pop error]")
				return err
			}

			ret := d.m.Offset(d.core.PeekStackWord())
			// Does the routine return to its original calling context?
			d.jsrs[val.RoutineAddress] = (ret == val.CallAddress)
			fmt.Printf("routine returns to calling: %t\n", (ret == val.CallAddress))
		}

		if ii.OpCode() == emu.OP_RTI {
			fmt.Println("Found RTI; aborting")
			return nil
		}

		if  _, hit := d.branches[off]; !hit && isBranch(ii.OpCode()) {
			fmt.Println("[not hit && isBranch]")
			newCpuState := d.core.GetState()
			newMapState := d.m.GetState()

			fmt.Printf("State %s -> ", flagsToString(newCpuState.Phlags))
			switch ii.OpCode() {
			case emu.OP_BCC, emu.OP_BCS:
				newCpuState.Phlags ^= emu.FLAG_CARRY
			case emu.OP_BEQ, emu.OP_BNE:
				newCpuState.Phlags ^= emu.FLAG_ZERO
			case emu.OP_BMI, emu.OP_BPL:
				newCpuState.Phlags = newCpuState.Phlags ^ emu.FLAG_NEGATIVE
			case emu.OP_BVC, emu.OP_BVS:
				newCpuState.Phlags ^= emu.FLAG_OVERFLOW
			}
			fmt.Printf("%s\n", flagsToString(newCpuState.Phlags))

			d.chunks = append(d.chunks, &Chunk{
				MapperState: newMapState,
				CpuState:    newCpuState,
			})
		}

		//fmt.Printf("%s | %s\n", t.String(), op.Name())//op.AddressMeta().Asm(d.core, d.core.PC))
		d.tokens[off] = t

		d.core.Tick()
		fmt.Println(">", d.core.LastInstruction())
	}

	fmt.Println("outside of processChunk() for loop")
	return nil
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
