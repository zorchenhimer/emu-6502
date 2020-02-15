package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zorchenhimer/emu-6502"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing rom")
		return
	}

	rom, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	core, err := emu.NewRWCore(rom, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// vectors have traps
	core.PC = 0x8000
	core.Debug = true

	err = core.Run()
	if err != nil {
		fmt.Println(err)
		core.DumpRegisters()
		fmt.Printf("Ticks: %d\n", core.Ticks())
		core.DumpPage(0x01)
		//core.DumpMemoryToFile("dump.txt")
		return
	}
}
