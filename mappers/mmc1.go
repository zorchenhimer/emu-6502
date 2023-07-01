package mappers

import (
	"bytes"
	"fmt"
	"strings"
)

func init() {
	registerMapper(1, NewMMC1)
}

type MMC1 struct {
	rom []byte
	ram [0x0800]byte
	wram [0x2000]byte

	hasRam bool

	// 0 - one-screen, lower bank
	// 1 - one-srceen, upper bank
	// 2 - vertical
	// 3 - horizontal
	Mirroring uint8

	// 0, 1 - switch 32kb at $8000, ignoring low bit of bank number
	// 2 - fix first bank to $8000, switch 16kb bank at $C000
	// 3 - fix last bank at $C000, switch 16kb bank at $8000
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

func (m *MMC1) GetState() interface{} {
	state := &MMC1{
		hasRam: m.hasRam,
		Mirroring: m.Mirroring,
		PrgBankMode: m.PrgBankMode,
		ChrBankMode: m.ChrBankMode,
		ChrBank0: m.ChrBank0,
		ChrBank1: m.ChrBank1,
		PrgBank: m.PrgBank,
		shiftReg: m.shiftReg,
		shiftCount: m.shiftCount,

		// Pointer, not duplicated
		rom: m.rom,
		ram: [0x0800]byte{},
	}

	if m.hasRam {
		state.wram = [0x2000]byte{}
		wramCopy(&state.wram, &m.wram)
	}

	ramCopy(&state.ram, &m.ram)

	return state
}

func (m *MMC1) SetState(data interface{}) error {
	state, ok := data.(MMC1)
	if !ok {
		return fmt.Errorf("Invalid state given")
	}

	m.hasRam = state.hasRam
	m.Mirroring = state.Mirroring
	m.PrgBankMode = state.PrgBankMode
	m.ChrBankMode = state.ChrBankMode
	m.ChrBank0 = state.ChrBank0
	m.ChrBank1 = state.ChrBank1
	m.PrgBank = state.PrgBank
	m.shiftReg = state.shiftReg
	m.shiftCount = state.shiftCount

	if m.hasRam {
		m.wram = [0x2000]byte{}
		wramCopy(&m.wram, &state.wram)
	}
	ramCopy(&m.ram, &state.ram)

	return nil
}

func NewMMC1(data []byte, hasRam bool) (Mapper, error) {
	// FIXME: data doesn't account for CHR
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
		mmc1.wram = [0x2000]byte{}
	}

	if len(data) % 0x8000 != 0 {
		fmt.Printf("len(data): %d\nlen(data) %% 0x8000: %d\n", len(data), len(data) % 0x8000)
		return nil, ErrRomSize
	}

	return mmc1, nil
}

func (m *MMC1) Info() Info {
	info := Info{
		PrgSize: uint(len(m.rom)),
		PrgStartAddress: 0x8000,
		PrgRamStartAddress: 0x6000,

		// TODO: CHR stuff
		ChrSize: 0,
		ChrRamSize: 0,
		ChrBankSize: 0,
	}

	if m.hasRam {
		info.PrgRamSize = 0x2000
	}

	switch m.PrgBankMode {
	case 0, 1:
		info.PrgBankSize = 0x8000
	case 2, 3:
		info.PrgBankSize = 0x4000
	}

	return info
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

func (m *MMC1) ReadWord(address uint16) uint16 {
	return uint16(m.ReadByte(address)) | (uint16(m.ReadByte(address+1)) << 8)
}

func (m *MMC1) Offset(address uint16) uint32 {
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
			lastBank := uint32(len(m.rom) / 0x4000) - 1
			romAddr = (uint32(address) % 0x4000) + (0x4000 * lastBank)
		}
		//fmt.Printf("[3] %04X -> %08X\n", address, romAddr)
	default:
		panic(fmt.Sprintf("Invalid PrgBankMode: %02X"))
	}

	if int(romAddr) > len(m.rom) {
		panic(fmt.Sprintf("address out of range for ROM: $%04X -> 0x%06X; len: 0x%06X [%s]",
			address, romAddr, len(m.rom), m.State()))
	}

	return romAddr
}

func (m *MMC1) ReadByte(address uint16) uint8 {
	// RAM
	if address < 0x2000 {
		return m.ram[address % 0x0800]
	} else if address < 0x6000 {
		return 0
	} else if address < 0x8000 && m.hasRam {
		return m.wram[address - 0x6000]
	}

	return m.rom[m.Offset(address)]
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

func (m *MMC1) ClearRam() {
	m.ram = [0x0800]byte{}
	if m.hasRam {
		m.wram = [0x2000]byte{}
	}
}

func (m *MMC1) DumpFullStack() string {
	st := []string{}
	for i := 0; i < 256; i++ {
		st = append(st, fmt.Sprintf("$%02X", m.ram[0x100+i]))
	}
	return strings.Join(st, " ")
}

func (m *MMC1) MemoryType(address uint16) string {
	if m.hasRam && address >= 0x6000 && address < 0x8000 {
		return "NesWorkRam"
	} else if address >= 0x8000 {
		return "NesPrgRom"
	}
	return "NesOpenBus"
}
