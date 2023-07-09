package mmu

import (
	"io"
	"sort"
	"fmt"

	"github.com/zorchenhimer/emu-6502/labels"
	"github.com/zorchenhimer/emu-6502/mappers"
	//dis "github.com/zorchenhimer/emu-6502/disassembly"
	//"github.com/zorchenhimer/go-nes/mesen"
)

type Disassembly struct {
	//OpCode byte
	Value string
	Size uint
	Address uint // address in bank space, not CPU space
}

type NES struct {
	mapper mappers.Mapper
	ram [0x0800]byte
	labels map[labels.MemoryType]labels.LabelMap

	dasmRom map[uint]string
	dasmRam map[uint]string

	dasm []*Disassembly
}

func NewNES(mapper mappers.Mapper) *NES {
	return &NES{
		mapper: mapper,
		ram: [0x0800]byte{},

		dasmRom: make(map[uint]string),
		dasmRam: make(map[uint]string),
		dasm: make([]*Disassembly, mapper.Info().PrgSize),
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

func (n *NES) GetZpLabel(address uint8) string {
	lbl := n.lookupLabel(uint16(address))
	if lbl != "" {
		return lbl
	}
	return fmt.Sprintf("$%02X", address)
}

func (n *NES) GetLabel(address uint16) string {
	lbl := n.lookupLabel(address)
	if lbl != "" {
		return lbl
	}
	return fmt.Sprintf("$%04X", address)
}

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

func (n *NES) AddDasm(address uint16, src string, size uint) {
	offset := uint(n.mapper.Offset(address))
	if n.MemoryType(address) != labels.NesPrgRom {
		return
	}

	instr := &Disassembly{
		Value: src,
		Address: offset,
		Size: size,
	}

	for i := uint(0); i < size; i++ {
		n.dasm[i+offset] = instr
	}

	//n.dasm.Add(&dis.Instruction{
	//	Address: offset,
	//	Value: src,
	//	Size: size,
	//})

	//switch n.MemoryType(address) {
	//case NesWorkRam:
	//	n.dasmRam[offset] = src
	//case NesPrgRom:
	//	n.dasmRom[offset] = src
	//default:
	//	// do nothing for now
	//}
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
		if n.dasm[i] == nil {
			if nothing == 0 {
				start = i
			}
			nothing++
			continue
		}

		if nothing != 0 {
			//err := n.writeline(writer, "", fmt.Sprintf("$%06X: unknown for %d bytes", start, nothing))
			////_, err := fmt.Fprintf(writer, "; $%06X: unknown for %d bytes\n", start, nothing)
			//if err != nil {
			//	return err
			//}

			for j := start; j < i; j++ {
				err := n.writeline(writer, fmt.Sprintf("    .byte $%02X", n.mapper.RomRead(j)), fmt.Sprintf("$%06X", j))
				//_, err := fmt.Fprintf(writer, "    .byte $%02X ; $%06X \n", n.mapper.RomRead(j), j)
				if err != nil {
					return err
				}
			}

			nothing = 0
		}

		if lbl, ok := n.labels[labels.NesPrgRom][i]; ok {
			if lbl.Comment != "" {
				err := n.writeline(writer, "", fmt.Sprintf("$%06X: %s", i, lbl.Comment))
				if err != nil {
					return err
				}
			}
			if lbl.Name != "" {
				//_, err := fmt.Fprintf(writer, "%s: ; $%06X\n", lbl.Name, i)
				err := n.writeline(writer, lbl.Name+":", fmt.Sprintf("$%06X", i))
				if err != nil {
					return err
				}
			}
		}

		err := n.writeline(writer, "    "+n.dasm[i].Value, fmt.Sprintf("$%06X", i))
		//_, err := fmt.Fprintf(writer, "  %-10s ; $%06X\n", n.dasm[i].Value, i)
		if err != nil {
			return err
		}

		if n.dasm[i].Size > 1 {
			i = i+n.dasm[i].Size-1
		}
	}

	return nil
	//return n.dasm.Write(writer)
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
