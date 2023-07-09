package mappers

import (
	"io"
	//"bytes"
	"fmt"
	//"os"
	"strings"
	//"encoding/binary"

	sbox "github.com/zorchenhimer/go-nes/studybox"
)

/*
BRAM is a set of eight 4k banks (16k total) mapped at CPU $5xxx with
$4400-$4fff being a mirror of the top 3k of bank 0.

PRAM is a set of four 8k banks (32k total) mapped at CPU $6000-$7fff.
*/

func init() {
	registerMapper(186, NewStudyBox)
}

type StudyBox struct {
	rom []byte
	mainRam [0x0800]byte
	ramA [0x4000]byte
	ramB [0x8000]byte
	RamABank uint8
	RamBBank uint8
	PrgBank uint8

	irqEnabled bool
	tapePages [][]byte
	currentPage int

	//tape *sbox.StudyBox
	tape [][]byte
	tapePage int
	tapeOffset int

	readRegisters map[uint16]sbReadRegisterFunction
	writeRegisters map[uint16]sbWriteRegisterFunction
}

type sbReadRegisterFunction func() uint8
type sbWriteRegisterFunction func(value uint8)

func (sb *StudyBox) PageCount() int {
	return len(sb.tape)
}

func NewStudyBox(raw []byte, hasRam bool) (Mapper, error) {
	sb := &StudyBox{
		rom: raw,
		//mainRam: [0x0800]byte{},
		//ramA: [0x8000]byte{},
		//ramB: [0x8000]byte{},
		writeRegisters: map[uint16]sbWriteRegisterFunction{},
		readRegisters: map[uint16]sbReadRegisterFunction{},
	}

	sb.writeRegisters[0x4200] = sb.write4200
	sb.writeRegisters[0x4201] = sb.write4201

	sb.readRegisters[0x4200] = sb.read4200
	sb.readRegisters[0x4201] = sb.read4201

	return sb, nil
}

// Read a byte from tape
func (sb *StudyBox) read4200() uint8 {
	if sb.tape == nil {
		panic("[read4200()] No tape is opened")
	}
	return 0
}

// Tape drive status stuff
func (sb *StudyBox) read4201() uint8 {
	// TODO
	return 0
}

// Set RAM bank
func (sb *StudyBox) write4200(value uint8) {
	sb.RamABank = value & 0x03
	sb.RamBBank = value & 0xC0
	return
}

// Set ROM bank
func (sb *StudyBox) write4201(value uint8) {
	sb.PrgBank = value & 0x0F
}

func (sb *StudyBox) LoadTape(reader io.Reader) error {
	tape, err := sbox.Read(reader)
	if err != nil {
		return err
	}
	//sb.tape = tape
	sb.readTape(tape)
	return nil
}

// Open a .studybox tape data file
func (sb *StudyBox) LoadTapeFile(filename string) error {
	tape, err := sbox.ReadFile(filename)
	if err != nil {
		return err
	}
	//sb.tape = tape
	sb.readTape(tape)
	return nil
}

func (sb *StudyBox) readTape(t *sbox.StudyBox) {
	sb.tape = [][]byte{}
	for _, page := range t.Data.Pages {
		pg := []byte{}
		for _, packet := range page.Packets {
			pg = append(pg, packet.RawBytes()...)
		}
		sb.tape = append(sb.tape, pg)
	}
}

func (sb *StudyBox) ReadByte(address uint16) uint8 {
	// Handle registers first
	if reg, ok := sb.readRegisters[address]; ok {
		return reg()
	}

	if address < 0x2000 {
		return sb.mainRam[address % 0x0800]

	} else if address < 0x4200 {
		return 0

	} else if address >= 0x4400 && address <= 0x4FFF {
		// fixed RAM A
		return sb.ramA[address % 0x1000]

	} else if address >= 0x5000 && address <= 0x5FFF {
		// Switched RAM A (32K)
		return sb.ramA[address % 0x1000 + (uint16(sb.RamABank) * 0x1000)]

	} else if address >= 0x6000 && address <= 0x7FFF {
		// Switched RAM B (32k)
		return sb.ramA[address % 0x2000 + (uint16(sb.RamBBank) * 0x2000)]

	} else if address >= 0x8000 && address < 0xC000 {
		// Switched ROM (256k)
		return sb.rom[(address % 0x4000) + (uint16(sb.PrgBank) * 0x4000)]

	} else if address >= 0xC000 {
		// Fixed ROM bank 0
		return sb.rom[address % 0x4000]

	} else {
		panic(fmt.Sprintf("Read from address that isn't mapped to anything: %04X", address))
		// between 0x4203 and 0x4400, exclusively
		// do nothing
	}

	return 0
}

func (sb *StudyBox) WriteByte(address uint16, value uint8) {
	if reg, ok := sb.writeRegisters[address]; ok {
		reg(value)
		return
	}

	if address < 0x2000 {
		sb.mainRam[address % 0x0800] = value

	} else if address >= 0x4400 && address <= 0x4FFF {
		// fixed RAM A
		sb.ramA[address % 0x1000] = value

	} else if address >= 0x5000 && address <= 0x5FFF {
		// Switched RAM A (32K)
		sb.ramA[address % 0x1000 + (uint16(sb.RamABank) * 0x1000)] = value

	} else if address >= 0x6000 && address <= 0x7FFF {
		// Switched RAM B (32k)
		sb.ramA[address % 0x2000 + (uint16(sb.RamBBank) * 0x2000)] = value

	//} else if address >= 0x8000 && address < 0xC000 {
	//	// Switched ROM (256k)
	//	sb.rom[(address % 0x4000) + (uint16(sb.PrgBank) * 0x4000)] = value

	//} else if address >= 0xC000 {
	//	// Fixed ROM bank 0
	//	sb.rom[address % 0x4000] = value

	} else {
		//panic(fmt.Sprintf("Read from address that isn't mapped to anything: %04X", address))
		// between 0x4203 and 0x4400, exclusively
		// do nothing
	}
}

func (sb *StudyBox) ReadWord(address uint16) uint16 {
	return uint16(sb.ReadByte(address)) | (uint16(sb.ReadByte(address+1)) << 8)
}

func (sb *StudyBox) Offset(address uint16) uint32 {
	panic("Offset not implemented yet")
	return 0
}

func (sb *StudyBox) GetState() interface{} {
	return nil
}

func (sb *StudyBox) SetState(data interface{}) error {
	return fmt.Errorf("SetState() not implemented for StudyBox mapper")
}

func (sb StudyBox) Info() Info {
	panic("StudyBox.Info() not implemented yet")
	return Info{}
}

func (sb StudyBox) Name() string {
	return "StudyBox"
}

func (sb StudyBox) State() string {
	return "State() not implemented for StudyBox mapper"
}

func (sb StudyBox) DumpFullStack() string {
	st := []string{}
	for i := 0; i < 256; i++ {
		st = append(st, fmt.Sprintf("$%02X", sb.mainRam[0x100+i]))
	}
	return strings.Join(st, " ")
}

func (sb *StudyBox) ClearRam() {
	sb.mainRam = [0x0800]byte{}
	sb.ramA = [0x4000]byte{}
	sb.ramB = [0x8000]byte{}
}

func (sb *StudyBox) MemoryType(address uint16) string {
	if address >= 0x4400 && address < 0x8000 {
		return "NesWorkRam"
	}
	return "NesPrgRom"
}

func (sb *StudyBox) RomRead(offset uint) byte {
	panic("RomRead() not implemented")
}
