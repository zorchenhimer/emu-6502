package mappers

import (
	"bytes"
	"fmt"
)

type MMC1 struct {
	rom []byte
	ram [0x0800]byte
	wram []byte

	hasRam bool

	// 0 - one-screen, lower bank
	// 1 - one-srceen, upper bank
	// 2 - vertical
	// 3 - horizontal
	Mirroring uint8

	// 0, 1 - switch 32kb at $8000, ignoring low bit of bank number
	// 2 - fix first bank to $8000, switch 16kb bank at $C000
	// 3 - fix last bank at $C000, switch 16kb bank at $C000
	PrgBankMode uint8

	// 0 - switch 8kb at a time
	// 1 - switch two separate 4kb banks
	ChrBankMode uint8

	ChrBank0 uint8
	ChrBank1 uint8
	PrgBank  uint8

	// Temporary values for shifting values
	shiftReg uint8
	shiftCount uint8
}

func NewMMC1(data []byte, hasRam bool) (Mapper, error) {
	//panic("Rewrite the NROM stuff")
	mmc1 := &MMC1{
		rom: data,
		ram: [0x0800]byte{},
		hasRam: hasRam,

		Mirroring:   0,
		PrgBankMode: 3,
		ChrBankMode: 8,

		ChrBank0: 0,
		ChrBank1: 0,
		PrgBank:  0,
	}

	if hasRam {
		mmc1.wram = make([]byte, 0x2000)
	}

	if len(data) % 0x8000 != 0 {
		return nil, ErrRomSize
	}

	return mmc1, nil
}

func (m *MMC1) Name() string {
	return "MMC1"
}

func (m *MMC1) State() string {
	var out bytes.Buffer

	out.WriteString("MMC1 ")
	out.WriteString(fmt.Sprintf("PrgMode: %X ", m.PrgBankMode))
	out.WriteString(fmt.Sprintf("PrgBank: %02X ", m.PrgBank))
	out.WriteString(fmt.Sprintf("ChrMode: %X ", m.ChrBankMode))
	out.WriteString(fmt.Sprintf("ChrBank0: %02X ", m.ChrBank0))
	out.WriteString(fmt.Sprintf("ChrBank1: %02X ", m.ChrBank1))

	if m.hasRam {
		out.WriteString("w/ WRAM")
	}

	return out.String()
}

func (m *MMC1) ReadByte(address uint16) uint8 {
	if address < 0x2000 {
		return m.ram[address % 0x0800]
	} else if address < 0x6000 {
		return 0
	} else if address < 0x8000 && m.wram != nil {
		return m.wram[address - 0x6000]
	}

	romAddr := uint32(address - 0x8000)
	switch m.PrgBankMode {
	case 0, 1:
		romAddr = ((uint32(m.PrgBank) & 0xFE) * 0x4000) + romAddr
	case 2:
		if romAddr > 0x4000 {
			romAddr = (uint32(m.PrgBank) * 0x4000) + romAddr
		}
	case 3:
		if romAddr < 0x4000 {
			romAddr = (uint32(m.PrgBank) * 0x4000) + romAddr
		} else {
			romAddr = (14 * 0x4000) + romAddr
		}
	default:
		panic(fmt.Sprintf("Invalid PrgBankMode: %02X"))
	}

	if int(romAddr) > len(m.rom) {
		panic(fmt.Sprintf("address out of range for ROM: $%04X -> 0x%06X; len: 0x%06X [%s]",
			address, romAddr, len(m.rom), m.State()))
	}
	return m.rom[romAddr]
}

func (m *MMC1) WriteByte(address uint16, value uint8) {
	if address < 0x2000 {
		m.ram[address % 0x0800] = value
	} else if address < 0x6000 {
		// do nothing
	} else if 0x6000 <= address && address < 0x8000 {
		m.wram[address - 0x6000] = value
	} else {
		// reset mapper shift
		if value & 0x80 != 0 {
			m.shiftReg = 0
			m.shiftCount = 0
		}

		m.shiftReg |= value << m.shiftCount
		m.shiftCount++

		if m.shiftCount == 5 {
			m.newState(address)
		}
	}
}

func (m *MMC1) newState(address uint16) {
	defer func() {
		m.shiftReg = 0
		m.shiftCount = 0
	}()

	// Control
	if address < 0xA000 {
		m.Mirroring = m.shiftReg & 0x03
		m.PrgBankMode = (m.shiftReg & 0x1C) >> 2
		m.ChrBankMode = (m.shiftReg & 0x20) >> 4
		return
	}

	// CHR Bank 0
	if address < 0xC000 {
		m.ChrBank0 = m.shiftReg
		return
	}

	// CHR Bank 1
	if address < 0xE000 {
		m.ChrBank0 = m.shiftReg
		return
	}

	// PRG Bank
	m.PrgBank = m.shiftReg & 0x0F
	if m.shiftReg & 0x10 != 0 {
		m.hasRam = true
	}
}

