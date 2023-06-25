package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	//"strings"

	"github.com/zorchenhimer/emu-6502"
	"github.com/zorchenhimer/emu-6502/mmu"
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

	mem, err := mmu.NewFullRam(rom)
	if err != nil {
		fmt.Println(err)
		return
	}

	core, err := emu.NewCore(mem)
	if err != nil {
		fmt.Println(err)
		return
	}

	//core.Breakpoints.Register(emu.EXECUTE, "Test success!", 0x3D78, func(c *emu.Core, event uint8, value uint8) {
	//	fmt.Println("\nTESTS PASS!\n")
	//	c.Halt()
	//})

	core.PC = 0x0400
	core.CheckStuck = true

	//core.DebugFile = dbgFile

	err = core.DumpMemoryToFile("before.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = core.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(core.Registers())
		fmt.Printf("Ticks: %d\n", core.Ticks())
		//core.DumpPage(0x01)
		//core.DumpPage(0x02)
		core.DumpMemoryToFile("memory.txt")

		if core.PC == 0x3D78 {
			fmt.Println("\nTests Passed!\n")
		} else {
			fmt.Println("\nTests Failed!\n")
		}

		return
	}

	err = core.DumpMemoryToFile("after.txt")
	if err != nil {
		fmt.Println(err)
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
