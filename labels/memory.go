package labels

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
