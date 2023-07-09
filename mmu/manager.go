package mmu

import (
	"io"

	"github.com/zorchenhimer/emu-6502/labels"
)

type Manager interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)

	// Find label name by address
	GetLabel(address uint16) string
	GetZpLabel(address uint8) string

	// Find label address by name
	FindLabel(name string) (uint, labels.MemoryType)

	// Return all known labels for the given memory type
	Labels(t labels.MemoryType) labels.LabelMap

	AddDasm(address uint16, src string, size uint)
	//UpdateDasm(address uint16, instr 
	WriteDasm(writer io.Writer) error

	ClearRam()
}

