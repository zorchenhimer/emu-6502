package dnasm

import (
	"fmt"

	"github.com/zorchenhimer/emu-6502"
)

type Chunk struct {
	MapperState interface{}
	CpuState emu.CpuState
	FromJsr bool
	Address uint16
}

// TODO: get rid of this.  Use Chunk.Address and create a new chunk for
//       each JSR instead of following it.  Chunk.FromJsr would need to
//       be utilized and copied "recursively"
type JsrStack struct {
	entries []*JsrElement
}

type JsrElement struct {
	CallAddress uint32
	RoutineAddress uint32
}

func (js *JsrStack) Push(val *JsrElement) {
	if js.entries == nil {
		js.entries = []*JsrElement{}
	}

	fmt.Printf("{push %08X -> %08X}\n", val.CallAddress, val.RoutineAddress)
	js.entries = append(js.entries, val)

	fmt.Println(js.entries)
}

func (js *JsrStack) Pop() (*JsrElement, error) {
	if js.entries == nil {
		return nil, fmt.Errorf("JsrStack Pop() before initial Push()")
	}

	if len(js.entries) == 0 {
		return nil, fmt.Errorf("JsrStack Pop() on empty stack")
	}

	fmt.Println(js.entries)
	val := js.entries[len(js.entries)-1]
	js.entries = js.entries[:len(js.entries)-1]

	fmt.Printf("{pop %08X <- %08X}\n", val.CallAddress, val.RoutineAddress)
	return val, nil
}

