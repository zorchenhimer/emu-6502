package mmu

import (
	"io"
	"sort"
	"fmt"

	"github.com/zorchenhimer/emu-6502/labels"
	"github.com/zorchenhimer/emu-6502/mappers"
	"github.com/zorchenhimer/go-nes/mesen"
)

type NES struct {
	mapper mappers.Mapper
	ram [0x0800]byte
	labels map[mesen.MemoryType]labels.LabelMap

	dasmRom map[uint]string
	dasmRam map[uint]string
}

func NewNES(mapper mappers.Mapper) *NES {
	return &NES{
		mapper: mapper,
		ram: [0x0800]byte{},

		dasmRom: make(map[uint]string),
		dasmRam: make(map[uint]string),
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

func (n *NES) GetLabel(address uint16) string {
	return ""
}

func (n *NES) LoadLabelsMesen2(filename string) error {
	var err error
	n.labels, err = labels.LoadMesen2(filename)
	return err
}

func (n *NES) AddDasm(address uint16, src string) {
	offset := uint(n.mapper.Offset(address))
	switch n.MemoryType(address) {
	case NesWorkRam:
		n.dasmRam[offset] = src
	case NesPrgRom:
		n.dasmRom[offset] = src
	default:
		// do nothing for now
	}
}

func (n *NES) MemoryType(address uint16) MemoryType {
	if address >= 0x4020 {
		return MemoryType(n.mapper.MemoryType(address))
	} else if address < 0x2000 {
		return NesInternalRam
	}
	return NesMemory
}

func (n *NES) WriteDasm(writer io.Writer) error {
	addrs := []uint{}

	for addr, _ := range n.dasmRom {
		addrs = append(addrs, addr)
	}

	sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })

	for _, addr := range addrs {
		_, err := fmt.Fprintln(writer, n.dasmRom[addr])
		if err != nil {
			return err
		}
	}

	return nil
}
