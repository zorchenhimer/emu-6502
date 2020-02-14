// +build ignore

package emu

import (
	"fmt"
)

/*
	Addressing functions return the absolute destination address
	to write to.

	The input to these functions is the address following the
	instruction.
*/

func (c *Core) addrAbsoluteX(addr uint16) uint16 {
	return c.ReadWord(addr) + uint16(c.X)
}

func (c *Core) addrAbsoluteY(addr uint16) uint16 {
	return c.ReadWord(addr) + uint16(c.Y)
}

func (c *Core) addrIndirectY(addr uint16) uint16 {
	return c.ReadWord(uint16(c.ReadByte(uint16(addr)))) + uint16(c.Y)
}

func (c *Core) addrIndirectX(addr uint16) uint16 {
	return c.ReadWord(uint16(c.ReadByte(addr) + c.X))
}

func (c *Core) addrZeroPage(addr uint16) uint16 {
	return uint16(c.ReadByte(addr))
}

func (c *Core) addrZeroPageX(addr uint16) uint16 {
	return uint16(c.ReadByte(addr) + c.X)
}

func (c *Core) addrZeroPageY(addr uint16) uint16 {
	return uint16(c.ReadByte(addr) + c.Y)
}

func exec_CLD(c *Core) {
	// lel
	c.PC += 1
}

func (c *Core) branch(flag uint8, set bool) {
	var v uint8 = 0
	if set {
		v = 1
	}

	prevPc := c.PC
	if c.Phlags&flag == v {
		c.PC = c.addrRelative(c.ReadByte(c.PC + 1))
	} else {
		c.PC += 2
	}

	fmt.Printf("branch: %s set: %t [%04X] -> [%04X]\n",
		flagToString(flag),
		set,
		prevPc,
		c.PC,
	)
	c.DumpRegisters()
}

/* Branches */
func exec_BCC(c *Core) {
	c.branch(FLAG_CARRY, false)
}

func exec_BCS(c *Core) {
	c.branch(FLAG_CARRY, true)
}

func exec_BEQ(c *Core) {
	c.branch(FLAG_ZERO, true)
}

func exec_BNE(c *Core) {
	c.branch(FLAG_ZERO, false)
}

func exec_BMI(c *Core) {
	c.branch(FLAG_NEGATIVE, true)
}

func exec_BPL(c *Core) {
	c.branch(FLAG_NEGATIVE, false)
}

func exec_BVC(c *Core) {
	c.branch(FLAG_OVERFLOW, false)
}

func exec_BVS(c *Core) {
	c.branch(FLAG_OVERFLOW, true)
}

