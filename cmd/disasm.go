package main

import (
	"fmt"
	"os"
	"bytes"

	"github.com/zorchenhimer/emu-6502"
	"github.com/zorchenhimer/emu-6502/mappers"
)

const (
	CheckPointCollide uint16 = 0x8EC9
)

func main() {
	mapper, err := mappers.LoadFromFile("breakout.nes")
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := emu.NewCore(mapper)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Disassemble(CheckPointCollide)
	if err != nil {
		fmt.Println(err)
		return
	}

	buff := &bytes.Buffer{}
	err = c.WriteCdl(buff)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("dasm.cdl.dat", buff.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	buff.Reset()
	err = c.WriteVisited(buff)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("dasm.code.asm", buff.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
