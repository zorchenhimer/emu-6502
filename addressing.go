package emu

import (
	"fmt"
)

type AddressModeMeta struct {
	Name    string
	Asm     func(c *Core, oppc uint16) string
	Address func(c *Core) (uint16, uint8)
}

var ADDR_Accumulator = AddressModeMeta{
	Name: "Accumulator",
	Asm: func(c *Core, oppc uint16) string {
		return "A"
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.PC, 1
	},
}

var ADDR_Absolute = AddressModeMeta{
	Name: "Absolute",
	Asm: func(c *Core, oppc uint16) string {
		return fmt.Sprintf("$%04X", c.ReadWord(oppc+1))
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(c.PC + 1), 3
	},
}

var ADDR_AbsoluteX = AddressModeMeta{
	Name: "Absolute, X",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadWord(oppc + 1)
		return fmt.Sprintf("$%04X, X @ $%04X",
			value,
			value+uint16(c.X),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(c.PC+1) + uint16(c.X), 3
	},
}

var ADDR_AbsoluteY = AddressModeMeta{
	Name: "Absolute, Y",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadWord(oppc + 1)
		return fmt.Sprintf("$%04X, Y @ $%04X",
			value,
			value+uint16(c.Y),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(c.PC+1) + uint16(c.Y), 3
	},
}

var ADDR_Immediate = AddressModeMeta{
	Name: "#Immediate",
	Asm: func(c *Core, oppc uint16) string {
		return fmt.Sprintf("#$%02X", c.ReadByte(oppc+1))
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.PC + 1, 2
	},
}

var ADDR_Implied = AddressModeMeta{
	Name: "Implied",
	Asm: func(c *Core, oppc uint16) string {
		return ""
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.PC, 1
	},
}

var ADDR_Indirect = AddressModeMeta{
	Name: "(Indirect)",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadWord(oppc + 1)
		return fmt.Sprintf("($%04X) @ $%04X",
			value,
			c.ReadWord(value),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(c.ReadWord(c.PC + 1)), 3
	},
}

var ADDR_IndirectX = AddressModeMeta{
	Name: "(Indirect), X",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadByte(oppc + 1)
		return fmt.Sprintf("($%02X), X @ $%04X",
			value,
			c.ReadWord(uint16(value+c.X)),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(uint16(c.ReadByte(c.PC+1) + c.X)), 2
	},
}

var ADDR_IndirectY = AddressModeMeta{
	Name: "(Indirect, Y)",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadByte(oppc + 1)
		return fmt.Sprintf("($%02X, Y) @ $%04X",
			value,
			c.ReadWord(uint16(value))+uint16(c.Y),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return c.ReadWord(uint16(c.ReadByte(uint16(c.PC+1)))) + uint16(c.Y), 2
	},
}

var ADDR_ZeroPage = AddressModeMeta{
	Name: "ZeroPage",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadByte(oppc + 1)
		return fmt.Sprintf("$%02X = %02X", value, c.ReadByte(uint16(value)))
	},
	Address: func(c *Core) (uint16, uint8) {
		return uint16(c.ReadByte(c.PC + 1)), 2
	},
}

var ADDR_ZeroPageX = AddressModeMeta{
	Name: "ZeroPage, X",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadByte(oppc + 1)
		return fmt.Sprintf("$%02X, X   @ $%04X",
			value,
			uint16(value+c.X),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return uint16(c.ReadByte(c.PC+1) + c.X), 2
	},
}

var ADDR_ZeroPageY = AddressModeMeta{
	Name: "ZeroPage, Y",
	Asm: func(c *Core, oppc uint16) string {
		value := c.ReadByte(oppc + 1)
		return fmt.Sprintf("$%02X, Y   @ $%04X",
			value,
			uint16(value+c.Y),
		)
	},
	Address: func(c *Core) (uint16, uint8) {
		return uint16(c.ReadByte(c.PC+1) + c.Y), 2
	},
}

var ADDR_Relative = AddressModeMeta{
	Name: "Relative",
	Asm: func(c *Core, oppc uint16) string {
		value := c.addrRelative(oppc, c.ReadByte(oppc+1))
		n, neg := TwosCompInv(c.ReadByte(oppc + 1))
		num := int(n)
		if neg {
			num *= -1
		}
		return fmt.Sprintf("$%02X   (%d)", value, num)
	},
	Address: func(c *Core) (uint16, uint8) {
		panic("branch Address()")
		return c.addrRelative(c.PC, c.ReadByte(c.PC+1)), 2
	},
}
