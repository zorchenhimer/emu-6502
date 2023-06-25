package mmu

import (
	"fmt"

	//"github.com/zorchenhimer/emu-6502/mappers"
)

type FullRam struct {
	ram [0x10000]byte
}

func NewFullRam(rombytes []byte) (*FullRam, error) {
	if len(rombytes) > 0x10000 {
		return nil, fmt.Errorf("rom too large")
	}

	fr := &FullRam{}
	for i, b := range rombytes {
	//for i := 0; i < len(rombytes); i++ {
		fr.ram[i] = b
	}

	return fr, nil
}

func (fr *FullRam) ReadByte(address uint16) uint8 {
	return fr.ram[address]
}

func (fr *FullRam) ReadWord(address uint16) uint16 {
	return uint16(fr.ram[address]) | (uint16(fr.ram[address+1]) << 8)
}

func (fr *FullRam) WriteByte(address uint16, value uint8) {
	fr.ram[address] = value
}

func (fr *FullRam) ClearRam() {
	// do nothing
}
