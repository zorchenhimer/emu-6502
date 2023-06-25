package mappers

import (
	"bytes"
	"fmt"
	"strings"
)

func init() {
	registerMapper(0, NewNROM)
}

type NROM struct {
	rom []byte
	ram [0x0800]byte
	wram [0x2000]byte

	hasRam bool
	isHalf bool // if true, mirror 0xC000
}

func (nr *NROM) GetState() interface{} {
	state := &NROM{
		hasRam: nr.hasRam,
		isHalf: nr.isHalf,

		rom: nr.rom,
		ram: [0x0800]byte{},
	}

	if nr.hasRam {
		state.wram = [0x2000]byte{}
		wramCopy(&state.wram, &nr.wram)
	}

	ramCopy(&state.ram, &nr.ram)

	return state
}

func (nr *NROM) SetState(data interface{}) error {
	state, ok := data.(NROM)
	if !ok {
		return fmt.Errorf("Invalid state given")
	}

	nr.hasRam = state.hasRam
	nr.isHalf = state.isHalf

	if nr.hasRam {
		wramCopy(&nr.wram, &state.wram)
	}

	ramCopy(&nr.ram, &state.ram)

	return nil
}

func NewNROM(data []byte, hasRam bool) (Mapper, error) {
	nrom := &NROM{
		rom: data,
		ram: [0x0800]byte{},
		hasRam: hasRam,
	}

	if hasRam {
		nrom.wram = [0x2000]byte{}
	}

	if len(data) == 0x4000 {
		nrom.isHalf = true
	} else if len(data) != 0x8000 {
		return nil, ErrRomSize
	}

	return nrom, nil
}

func (nr *NROM) Info() Info {
	info := Info{
		PrgRamStartAddress: 0x6000,

		// TODO: CHR stuff
		ChrSize: 0,
		ChrRamSize: 0,
		ChrBankSize: 0,
	}

	if nr.hasRam {
		info.PrgRamSize = 0x2000
	}

	if nr.isHalf {
		info.PrgSize = 0x4000
		info.PrgBankSize = 0x4000
		info.PrgStartAddress = 0xC000
	} else {
		info.PrgSize = 0x8000
		info.PrgBankSize = 0x8000
		info.PrgStartAddress = 0x8000
	}

	return info
}

func (nr *NROM) Name() string {
	return "NROM"
}

func (nr *NROM) State() string {
	var out bytes.Buffer

	out.WriteString("NROM ")

	if nr.isHalf {
		out.WriteString("16k ")
	} else {
		out.WriteString("32k ")
	}

	if nr.hasRam {
		out.WriteString("w/ WRAM")
	}

	return out.String()
}

func (nr *NROM) Offset(address uint16) uint32 {
	// Minus 8k to put the ROM start at the start of the
	// address space, plus 16 to account for the header.
	return uint32(address) - 0x8000 + 16
}

func (nr *NROM) ReadWord(address uint16) uint16 {
	return uint16(nr.ReadByte(address)) | (uint16(nr.ReadByte(address+1)) << 8)
}

func (nr *NROM) ReadByte(address uint16) uint8 {
	if address < 0x2000 {
		return nr.ram[address % 0x0800]
	} else if address < 0x6000 {
		return 0
	} else if address < 0x8000 && nr.hasRam {
		return nr.wram[address - 0x6000]
	}

	address -= 0x8000
	if nr.isHalf {
		address = address % 0x4000
	}

	return nr.rom[address]
}

func (nr *NROM) WriteByte(address uint16, value uint8) {
	if address < 0x2000 {
		nr.ram[address % 0x0800] = value
	} else if nr.hasRam && 0x6000 <= address && address < 0x8000 {
		nr.wram[address - 0x6000] = value
	}
}

func (nr *NROM) ClearRam() {
	nr.ram = [0x0800]byte{}
	if nr.hasRam {
		nr.wram = [0x2000]byte{}
	}
}

func (nr *NROM) DumpFullStack() string {
	st := []string{}
	for i := 0; i < 256; i++ {
		st = append(st, fmt.Sprintf("$%02X", nr.ram[0x100+i]))
	}
	return strings.Join(st, " ")
}
