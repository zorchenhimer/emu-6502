package emu

import (
	"fmt"
	"strings"
	"testing"
)

var testsRun int = 0

// Tests that check memory values in addition to register flags
var memoryBased = []memTest{
	// DEC
	memTest{
		"OP_DEC_AB",
		[]byte{OP_DEC_AB, 0x00, 0x03},
		memVal{0x0300, 0xFF},
		regState{},
		regState{0x00, 0x00, 0x00, 0x8003, FLAG_NEGATIVE, 0x00}},
	memTest{
		"OP_DEC_AX",
		[]byte{OP_DEC_AX, 0x00, 0x03},
		memVal{0x0302, 0xFF},
		regState{x: 2},
		regState{0x00, 0x02, 0x00, 0x8003, FLAG_NEGATIVE, 0x00}},
	memTest{
		"OP_DEC_ZP",
		[]byte{OP_DEC_ZP, 0x03},
		memVal{0x0003, 0x02},
		regState{},
		regState{0x00, 0x00, 0x00, 0x8002, 0x00, 0x00}},
	memTest{
		"OP_DEC_ZX",
		[]byte{OP_DEC_ZX, 0x03},
		memVal{0x0005, 0x04},
		regState{x: 2},
		regState{0x00, 0x02, 0x00, 0x8002, 0x00, 0x00}},

	// INC
	memTest{
		"OP_INC_AB",
		[]byte{OP_INC_AB, 0x00, 0x03},
		memVal{0x0300, 0x01},
		regState{},
		regState{0x00, 0x00, 0x00, 0x8003, 0x00, 0x00}},
	memTest{
		"OP_INC_AX",
		[]byte{OP_INC_AX, 0x00, 0x03},
		memVal{0x0302, 0x01},
		regState{x: 2},
		regState{0x00, 0x02, 0x00, 0x8003, 0x00, 0x00}},
	memTest{
		"OP_INC_ZP",
		[]byte{OP_INC_ZP, 0x03},
		memVal{0x0003, 0x04},
		regState{},
		regState{0x00, 0x00, 0x00, 0x8002, 0x00, 0x00}},
	memTest{
		"OP_INC_ZX",
		[]byte{OP_INC_ZX, 0x03},
		memVal{0x0005, 0x06},
		regState{x: 2},
		regState{0x00, 0x02, 0x00, 0x8002, 0x00, 0x00}},

	// STA
	memTest{
		"OP_STA_AB",
		[]byte{OP_STA_AB, 0x00, 0x03},
		memVal{0x0300, OP_STA_AB},
		regState{a: OP_STA_AB},
		regState{OP_STA_AB, 0x00, 0x00, 0x8003, 0x00, 0x00}},

	memTest{
		"OP_STA_ZP",
		[]byte{OP_STA_ZP, 0x03},
		memVal{0x0003, OP_STA_ZP},
		regState{a: OP_STA_ZP},
		regState{OP_STA_ZP, 0x00, 0x00, 0x8002, 0x00, 0x00}},
	memTest{
		"OP_STA_ZX",
		[]byte{OP_STA_ZX, 0x03},
		memVal{0x0006, OP_STA_ZX},
		regState{a: OP_STA_ZX, x: 03},
		regState{OP_STA_ZX, 0x03, 0x00, 0x8002, 0x00, 0x00}},

	memTest{
		"OP_STA_AX",
		[]byte{OP_STA_AX, 0x03, 0x03},
		memVal{0x0306, OP_STA_AX},
		regState{a: OP_STA_AX, x: 03},
		regState{OP_STA_AX, 0x03, 0x00, 0x8003, 0x00, 0x00}},
	memTest{
		"OP_STA_AY",
		[]byte{OP_STA_AY, 0x03, 0x03},
		memVal{0x0306, OP_STA_AY},
		regState{a: OP_STA_AY, y: 03},
		regState{OP_STA_AY, 0x00, 0x03, 0x8003, 0x00, 0x00}},

	memTest{
		"OP_STA_IX",
		// pointer is at $0001 + 1 (x) = $0002
		// should be a pointer val of $0302
		[]byte{OP_STA_IX, 0x00},
		memVal{0x0302, OP_STA_IX},
		regState{a: OP_STA_IX, x: 02},
		regState{OP_STA_IX, 0x02, 0x00, 0x8002, 0x00, 0x00}},
	memTest{
		"OP_STA_IY",
		// pointer is $0002 = $0302 + 1 (y) = $0303
		[]byte{OP_STA_IY, 0x02},
		memVal{0x0303, OP_STA_IY},
		regState{a: OP_STA_IY, y: 01},
		regState{OP_STA_IY, 0x00, 0x01, 0x8002, 0x00, 0x00}},

	memTest{
		"OP_STX_AB",
		[]byte{OP_STX_AB, 0x00, 0x03},
		memVal{0x0300, OP_STX_AB},
		regState{x: OP_STX_AB},
		regState{x: OP_STX_AB, pc: 0x8003}},
	memTest{
		"OP_STX_ZP",
		[]byte{OP_STX_ZP, 0x03},
		memVal{0x0003, OP_STX_AB},
		regState{x: OP_STX_AB},
		regState{x: OP_STX_AB, pc: 0x8002}},
	memTest{
		"OP_STX_ZY",
		[]byte{OP_STX_ZY, 0x03},
		memVal{0x0006, OP_STX_AB},
		regState{x: OP_STX_AB, y: 3},
		regState{x: OP_STX_AB, y: 3, pc: 0x8002}},
}

