package mmu

import (
	"io"
	"sort"
	"fmt"

	"github.com/zorchenhimer/emu-6502/labels"
	"github.com/zorchenhimer/emu-6502/mappers"
	//"github.com/zorchenhimer/go-nes/mesen"
)

type Disassembly struct {
	//OpCode byte
	Value string
	//Raw []byte
	Instr Instruction
	Argument uint
	Size uint
	Address uint // address in bank space, not CPU space
}

type NES struct {
	mapper mappers.Mapper
	ram [0x0800]byte
	labels map[labels.MemoryType]labels.LabelMap

	dasmRom map[uint]string
	dasmRam map[uint]string
	ppu2ktoggle bool

	dasm []*Disassembly
	dasmVisited map[uint]map[uint16]any // list of CPU addresses
}

func NewNES(mapper mappers.Mapper) *NES {
	n := &NES{
		mapper: mapper,
		ram: [0x0800]byte{},

		dasmRom: make(map[uint]string),
		dasmRam: make(map[uint]string),
		dasm: make([]*Disassembly, mapper.Info().PrgSize),
		labels: make(map[labels.MemoryType]labels.LabelMap),
		dasmVisited: make(map[uint]map[uint16]any),
	}

	n.labels[labels.NesPrgRom] = make(labels.LabelMap)

	return n
}

func (n *NES) Visited(address uint16) bool {
	_, visited := n.dasmVisited[n.Offset(address)]
	return visited
}

// Mark and address as visited and record the CPU address it was mapped to
func (n *NES) Visit(address uint16) {
	offset := n.Offset(address)
	if _, ok := n.dasmVisited[offset]; !ok {
		n.dasmVisited[offset] = make(map[uint16]any)
	}
	n.dasmVisited[offset][address] = nil
}

func (n *NES) Offset(address uint16) uint {
	return uint(n.mapper.Offset(address))
}

func (n *NES) ReadByte(address uint16) uint8 {
	if address < 0x2000 {
		return n.ram[address % 0x0800]
	} else if address >= 0x4020 { // $4020 is the start of cart space
		return n.mapper.ReadByte(address)
	} else if address == 0x2002 {
		n.ppu2ktoggle = !n.ppu2ktoggle
		if n.ppu2ktoggle {
			return 0x80
		} else {
			return 0x00
		}
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

func (n *NES) GetZpLabel(address uint8) string {
	lbl := n.lookupLabel(uint16(address))
	if lbl != "" {
		return lbl
	}
	return fmt.Sprintf("$%02X", address)
}

// rom address space
func (n *NES) GetLabel(address uint) string {
	if lbl, ok := n.labels[labels.NesPrgRom][address]; ok {
		return lbl.Name
	}
	return fmt.Sprintf("$%04X", address)
}

// TODO: this needs to be in rom space, not CPU space
func (n *NES) lookupLabel(address uint16) string {
	switch n.MemoryType(address) {
	case labels.NesInternalRam:
		if lbl, ok := n.labels[labels.NesInternalRam][uint(address%0x800)]; ok {
			return lbl.Name
		}
	case labels.NesPrgRom, labels.NesWorkRam, labels.NesSaveRam:
		if lbl, ok := n.labels[labels.MemoryType(n.MemoryType(address))][uint(n.mapper.Offset(address))]; ok {
			return lbl.Name
		}
	}
	return ""
}

func (n *NES) FindLabel(name string) (uint, labels.MemoryType) {
	for t, list := range n.labels {
		addr, found := list.FindLabel(name)
		if found {
			return addr, t
		}
	}

	return 0, labels.NesOpenBus
}

func (n *NES) LoadLabelsMesen2(filename string) error {
	var err error
	n.labels, err = labels.LoadMesen2(filename)
	return err
}

func (n *NES) AddAutoLabel(address uint16) {
	if address < n.mapper.Info().PrgStartAddress {
		return
	}

	offset := n.mapper.Offset(address)
	lbl := n.labels[labels.NesPrgRom].GetLabel(uint(offset))
	if lbl != "" {
		return
	}

	n.labels[labels.NesPrgRom][uint(offset)] = &labels.Label{
		Name: fmt.Sprintf("L%06X", offset), Size: 1}
}

//func (n *NES) AddDasm(address uint16, src string, size uint) {
// CPU address, instruction, arugement (which is sometimes an address)
func (n *NES) AddDasm(address uint16, instr Instruction, arg uint) {
	offset := uint(n.mapper.Offset(address))
	if n.MemoryType(address) != labels.NesPrgRom {
		return
	}

	d := &Disassembly{
		//Value: src,
		Instr: instr,
		Argument: arg,
		Address: offset,
	}

	for i := uint(0); i < uint(d.Instr.Length()); i++ {
		n.dasm[i+offset] = d
	}
}

func (n *NES) MemoryType(address uint16) labels.MemoryType {
	if address >= 0x4020 {
		return labels.MemoryType(n.mapper.MemoryType(address))
	} else if address < 0x2000 {
		return labels.NesInternalRam
	}
	return labels.NesMemory
}

func (n *NES) WriteDasm(writer io.Writer) error {
	nothing := 0
	start := uint(0)
	for i := uint(0); i < uint(len(n.dasm)); i++ {
		// keep track of how long we haven't found anything
		if n.dasm[i] == nil {
			if nothing == 0 {
				start = i
			}
			nothing++
			continue
		}

		if nothing != 0 {
			// Write out raw bytes now that we found something.
			for j := start; j < i; j++ {
				err := n.writeline(writer, fmt.Sprintf("    .byte $%02X", n.mapper.RomRead(j)), fmt.Sprintf("$%06X", j))
				if err != nil {
					return err
				}
			}

			nothing = 0
		}

		// Code label
		if lbl, ok := n.labels[labels.NesPrgRom][i]; ok {
			// write comment if there is one
			if lbl.Comment != "" {
				err := n.writeline(writer, "", fmt.Sprintf("$%06X: %s", i, lbl.Comment))
				if err != nil {
					return err
				}
			}

			// write name if there is one (unlabeled comments are valid)
			if lbl.Name != "" {
				err := n.writeline(writer, lbl.Name+":", fmt.Sprintf("$%06X", i))
				if err != nil {
					return err
				}
			}
		}

		src := n.dasm[i].Instr.Decode(n.dasm[i].Argument, n)
		err := n.writeline(writer, "    "+src, fmt.Sprintf("$%06X", i))
		if err != nil {
			return err
		}

		if n.dasm[i].Instr.Length() > 1 {
			i = i+uint(n.dasm[i].Instr.Length())-1
		}
	}

	return nil
}

func (n *NES) writeline(w io.Writer, src, comment string) error {
	_, err := fmt.Fprintf(w, "%-30s ; %s\n", src, comment)
	return err
}

func (n *NES) WriteDasmWhat(writer io.Writer) error {
	addrs := []uint{}

	for addr, _ := range n.dasmRom {
		addrs = append(addrs, addr)
	}

	sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })

	for _, addr := range addrs {
		if lbl, ok := n.labels[labels.NesPrgRom][addr]; ok {
			if lbl.Comment != "" {
				_, err := fmt.Fprintln(writer, ";"+lbl.Comment)
				if err != nil {
					return err
				}
			}
			if lbl.Name != "" {
				_, err := fmt.Fprintln(writer, lbl.Name+":")
				if err != nil {
					return err
				}
			}
		}
		_, err := fmt.Fprintln(writer, n.dasmRom[addr])
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NES) Labels(t labels.MemoryType) labels.LabelMap {
	if l, ok := n.labels[t]; ok {
		return l
	}
	return nil

}
