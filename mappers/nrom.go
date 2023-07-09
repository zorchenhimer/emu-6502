package mappers

import (
	"bytes"
	"fmt"
)

func init() {
	registerMapper(0, NewNROM)
}

type NROM struct {
	rom []byte
	wram [0x2000]byte

	hasRam bool
	isHalf bool // if true, mirror 0xC000
}

func (nr *NROM) GetState() interface{} {
	state := &NROM{
		hasRam: nr.hasRam,
		isHalf: nr.isHalf,

		rom: nr.rom,
	}

	if nr.hasRam {
		state.wram = [0x2000]byte{}
		wramCopy(&state.wram, &nr.wram)
	}


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

	return nil
}

func NewNROM(data []byte, hasRam bool) (Mapper, error) {
	nrom := &NROM{
		rom: data,
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
	// This one's wram, probably
	if address < 0x8000 {
		return uint32(address - 0x6000)
	}

	// Minus 8k to put the ROM start at the start of the
	// address space.
	return uint32(address) - 0x8000
}

func (nr *NROM) MemoryType(address uint16) string {
	if address >= 0x8000 {
		return "NesPrgRom"
	} else if address >= 0x6000 {
		return "NesWorkRam"
	}
	return "NesOpenBus"
}

func (nr *NROM) ReadByte(address uint16) uint8 {
	if nr.hasRam && address >= 0x6000 && address < 0x8000 {
		return nr.wram[address - 0x6000]
	} else if address >= 0x8000 {
		address -= 0x8000
		if nr.isHalf {
			address = address % 0x4000
		}
		return nr.rom[address]
	}

	return 0
}

func (nr *NROM) WriteByte(address uint16, value uint8) {
	if nr.hasRam && 0x6000 <= address && address < 0x8000 {
		nr.wram[address - 0x6000] = value
	}
}

func (nr *NROM) ClearRam() {
	if nr.hasRam {
		nr.wram = [0x2000]byte{}
	}
}

func (nr *NROM) RomRead(offset uint) byte {
	if offset > uint(len(nr.rom)) {
		return 0
	}

	return nr.rom[offset]
}