var testData_A = []basicTest{

	// Decrements
	basicTest{
		"OP_DEX",
		[]byte{OP_DEX},
		regState{},
		regState{0x00, 0xFF, 0x00, 0x8001, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_DEY",
		[]byte{OP_DEY},
		regState{},
		regState{0x00, 0x00, 0xFF, 0x8001, FLAG_NEGATIVE, 0x00}},

	// Increments
	basicTest{
		"OP_DEX",
		[]byte{OP_INX},
		regState{},
		regState{0x00, 0x01, 0x00, 0x8001, 0x00, 0x00}},
	basicTest{
		"OP_DEY",
		[]byte{OP_INY},
		regState{},
		regState{0x00, 0x00, 0x01, 0x8001, 0x00, 0x00}},

	// LDA
	basicTest{
		"OP_LDA_IM",
		[]byte{OP_LDA_IM, 0x01},
		regState{},
		regState{0x01, 0x00, 0x00, 0x8002, 0x00, 0x00}},
	basicTest{
		"OP_LDA_AB",
		[]byte{OP_LDA_AB, 0x00, 0x80},
		regState{},
		regState{OP_LDA_AB, 0x00, 0x00, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDA_ZP",
		[]byte{OP_LDA_ZP, 0x01},
		regState{},
		regState{0x01, 0x00, 0x00, 0x8002, 0x00, 0x00}},
	basicTest{
		"OP_LDA_AX",
		[]byte{OP_LDA_AX, 0xFD, 0x7F}, // address is three before OpCode
		regState{x: 0x03},
		regState{OP_LDA_AX, 0x03, 0x00, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDA_AY",
		[]byte{OP_LDA_AY, 0xFD, 0x7F}, // address is three before OpCode
		regState{y: 0x03},
		regState{OP_LDA_AY, 0x00, 0x03, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDA_ZX",
		[]byte{OP_LDA_ZX, 0x03}, // address is three before OpCode
		regState{x: 0x03},
		regState{0x06, 0x03, 0x00, 0x8002, 0x00, 0x00}},

	//These will probably need to be in their own tests
	basicTest{
		"OP_LDA_IX",
		[]byte{
			OP_LDA_IM, 0x00,
			OP_STA_ZP, 0x01,
			OP_LDA_IM, 0x80,
			OP_STA_ZP, 0x02,
			OP_LDA_IX, 0x00,
		},
		regState{x: 1},
		regState{OP_LDA_IM, 1, 0x00, 0x800A, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDA_IY",
		[]byte{OP_LDA_IY, 0x7E}, // pointer should be $7F7E
		regState{y: 130},        // *should* be OP_LDA_IY
		regState{OP_LDA_IY, 0x00, 130, 0x8002, FLAG_NEGATIVE, 0x00}},

	// LDX
	basicTest{
		"OP_LDX_AB",
		[]byte{OP_LDX_AB, 0x00, 0x80},
		regState{},
		regState{0x00, OP_LDX_AB, 0x00, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDX_AY",
		[]byte{OP_LDX_AY, 0xFD, 0x7F},
		regState{y: 0x03},
		regState{0x00, OP_LDX_AY, 0x03, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDX_IM",
		[]byte{OP_LDX_IM, 0x04},
		regState{},
		regState{0x00, 0x04, 0x00, 0x8002, 0x00, 0x00}},
	basicTest{
		"OP_LDX_ZP",
		[]byte{OP_LDX_ZP, 0x04},
		regState{},
		regState{0x00, 0x04, 0x00, 0x8002, 0x00, 0x00}},
	basicTest{
		"OP_LDX_ZY",
		[]byte{OP_LDX_ZY, 0x01},
		regState{y: 0x03},
		regState{0x00, 0x04, 0x03, 0x8002, 0x00, 0x00}},

	// LDY
	basicTest{
		"OP_LDY_IM",
		[]byte{OP_LDY_IM, 0x04},
		regState{},
		regState{0x00, 0x00, 0x04, 0x8002, 0x00, 0x00}},
	basicTest{
		"OP_LDY_AB",
		[]byte{OP_LDY_AB, 0x00, 0x80},
		regState{},
		regState{0x00, 0x00, OP_LDY_AB, 0x8003, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDY_ZP",
		[]byte{OP_LDY_ZP, 0xF4},
		regState{},
		regState{0x00, 0x00, 0xF4, 0x8002, FLAG_NEGATIVE, 0x00}},
	basicTest{
		"OP_LDY_ZX",
		[]byte{OP_LDY_ZX, 0x04},
		regState{x: 0x03},
		regState{x: 0x03, y: 0x07, pc: 0x8002}},
	basicTest{
		"OP_LDY_ZX",
		[]byte{OP_LDY_ZX, 0x01},
		regState{x: 0x03},
		regState{0x00, 0x03, 0x04, 0x8002, 0x00, 0x00}},

	basicTest{
		"OP_NOP",
		[]byte{OP_NOP},
		regState{},
		regState{pc: 0x8001}},

	basicTest{
		"OP_TAX",
		[]byte{OP_TAX},
		regState{a: 4},
		regState{a: 4, x: 4, pc: 0x8001}},
	basicTest{
		"OP_TAY",
		[]byte{OP_TAY},
		regState{a: 4},
		regState{a: 4, y: 4, pc: 0x8001}},
	basicTest{
		"OP_TSX",
		[]byte{OP_TSX},
		regState{stack: 4},
		regState{x: 0x04, pc: 0x8001, stack: 4}},
	basicTest{
		"OP_TXA",
		[]byte{OP_TXA},
		regState{x: 4},
		regState{a: 4, x: 4, pc: 0x8001}},
	basicTest{
		"OP_TXS",
		[]byte{OP_TXS},
		regState{x: 4},
		regState{x: 0x04, pc: 0x8001, stack: 0x04}},
}

type regState struct {
	a      uint8
	x      uint8
	y      uint8
	pc     uint16
	phlags uint8
	stack  uint8
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

func (bt basicTest) GetExpected() regState {
	return bt.regExpected
}

func padToPage(input []byte) []byte {
	for len(input)%256 != 0 {
		input = append(input, 0x00)
	}

	return input
}

func TestBasic(t *testing.T) {
	core := newTestCore(t)
	for _, bt := range testData_A {
		t.Run(bt.name, func(t *testing.T) {
			testsRun++

			err := core.resetTest(t, bt.rom, nil)
			if err != nil {
				t.Errorf("%s: %v", bt.name, err)
			}

			core.setRegisters(t, bt.regInitial)

			ticksran := 0

			//for i := 0; i < bt.ticks; i++ {
			for !core.testDone {
				err = core.tick()
				if err != nil {
					//core.dumpPage(0, t)
					t.Fatalf("%s: %v", bt.name, err)
				}
				ticksran++
				if ticksran > 1000 {
					t.Error("Tick limit hit")
					break
				}
			}
			ticksran -= 1 // remove the 0xFF "test done" tick

			core.checkRegisters(t, bt.name, bt.regExpected)
			if t.Failed() {
				core.dumpPage(0, t)
				core.dumpPage(0x80, t)
				core.dumpReg(t)
			}
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

func (mt memTest) GetExpected() regState {
	return mt.regExpected
}

//func (mt memTest) GetInitialState() {
//	return mt.regInitial
//}

func TestMemory(t *testing.T) {
	core := newTestCore(t)
	for _, mt := range memoryBased {
		t.Run(mt.name, func(t *testing.T) {
			testsRun++

			err := core.resetTest(t, mt.rom, nil)
			if err != nil {
				t.Errorf("%s: %v", mt.name, err)
			}

			core.setRegisters(t, mt.regInitial)
			ticksran := 0

			for !core.testDone {
				err = core.tick()
				if err != nil {
					t.Fatalf("%s: %v", mt.name, err)
				}
				ticksran++
			}
			ticksran -= 1 // remove the 0xFF "test done" tick

			core.checkRegisters(t, mt.name, mt.regExpected)

			if core.ReadByte(mt.mem.addr) != mt.mem.val {
				t.Errorf("%s: Incorrect memory value at $%04X: Exp:$%02X Got:$%02X", mt.name, mt.mem.addr, mt.mem.val, core.ReadByte(mt.mem.addr))
			}

			if t.Failed() {
				core.dumpPage(0x00, t)
				core.dumpPage(0x03, t)
				core.dumpPage(0x80, t)
				core.dumpReg(t)
			}
		})
	}
}

func TestEnd(t *testing.T) {
	t.Logf("Tests run: %d", testsRun)
}

// Core helper functions
func (core *Core) checkRegisters(t *testing.T, name string, e regState) {
	t.Helper()
	if core.A != e.a {
		t.Errorf("%s: Incorrect A: Exp:$%02X Got:$%02X", name, e.a, core.A)
	}

	if core.X != e.x {
		t.Errorf("%s: Incorrect X: Exp:$%02X Got:$%02X", name, e.x, core.X)
	}

	if core.Y != e.y {
		t.Errorf("%s: Incorrect Y: Exp:$%02X Got:$%02X", name, e.y, core.Y)
	}

	if core.PC != e.pc {
		t.Errorf("%s: Incorrect PC: Exp:$%04X Got:$%04X", name, e.pc, core.PC)
	}

	if core.Phlags != e.phlags {
		t.Errorf("%s: Incorrect Phlags: Exp:$%02X Got:$%02X", name, e.phlags, core.Phlags)
	}

	if core.SP != e.stack {
		t.Errorf("%s: Incorrect Stack Pointer: Exp:$%02X Got:$%02X", name, e.stack, core.SP)
	}
}

func (core *Core) setRegisters(t *testing.T, r regState) {
	t.Helper()
	core.A = r.a
	core.X = r.x
	core.Y = r.y
	//core.PC = r.pc
	core.Phlags = r.phlags
	core.SP = r.stack
}

func (c *Core) resetTest(t *testing.T, rom, ram []byte) error {
	t.Helper()
	rom = PadWithVectors(rom, 0x8000, 0x8000, 0x8000)
	if len(rom)%256 != 0 {
		return fmt.Errorf("ROM is not divisible by 256: %d", len(rom))
	}

	c.rom = rom
	c.PC = c.ReadWord(VECTOR_RESET)

	if ram != nil {
		c.memory = ram
	} else {
		c.memory = make([]byte, 0x1000)
	}

	c.testDone = false

	// fill zero page with some data
	for i := 0; i < 256; i++ {
		c.memory[i] = uint8(i)
	}

	return nil
}

func (c *Core) dumpReg(t *testing.T) {
	t.Logf("A: %02X (%d) X: %02X (%d) Y: %02X (%d) PC: %04X Phlags: %08b",
		c.A,
		c.A,
		c.X,
		c.X,
		c.Y,
		c.Y,
		c.PC,
		c.Phlags,
	)
}

func (c *Core) dumpPage(page uint8, t *testing.T) {
	vals := []string{}
	base := uint16(page) << 8
	for i := uint16(0); i < 256; i++ {
		vals = append(vals, fmt.Sprintf("%02X", c.ReadByte(base+i)))
	}

	for i := 0; i < 256; i += 16 {
		t.Logf("%04X: %s", int(base)+i, strings.Join(vals[i:i+16], " "))
	}
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
