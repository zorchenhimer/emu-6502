package emu

import (
	"fmt"
)

func exec_BRK(c *Core) error {
	return fmt.Errorf("BRK Unimplimented")
}

/* JMP */
func exec_JMP_AB(c *Core) error {
	c.PC = c.ReadWord(c.PC + 1)
	return nil
}

func exec_JMP_ID(c *Core) error {
	c.PC = c.ReadWord(c.ReadWord(c.PC + 1))
	return nil
}

/* LDA */
func exec_LDA_IM(c *Core) error {
	c.A = c.ReadByte(c.PC + 1)
	c.PC += 2
	return nil
}

func exec_LDA_AB(c *Core) error {
	c.A = c.ReadByte(c.ReadWord(c.PC + 1))
	c.PC += 3
	return nil
}

func exec_LDA_AX(c *Core) error {
	c.A = c.ReadByte(c.ReadWord(c.PC+1) + uint16(c.X))
	c.PC += 3
	return nil
}

func exec_LDA_AY(c *Core) error {
	c.A = c.ReadByte(c.ReadWord(c.PC+1) + uint16(c.Y))
	c.PC += 3
	return nil
}

func exec_LDA_IX(c *Core) error {
	return fmt.Errorf("OP_ADC_AX not implemented")
}

func exec_LDA_IY(c *Core) error {
	return fmt.Errorf("OP_ADC_AY not implemented")
}

func exec_LDA_ZP(c *Core) error {
	c.A = c.ReadByte(uint16(c.ReadByte(c.PC + 1)))
	c.PC += 2
	return nil
}

func exec_LDA_ZX(c *Core) error {
	return fmt.Errorf("OP_ADC_ZX not implemented")
}

/* LDX */
func exec_LDX_AB(c *Core) error {
	c.X = c.ReadByte(c.ReadWord(c.PC + 1))
	c.PC += 3
	return nil
}

func exec_LDX_AY(c *Core) error {
	c.X = c.ReadByte(c.ReadWord(c.PC+1) + uint16(c.Y))
	c.PC += 3
	return nil
}

func exec_LDX_IM(c *Core) error {
	c.X = c.ReadByte(c.PC + 1)
	c.PC += 2
	return nil
}

func exec_LDX_ZP(c *Core) error {
	c.X = c.ReadByte(uint16(c.ReadByte(c.PC + 1)))
	c.PC += 2
	return nil
}

/* LDY */
func exec_LDY_AB(c *Core) error {
	c.Y = c.ReadByte(c.ReadWord(c.PC + 1))
	c.PC += 3
	return nil
}

func exec_LDY_AX(c *Core) error {
	c.Y = c.ReadByte(c.ReadWord(c.PC+1) + uint16(c.X))
	c.PC += 3
	return nil
}

func exec_LDY_IM(c *Core) error {
	c.Y = c.ReadByte(c.PC + 1)
	c.PC += 2
	return nil
}

func exec_LDY_ZP(c *Core) error {
	c.Y = c.ReadByte(uint16(c.ReadByte(c.PC + 1)))
	c.PC += 2
	return nil
}

/* ADC */
func exec_ADC_AB(c *Core) error {
	return fmt.Errorf("OP_ADC_AB not implemented")
}

func exec_ADC_AX(c *Core) error {
	return fmt.Errorf("OP_ADC_AX not implemented")
}

func exec_ADC_AY(c *Core) error {
	return fmt.Errorf("OP_ADC_AY not implemented")
}

func exec_ADC_IM(c *Core) error {
	return fmt.Errorf("OP_ADC_IM not implemented")
}

func exec_ADC_IX(c *Core) error {
	return fmt.Errorf("OP_ADC_IX not implemented")
}

func exec_ADC_IY(c *Core) error {
	return fmt.Errorf("OP_ADC_IY not implemented")
}

func exec_ADC_ZP(c *Core) error {
	return fmt.Errorf("OP_ADC_ZP not implemented")
}

func exec_ADC_ZX(c *Core) error {
	return fmt.Errorf("OP_ADC_ZX not implemented")
}

func exec_NOP(c *Core) error {
	c.PC += 1
	return nil
}
