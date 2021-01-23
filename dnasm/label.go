package dnasm

import (
	"fmt"
)

type LabelType int
const (
	LT_Branch LabelType = iota
	LT_Jump
	LT_Jsr
	LT_Vector
	LT_Data
	LT_Ram
)

type LabelMeta struct {
	Address uint16
	Offset uint32
	Type LabelType
}

func (lm LabelMeta) Label() string {
	t := "L"

	switch lm.Type {
	case LT_Branch:
		t = "B"
	case LT_Vector:
		t = "V"
	//case TT_Jsr:
	//	t = "R"
	}

	return fmt.Sprintf("%s%08X", t, lm.Offset)
}

type RamLabelType int
const (
	RLT_Byte RamLabelType = iota
	RLT_Pointer
	RLT_Table
)

type RamLabelMeta struct {
	Address uint16
	Type RamLabelType

	Used map[uint32]interface{}
}
