package mappers

import (
	"fmt"
)

// this "mapper" cannot be automatically detected and loaded.

type FullRW struct {
	rom []byte
}

func NewFullRW(data []byte) (Mapper, error) {
	if len(data) != 0x10000 {
		return nil, ErrRomSize
	}

	return &FullRW{
		rom: data,
	}, nil
}

func (rw *FullRW) GetState() interface{} {
	panic("\"Mapper\" FullRW does not support GetState()")
}

func (rw *FullRW) SetState(interface{}) error {
	return fmt.Errorf("\"Mapper\" FullRW does not support GetState()")
}

func (rw *FullRW) Info() Info {
	return Info{
		PrgSize: 0,
		PrgRamSize: 0x10000,
		PrgBankSize: 0x10000,

		ChrSize: 0,
		ChrRamSize: 0,
		ChrBankSize: 0,

		PrgStartAddress: 0,
		PrgRamStartAddress: 0,
	}
}

func (rw *FullRW) Name() string {
	return "FullRW"
}

func (rw *FullRW) State() string {
	return rw.Name()
}

func (rw *FullRW) Offset(address uint16) uint32 {
	// Minus 8k to put the ROM start at the start of the
	// address space.
	return uint32(address) - 0x8000
}

func (rw *FullRW) ReadWord(address uint16) uint16 {
	return uint16(rw.ReadByte(address)) | (uint16(rw.ReadByte(address+1)) << 8)
}

func (rw *FullRW) ReadByte(address uint16) uint8 {
	return rw.rom[address]
}

func (rw *FullRW) WriteByte(address uint16, value uint8) {
	rw.rom[address] = value
}

func (rw *FullRW) ClearRam() {}

func (rw *FullRW) DumpFullStack() string {
	return "FullRW DumpFullStack() not implemented"
}

func (rw *FullRW) MemoryType(address uint16) string {
	return "NesWorkRam"
}

func (rw *FullRW) RomRead(offset uint) byte {
	panic("RomRead() not implemented")
}