/* Decrements */
func exec_DEC_AB(c *Core) {
	val := c.ReadByte(c.ReadWord(c.PC + 1))
	val -= 1
	c.WriteByte(c.ReadWord(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 3
}

func exec_DEC_AX(c *Core) {
	val := c.ReadByte(c.addrAbsoluteX(c.PC + 1))
	val -= 1
	c.WriteByte(c.addrAbsoluteX(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 3
}

func exec_DEC_ZP(c *Core) {
	val := c.ReadByte(c.addrZeroPage(c.PC + 1))
	val -= 1
	c.WriteByte(c.addrZeroPage(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 2
}

func exec_DEC_ZX(c *Core) {
	val := c.ReadByte(c.addrZeroPageX(c.PC + 1))
	val -= 1
	c.WriteByte(c.addrZeroPageX(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 2
}

func exec_DEX(c *Core) {
	c.X -= 1
	c.setZeroNegative(c.X)
	c.PC += 1
}

func exec_DEY(c *Core) {
	c.Y -= 1
	c.setZeroNegative(c.Y)
	c.PC += 1
}

/* Increments */
func exec_INC_AB(c *Core) {
	val := c.ReadByte(c.ReadWord(c.PC + 1))
	val += 1
	c.WriteByte(c.ReadWord(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 3
}

func exec_INC_AX(c *Core) {
	val := c.ReadByte(c.addrAbsoluteX(c.PC + 1))
	val += 1
	c.WriteByte(c.addrAbsoluteX(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 3
}

func exec_INC_ZP(c *Core) {
	val := c.ReadByte(c.addrZeroPage(c.PC + 1))
	val += 1
	c.WriteByte(c.addrZeroPage(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 2
}

func exec_INC_ZX(c *Core) {
	val := c.ReadByte(c.addrZeroPageX(c.PC + 1))
	val += 1
	c.WriteByte(c.addrZeroPageX(c.PC+1), val)
	c.setZeroNegative(val)
	c.PC += 2
}

func exec_INX(c *Core) {
	c.X += 1
	c.setZeroNegative(c.X)
	c.PC += 1
}

func exec_INY(c *Core) {
	c.Y += 1
	c.setZeroNegative(c.Y)
	c.PC += 1
}

/* JMP */
func exec_JMP_AB(c *Core) {
	c.PC = c.ReadWord(c.PC + 1)
}

func exec_JMP_ID(c *Core) {
	c.PC = c.ReadWord(c.ReadWord(c.PC + 1))
}

/* LDA */
func exec_LDA_IM(c *Core) {
	c.A = c.ReadByte(c.PC + 1)
	c.setZeroNegative(c.A)
	c.PC += 2
}

func exec_LDA_AB(c *Core) {
	c.A = c.ReadByte(c.ReadWord(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 3
}

func exec_LDA_AX(c *Core) {
	c.A = c.ReadByte(c.addrAbsoluteX(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 3
}

func exec_LDA_AY(c *Core) {
	c.A = c.ReadByte(c.addrAbsoluteY(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 3
}

func exec_LDA_IX(c *Core) {
	c.A = c.ReadByte(c.addrIndirectX(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 2
}

func exec_LDA_IY(c *Core) {
	c.A = c.ReadByte(c.addrIndirectY(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 2
}

func exec_LDA_ZP(c *Core) {
	c.A = c.ReadByte(c.addrZeroPage(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 2
}

func exec_LDA_ZX(c *Core) {
	c.A = c.ReadByte(c.addrZeroPageX(c.PC + 1))
	c.setZeroNegative(c.A)
	c.PC += 2
}

/* LDX */
func exec_LDX_AB(c *Core) {
	c.X = c.ReadByte(c.ReadWord(c.PC + 1))
	c.setZeroNegative(c.X)
	c.PC += 3
}

func exec_LDX_AY(c *Core) {
	c.X = c.ReadByte(c.addrAbsoluteY(c.PC + 1))
	c.setZeroNegative(c.X)
	c.PC += 3
}

func exec_LDX_IM(c *Core) {
	c.X = c.ReadByte(c.PC + 1)
	c.setZeroNegative(c.X)
	c.PC += 2
}

func exec_LDX_ZP(c *Core) {
	c.X = c.ReadByte(c.addrZeroPage(c.PC + 1))
	c.setZeroNegative(c.X)
	c.PC += 2
}

func exec_LDX_ZY(c *Core) {
	c.X = c.ReadByte(c.addrZeroPageY(c.PC + 1))
	c.setZeroNegative(c.X)
	c.PC += 2
}

/* LDY */
func exec_LDY_AB(c *Core) {
	c.Y = c.ReadByte(c.ReadWord(c.PC + 1))
	c.setZeroNegative(c.Y)
	c.PC += 3
}

func exec_LDY_AX(c *Core) {
	c.Y = c.ReadByte(c.addrAbsoluteX(c.PC + 1))
	c.setZeroNegative(c.Y)
	c.PC += 3
}

func exec_LDY_IM(c *Core) {
	c.Y = c.ReadByte(c.PC + 1)
	c.setZeroNegative(c.Y)
	c.PC += 2
}

func exec_LDY_ZP(c *Core) {
	c.Y = c.ReadByte(c.addrZeroPage(c.PC + 1))
	c.setZeroNegative(c.Y)
	c.PC += 2
}

func exec_LDY_ZX(c *Core) {
	c.Y = c.ReadByte(c.addrZeroPageX(c.PC + 1))
	c.setZeroNegative(c.Y)
	c.PC += 2
}

func exec_NOP(c *Core) {
	c.PC += 1
}

/* STA */
func exec_STA_AB(c *Core) {
	c.WriteByte(c.ReadWord(c.PC+1), c.A)
	c.PC += 3
}

func exec_STA_AX(c *Core) {
	c.WriteByte(c.addrAbsoluteX(c.PC+1), c.A)
	c.PC += 3
}

func exec_STA_IX(c *Core) {
	c.WriteByte(c.addrIndirectX(c.PC+1), c.A)
	c.PC += 2
}

func exec_STA_IY(c *Core) {
	c.WriteByte(c.addrIndirectY(c.PC+1), c.A)
	c.PC += 2
}

func exec_STA_AY(c *Core) {
	c.WriteByte(c.addrAbsoluteY(c.PC+1), c.A)
	c.PC += 3
}

func exec_STA_ZP(c *Core) {
	c.WriteByte(c.addrZeroPage(c.PC+1), c.A)
	c.PC += 2
}

func exec_STA_ZX(c *Core) {
	c.WriteByte(c.addrZeroPageX(c.PC+1), c.A)
	c.PC += 2
}

/* STX */
func exec_STX_AB(c *Core) {
	c.WriteByte(c.ReadWord(c.PC+1), c.X)
	c.PC += 3
}
func exec_STX_ZP(c *Core) {
	c.WriteByte(c.addrZeroPage(c.PC+1), c.X)
	c.PC += 2
}

func exec_STX_ZY(c *Core) {
	c.WriteByte(c.addrZeroPageY(c.PC+1), c.X)
	c.PC += 2
}

/* transfers */
func exec_TAX(c *Core) {
	c.X = c.A
	c.setZeroNegative(c.X)
	c.PC += 1
}

func exec_TAY(c *Core) {
	c.Y = c.A
	c.setZeroNegative(c.Y)
	c.PC += 1
}

func exec_TSX(c *Core) {
	c.X = c.SP
	c.setZeroNegative(c.X)
	c.PC += 1
}

func exec_TXA(c *Core) {
	c.A = c.X
	c.setZeroNegative(c.A)
	c.PC += 1
}

func exec_TXS(c *Core) {
	c.SP = c.X
	c.PC += 1
}

func exec_TYA(c *Core) {
	c.A = c.Y
	c.setZeroNegative(c.A)
	c.PC += 1
}
