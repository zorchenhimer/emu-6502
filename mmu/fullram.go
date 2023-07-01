package mmu

import (
	"fmt"
	"io"

	//"github.com/zorchenhimer/emu-6502/mappers"
	"github.com/zorchenhimer/emu-6502/labels"
)

type FullRam struct {
	ram [0x10000]byte
	lbmap labels.LabelMap
}

func NewFullRam(rombytes []byte) (*FullRam, error) {
	if len(rombytes) > 0x10000 {
		return nil, fmt.Errorf("rom too large")
	}

	fr := &FullRam{lbmap: make(labels.LabelMap)}
	for i, b := range rombytes {
	//for i := 0; i < len(rombytes); i++ {
		fr.ram[i] = b
	}

	return fr, nil
}

func (fr *FullRam) ReadByte(address uint16) uint8 {
	return fr.ram[address]
}

func (fr *FullRam) WriteByte(address uint16, value uint8) {
	fr.ram[address] = value
}

func (fr *FullRam) ClearRam() {
	// do nothing
}

func (fr *FullRam) GetLabel(address uint16) string {
	return ""
}

func (fr *FullRam) AddDasm(address uint16, src string) {
	panic("AddDasm() not implemented for FullRam")
}

func (fr *FullRam) WriteDasm(writer io.Writer) error {
	return fmt.Errorf("WriteDasm() not implemented for FullRam")
}
