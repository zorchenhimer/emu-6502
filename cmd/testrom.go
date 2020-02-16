package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	//"strings"

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

	instructions(core)

	file, err := os.Create("debug.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	core.DebugFile = file
	// vectors have traps
	core.PC = 0x8000
	//core.PC = 0x0400
	core.Debug = true

	err = core.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(core.Registers())
		fmt.Printf("Ticks: %d\n", core.Ticks())
		//core.DumpPage(0x01)
		//core.DumpPage(0x02)
		core.DumpMemoryToFile("memory.txt")
		return
	}
}

func instructions(core *emu.Core) {
	instr := core.Instructions()
	sort.Strings(instr)
	count := len(instr)

	//err = ioutil.WriteFile("instructions.txt", []byte(strings.Join(instr, "\n")), 0777)
	file, err := os.Create("instructions.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, i := range instr {
		fmt.Fprintln(file, i)
	}

	percent := (float32(count) / 151) * 100

	fmt.Fprintf(file, "Total implemented: %d/151: %0.2f%%\nUnimplemented: %d\n", count, percent, 151-count)
}
