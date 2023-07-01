package mmu

import (
	"io"
)

type Manager interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)

	GetLabel(address uint16) string

	AddDasm(address uint16, src string)
	WriteDasm(writer io.Writer) error

	ClearRam()
}

type MemoryType string

const (
	NesChrRam             MemoryType = "NesChrRam"
	NesChrRom             MemoryType = "NesChrRom"
	NesInternalRam        MemoryType = "NesInternalRam"
	NesMemory             MemoryType = "NesMemory"
	NesNametableRam       MemoryType = "NesNametableRam"
	NesPaletteRam         MemoryType = "NesPaletteRam"
	NesPrgRom             MemoryType = "NesPrgRom"
	NesSaveRam            MemoryType = "NesSaveRam"
	NesSecondarySpriteRam MemoryType = "NesSecondarySpriteRam"
	NesSpriteRam          MemoryType = "NesSpriteRam"
	NesWorkRam            MemoryType = "NesWorkRam"
	NesOpenBus            MemoryType = "NesOpenBus"
)
