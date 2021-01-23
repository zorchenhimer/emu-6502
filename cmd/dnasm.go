package main

import (
	"fmt"
	"io/ioutil"

	"github.com/zorchenhimer/emu-6502/dnasm"
)

func main() {
	raw, err := ioutil.ReadFile("runner.nes")
	if err != nil {
		fmt.Println("Unable to open ROM file:", err)
		return
	}

	d, err := dnasm.Disassemble(raw)
	if err != nil {
		fmt.Println("Unable to disassemble:", err)
		return
	}

	fmt.Println("Disassembly successful")

	err = d.WriteToFile("disassembled.asm")
	if err != nil {
		fmt.Println("Unable to write disassembly:", err)
	}

	fmt.Println("Done")
}
