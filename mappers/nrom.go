package mappers

import (
	"bytes"
)

type NROM struct {
	rom []byte
	ram [0x0800]byte
	wram []byte

	hasRam bool
	isHalf bool // if true, mirror 0xC000
}

func NewNROM(data []byte, hasRam bool) (Mapper, error) {
	nrom := &NROM{
		rom: data,
		ram: [0x0800]byte{},
		hasRam: hasRam,
	}

	if hasRam {
		nrom.wram = make([]byte, 0x2000)
	}

	if len(data) == 0x4000 {
		nrom.isHalf = true
	} else if len(data) != 0x8000 {
		return nil, ErrRomSize
	}

	return nrom, nil
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


func (nr *NROM) ReadByte(address uint16) uint8 {
	if address < 0x2000 {
		return nr.ram[address % 0x0800]
	} else if address < 0x6000 {
		return 0
	} else if address < 0x8000 && nr.wram != nil {
		return nr.wram[address - 0x6000]
	}

	address -= 0x8000
	if nr.isHalf {
		address = address % 0x4000
	}

	return nr.rom[address]
}

func (nr *NROM) WriteByte(address uint16, value uint8) {
	if nr.hasRam && 0x6000 <= address && address < 0x8000 {
		nr.ram[address - 0x6000] = value
	}
}

func (nr *NROM) ClearRam() {
	nr.ram = [0x0800]byte{}
	if nr.hasRam {
		nr.wram = make([]byte, 0x2000)
	}
}
