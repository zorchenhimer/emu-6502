package emu

import (
	"fmt"
)

type AddressModeMeta struct {
	Name string
	Asm func(value uint16) string
	Address func(c *Core) (uint16, uint8)
}

var ADDR_Absolute = AddressModeMeta{
		Name: "Absolute",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%04X", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(c.PC + 1), 3
		},
	}

var ADDR_AbsoluteX = AddressModeMeta{
		Name: "Absolute, X",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%04X, X", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(c.PC + 1) + uint16(c.X), 3
		},
	}

var ADDR_AbsoluteY = AddressModeMeta{
		Name: "Absolute, Y",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%04X, Y", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(c.PC + 1) + uint16(c.Y), 3
		},
	}

var ADDR_Immediate = AddressModeMeta{
		Name: "#Immediate",
		Asm: func(value uint16) string {
			return fmt.Sprintf("#$%02X", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.PC + 1, 2
		},
	}

var ADDR_Implied = AddressModeMeta{
		Name: "Implied",
		Asm: func(value uint16) string {
			return ""
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.PC, 1
		},
	}

var ADDR_Indirect = AddressModeMeta{
		Name: "(Indirect)",
		Asm: func(value uint16) string {
			return fmt.Sprintf("($%04X)", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(c.ReadWord(c.PC + 1)), 3
		},
	}

var ADDR_IndirectX = AddressModeMeta{
		Name: "(Indirect), X",
		Asm: func(value uint16) string {
			return fmt.Sprintf("($%04X), X", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(uint16(c.ReadByte(c.PC + 1) + c.X)), 2
		},
	}

var ADDR_IndirectY = AddressModeMeta{
		Name: "(Indirect, Y)",
		Asm: func(value uint16) string {
			return fmt.Sprintf("($%04X, Y)", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.ReadWord(uint16(c.ReadByte(uint16(c.PC + 1)))) + uint16(c.Y), 2
		},
	}

var ADDR_ZeroPage = AddressModeMeta{
		Name: "ZeroPage",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%02X", uint8(value))
		},
		Address: func(c *Core) (uint16, uint8) {
			return uint16(c.ReadByte(c.PC + 1)), 2
		},
	}

var ADDR_ZeroPageX = AddressModeMeta{
		Name: "ZeroPage, X",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%02X, X", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return uint16(c.ReadByte(c.PC + 1) + c.X), 2
		},
	}

var ADDR_ZeroPageY = AddressModeMeta{
		Name: "ZeroPage, Y",
		Asm: func(value uint16) string {
			return fmt.Sprintf("$%02X, Y", value)
		},
		Address: func(c *Core) (uint16, uint8) {
			return uint16(c.ReadByte(c.PC + 1) + c.Y), 2
		},
	}

var ADDR_Relative = AddressModeMeta{
		Name: "Relative",
		Asm: func(value uint16) string {
			n, neg := TwosCompInv(uint8(value))
			num := int(n)
			if neg {
				num *= -1
			}
			return fmt.Sprintf("$%02X   (%d)", value, num)
		},
		Address: func(c *Core) (uint16, uint8) {
			return c.addrRelative(c.ReadByte(c.PC+1)), 2
		},
}

