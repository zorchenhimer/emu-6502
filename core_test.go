package emu

import (
	"testing"
)

type regState struct {
	a      uint8
	x      uint8
	y      uint8
	pc     uint16
	phlags uint8
}

// Figure out how this'll work, lol
type ExpectedResults interface {
	GetExpected() regState
}

//type InitialState interface {
//    GetInitialState() regState
//}

// single test case.  single OP code, and register state
type basicTest struct {
	name string
	rom  []byte // should be no more than three bytes
	//ticks int // ticks to perform
	regInitial  regState
	regExpected regState
}

func (bt basicTest) GetExpected() {
	return bt.regExpected
}

//func (bt basicTest) GetInitialState() {
//	return bt.regInitial
//}

var testData_A = []basicTest{

	// LDA
	basicTest{
		"OP_LDA_IM",
		[]byte{OP_LDA_IM, 0x01},
		regState{},
		regState{0x01, 0x00, 0x00, 0x8002, 0x00}},
	basicTest{
		"OP_LDA_AB",
		[]byte{OP_LDA_AB, 0x00, 0x80},
		regState{},
		regState{OP_LDA_AB, 0x00, 0x00, 0x8003, 0x00}},
	basicTest{
		"OP_LDA_ZP",
		[]byte{OP_LDA_ZP, 0x01},
		regState{},
		regState{0x01, 0x00, 0x00, 0x8002, 0x00}},
	basicTest{
		"OP_LDA_AX",
		[]byte{OP_LDA_AX, 0xFD, 0x7F}, // address is three before OpCode
		regState{x: 0x03},
		regState{OP_LDA_AX, 0x03, 0x00, 0x8003, 0x00}},
	basicTest{
		"OP_LDA_AY",
		[]byte{OP_LDA_AY, 0xFD, 0x7F}, // address is three before OpCode
		regState{y: 0x03},
		regState{OP_LDA_AY, 0x00, 0x03, 0x8003, 0x00}},

	//    These will probably need to be in their own tests
	//    basicTest{
	//        "OP_LDA_IX",
	//        []byte{OP_LDA_IX, 0xFD, 0x7F}, // address is three before OpCode
	//        regState{x:0x03},
	//        regState{OP_LDA_IX, 0x03, 0x00, 0x8003, 0x00}},
	//    basicTest{
	//        "OP_LDA_AY",
	//        []byte{OP_LDA_IY, 0xFD, 0x7F}, // address is three before OpCode
	//        regState{y:0x03},
	//        regState{OP_LDA_IY, 0x00, 0x03, 0x8003, 0x00}},

	basicTest{
		"OP_NOP",
		[]byte{OP_NOP},
		regState{},
		regState{pc: 0x8001}},
}

func TestImmediate(t *testing.T) {
	core := newTestCore(t)
	for _, bt := range testData_A {
		t.Run(bt.name, func(t *testing.T) {
			err := core.resetTest(bt.rom)
			if err != nil {
				t.Errorf("%s: %v", bt.name, err)
			}

			core.setRegisters(t, bt.regInitial)

			ticksran := 0

			//for i := 0; i < bt.ticks; i++ {
			for !core.testDone {
				err = core.tick()
				if err != nil {
					t.Fatalf("%s: %v", bt.name, err)
				}
				ticksran++
			}
			ticksran -= 1 // remove the 0xFF "test done" tick

			core.checkRegisters(t, bt.name, bt)
			t.Logf("ticks ran: %d", ticksran)
			// TODO: phlags
		})
	}
}

type memVal struct {
	addr uint16
	val  byte
}

type memTest struct {
	name string
	rom  []byte

	// memory locations to check
	mem memVal

	regInitial  regState
	regExpected regState
}

func (mt memTest) GetExpected() {
	return mt.regExpected
}

//func (mt memTest) GetInitialState() {
//	return mt.regInitial
//}

// Tests that check memory values in addition to register flags
var memoryBased = memTest{
	// STA
	memTest{
		"OP_STA_AB",
		[]byte{OP_STA_AB, 0x00, 0x03},
		memVal{0x0300, OP_STA_AB},
		regState{a: OP_STA_AB},
		regState{OP_STA_AB, 0x00, 0x00, 0x8002, 0x00}},
}

func TestMemoryBased(t *testing.T) {
	core := newTestCore(t)
	for _, mt := range memoryBased {
		t.Run(mt.name, func(t *testing.T) {
			err := core.resetTest(mt.rom)
			if err != nil {
				t.Errorf("%s: %v", bt.name, err)
			}

			core.setRegisters(bt.regInitial)

			for !core.testDone {
				err = core.tick()
				if err != nil {
					t.Fatalf("%s: %v", bt.name, err)
				}
				ticksran++
			}
			ticksran -= 1 // remove the 0xFF "test done" tick

			core.checkRegisters(t, mt)

		})
	}
}

// Core helper functions
func (core *Core) checkRegisters(t *testing.T, name string, e ExpectedResults) {
	t.Helper()
	if core.A != e.GetExpected().a {
		t.Errorf("%s: Incorrect A: Exp:$%02X Got:$%02X", name, e.GetExpected().a, core.A)
	}

	if core.X != e.GetExpected().x {
		t.Errorf("%s: Incorrect X: Exp:$%02X Got:$%02X", name, e.GetExpected().x, core.X)
	}

	if core.Y != e.GetExpected().y {
		t.Errorf("%s: Incorrect Y: Exp:$%02X Got:$%02X", name, e.GetExpected().y, core.Y)
	}

	if core.PC != e.GetExpected().pc {
		t.Errorf("%s: Incorrect PC: Exp:$%04X Got:$%04X", name, e.GetExpected().pc, core.PC)
	}
}

func (core *Core) setRegisters(t *testing.T, r regState) {
	t.Helper()
	core.A = bt.regInitial.a
	core.X = bt.regInitial.x
	core.Y = bt.regInitial.y
	//core.PC = bt.regInitial.pc
	core.Phlags = bt.regInitial.phlags
}

func (c *Core) resetTest(t *testing.T, rom []byte) error {
	t.Helper()
	rom = padWithVectors(rom, 0x8000, 0x8000, 0x8000)
	if len(rom)%256 != 0 {
		return fmt.Errorf("ROM is not divisible by 256: %d", len(rom))
	}

	c.rom = rom
	c.PC = c.ReadWord(VECTOR_RESET)

	c.memory = make([]byte, 0x1000)

	c.testDone = false

	// fill zero page with some data
	for i := 0; i < 256; i++ {
		c.memory[i] = uint8(i)
	}

	return nil
}

func newTestCore(t *testing.T) *Core {
	t.Helper()
	return &Core{
		A:      0,
		X:      0,
		Y:      0,
		PC:     0,
		Phlags: 0,
		SP:     0,

		memory: make([]byte, 0x1000), // no registers, no WRAM, no ROM
		rom:    nil,

		InstructionLimit: 0,
		testing:          true,
		t:                t,
	}
}
