package dnasm

import (
	"github.com/zorchenhimer/emu-6502"
)

type Chunk struct {
	MapperState interface{}
	CpuState emu.CpuState
	FromJsr bool
	Address uint16

	FromNode Node
}

