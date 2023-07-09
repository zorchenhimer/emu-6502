package mmu

import (
	"fmt"
	"io"
	"sort"

	//"github.com/zorchenhimer/emu-6502/mappers"
	"github.com/zorchenhimer/emu-6502/labels"
)

type FullRam struct {
	ram [0x10000]byte
	lbmap labels.LabelMap
	dasm map[uint16]string
}

func NewFullRam(rombytes []byte) (*FullRam, error) {
	if len(rombytes) > 0x10000 {
		return nil, fmt.Errorf("rom too large")
	}

	fr := &FullRam{lbmap: make(labels.LabelMap), dasm: make(map[uint16]string)}
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

func (fr *FullRam) GetZpLabel(address uint8) string {
	return fmt.Sprintf("$%02X", address)
}

func (fr *FullRam) GetLabel(address uint16) string {
	return fmt.Sprintf("$%04X", address)
}

func (fr *FullRam) AddDasm(address uint16, src string, size uint) {
	//panic("AddDasm() not implemented for FullRam")
	fr.dasm[address] = src
}

func (fr *FullRam) WriteDasm(writer io.Writer) error {
	addrs := []uint16{}

	for addr, _ := range fr.dasm {
		addrs = append(addrs, addr)
	}

	sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })

	for _, addr := range addrs {
		_, err := fmt.Fprintln(writer, fr.dasm[addr])
		if err != nil {
			return err
		}
	}

	return nil
}
