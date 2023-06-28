package mmu

import (
	"github.com/zorchenhimer/emu-6502/mappers"
)

type NES struct {
	mapper mappers.Mapper
	ram [0x0800]byte
}

func NewNES(mapper mappers.Mapper) *NES {
	return &NES{
		mapper: mapper,
		ram: [0x0800]byte{},
	}
}

func (n *NES) ReadByte(address uint16) uint8 {
	if address < 0x2000 {
		return n.ram[address % 0x0800]
	} else if address >= 0x4020 { // $4020 is the start of cart space
		return n.mapper.ReadByte(address)
	}

	return 0
}

func (n *NES) WriteByte(address uint16, value uint8) {
	if address < 0x2000 {
		n.ram[address % 0x0800] = value
	} else if address >= 0x4020 { // $4020 is the start of cart space
		n.mapper.WriteByte(address, value)
	}
}

func (n *NES) ClearRam() {
	for i := 0; i < len(n.ram); i++ {
		n.ram[i] = 0
	}

	n.mapper.ClearRam()
}
