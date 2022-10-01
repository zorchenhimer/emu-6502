package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	//"os"
	//"sort"
	//"strings"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/zorchenhimer/emu-6502"
	"github.com/zorchenhimer/emu-6502/mappers"
	//cc "github.com/zorchenhimer/emu-6502/cc65"
)

const (
	// Routine addresses
	LoadChildMap uint16 = 0xCB23
	LoadMap uint16 = 0xCE58
	CheckBrickCollide uint16 = 0x8C04
	CheckPointCollide uint16 = 0x8E66

	// input variables used in routines
	BallDirection uint16 = 0x0080
	BallX uint16 = 0x0082
	BallY uint16 = 0x0084
	CurrentBoard uint16 = 0x00A0
	game_BoardOffsetX uint16 = 0x0096
	game_BoardOffsetY uint16 = 0x0095
	game_BoardHeight uint16 = 0x0098
	game_BoardWidth uint16 = 0x0097
	TmpX uint16 = 0x0010
	TmpY uint16 = 0x0011
	AddressPointer0 uint16 = 0x0000
)

func main() {
	//rom, err := ioutil.ReadFile("breakout.nes")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// TODO: read NES header

	//mapper, err := mappers.NewMMC1(rom[16:len(rom)], true)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	mapper, err := mappers.LoadFromFile("breakout.nes")
	if err != nil {
		fmt.Println(err)
		return
	}

	core, err := emu.NewCore(mapper)
	if err != nil {
		fmt.Println(err)
		return
	}

	//dbg, err := cc.ParseDebugFile("breakout.dbg")
	//if err != nil {
	//	fmt.Println("unable to parse dbg file:" + err.Error())
	//	return
	//}

	//addr_BallX, err := dbg.LabelAddress("BallX")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Printf("addr_BallX: %04X\n", addr_BallX)
	//return

	//for _, sym := range dbg.Symbols {
	//	fmt.Printf("%04X %s\n", sym.Val, sym.Name)
	//}
	//return

	// The current CheckBrickCollide routine
	//core.PC = 0x8C04
	//core.Debug = false

	//err = core.Run()
	//err = core.RunRoutine(CheckBrickCollide)
	//for i := uint8(0); i < 16; i++ {
	//	core.A = i
		//err = core.RunRoutine(LoadChildMap)
		//if err != nil {
		//	fmt.Println(core.Registers)
		//	fmt.Println(err)
		//	return
		//}
	//}

	var stop bool

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			stop = true
		}
	}()

	/*
		brick position
			ball position
				ball direction
	*/
	start := time.Now()
	for brickIdx := uint16(0); brickIdx < (6 * 12) - 1; brickIdx++ {
		if stop {
			fmt.Println("Stop received")
			fmt.Println(core.Registers())
			fmt.Printf("last brick index: %d\n", brickIdx - 1)
			break
		}
		core.HardReset()
		core.WriteByte(brickIdx + 0x6000, 0x41)
		core.WriteByte(brickIdx + 0x6000 + 1, 0x80)
		core.WriteByte(CurrentBoard, 0x80)	// set child board

		// magic numbers are wall locations
		for X := uint8(0x0A); X < 0xF5; X++ {
			core.WriteByte(BallX, X)
			core.WriteByte(TmpX, X)

			// bricks should never be below $70 pixels
			for Y := uint8(0x11); Y < 0x70; Y++ {
				core.WriteByte(BallY, Y)
				core.WriteByte(TmpY, Y)

				for D := uint8(0); D < 4; D++ {
					core.WriteByte(BallDirection, D)
					err = core.RunRoutine(
					err = core.RunRoutine(CheckPointCollide)
					if err != nil {
						fmt.Println(err)
						fmt.Println(core.Registers())
						return
					}

					brickAddr := core.ReadWord(AddressPointer0)
					if brickAddr != 0 && brickAddr != brickIdx + 0x6000 && brickAddr != brickIdx + 0x6000 + 1 {
						fmt.Println(core.Registers())
						fmt.Printf("invalid AddressPointer0.  Brick: $%04X; pointer: $%04X",
							brickIdx + 0x6000, brickAddr)
						return
					}
				}
			}
		}
	}
	fmt.Printf("time: %s\n", time.Now().Sub(start))

	ram := []byte{}
	for i := uint16(0x6000); i < 0x6100; i++ {
		ram = append(ram, core.ReadByte(i))
	}

	err = ioutil.WriteFile("breakout.ram", ram, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("xxd", "-u", "-o 24576", "breakout.ram")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("breakout.ram.txt", out.Bytes(), 0777)
	if err != nil {
		fmt.Println(err)
		return
	}


	fmt.Println("Done")
}
