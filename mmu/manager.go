package mmu

import (
	"io"

	"github.com/zorchenhimer/emu-6502/labels"
)

type Manager interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)

	Offset(address uint16) uint

	// Find label name by address (rom addr space)
	GetLabel(address uint) string
	// ram space
	GetZpLabel(address uint8) string

	// Find label address by name
	FindLabel(name string) (uint, labels.MemoryType)

	// Add a label if one doesn't already exist
	AddAutoLabel(address uint16)

	// Same as above, but using PRG address space instead of CPU space.
	//AddAutoLabelPRG(name string, address uint)

	// Return all known labels for the given memory type
	Labels(t labels.MemoryType) labels.LabelMap

	AddDasm(address uint16, instruction Instruction, arg uint)
	//UpdateDasm(address uint16, instr 
	WriteDasm(writer io.Writer) error

	ClearRam()
}

type Instruction interface {
	Name() string
	Length() uint8
	Decode(arg uint, mem Manager) string
}

