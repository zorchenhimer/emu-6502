package main

import (
	"fmt"
	"io/ioutil"
	//"os"
	//"sort"
	//"strings"

	"github.com/zorchenhimer/emu-6502"
	"github.com/zorchenhimer/emu-6502/mappers"
)

// Routine addresses
const (
	LoadChildMap uint16 = 0xCB23
	CheckBrickCollide uint16 = 0x8C04
)

func main() {
	rom, err := ioutil.ReadFile("breakout.nes")
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: read NES header

	mapper, err := mappers.NewMMC1(rom[16:len(rom)], true)
	if err != nil {
		fmt.Println(err)
		return
	}

	core, err := emu.NewCore(mapper)
	if err != nil {
		fmt.Println(err)
		return
	}

	// The current CheckBrickCollide routine
	//core.PC = 0x8C04
	//core.Debug = false

	//err = core.Run()
	err = core.RunRoutine(CheckBrickCollide)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done")
}
