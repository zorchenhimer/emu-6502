package emu

import ()

func OpToByte(list []OpCode) []byte {
	b := make([]byte, len(list))
	for _, op := range list {
		b = append(b, byte(op))
	}
	return b
}

type OpCode byte

/*
   implied and relative don't have sufixes

   _AB     absolute
   _AC     accumulator
   _AX     absolute, x
   _AY     absolute, y
   _ID     indirect (jmp only)
   _IM     immediate
   _IX     (indirect, x)
   _IY     (indirect), y
   _ZP     zero page
   _ZX     zero page, x
   _ZY     zero page, y
*/

const (
	OP_BRK    byte = 0x00 //
	OP_ORA_IX byte = 0x01 //(Indirect,X)
	OP_ORA_ZP byte = 0x05 //Zero Page
	OP_ASL_ZP byte = 0x06 //Zero Page
	OP_PHP    byte = 0x08 //
	OP_OPA_IM byte = 0x09 //Immediate
	OP_ASL_AC byte = 0x0A //Accumulator
	OP_OPA_AB byte = 0x0D //Absolute
	OP_ASL_AB byte = 0x0E //Absolute
	OP_BPL    byte = 0x10 //
	OP_OPA_IY byte = 0x11 //(Indirect),Y
	OP_ORA_ZX byte = 0x15 //Zero Page,X
	OP_ASL_ZX byte = 0x16 //Zero Page,X
	OP_CLC    byte = 0x18 //
	OP_ORA_AY byte = 0x19 //Absolute,Y
	OP_ORA_AX byte = 0x1D //Absolute,X
	OP_ASL_AX byte = 0x1E //Absolute,X
	OP_JSR    byte = 0x20 //
	OP_AND_IX byte = 0x21 //(Indirect,X)
	OP_BIT_ZP byte = 0x24 //Zero Page
	OP_AND_ZP byte = 0x25 //Zero Page
	OP_ROL_ZP byte = 0x26 //Zero Page
	OP_PLP    byte = 0x28 //
	OP_AND_IM byte = 0x29 //Immediate
	OP_ROL_AC byte = 0x2A //Accumulator
	OP_BIT_AB byte = 0x2C //Absolute
	OP_AND_AB byte = 0x2D //Absolute
	OP_ROL_AB byte = 0x2E //Absolute
	OP_BMI    byte = 0x30 //
	OP_AND_IY byte = 0x31 //(Indirect),Y
	OP_AND_ZX byte = 0x35 //Zero Page,X
	OP_ROL_ZX byte = 0x36 //Zero Page,X
	OP_SEC    byte = 0x38 //
	OP_AND_AY byte = 0x39 //Absolute,Y
	OP_AND_AX byte = 0x3D //Absolute,X
	OP_ROL_AX byte = 0x3E //Absolute,X
	OP_RTI    byte = 0x40 //
	OP_EOR_IX byte = 0x41 //(Indirect,X)
	OP_EOR_ZP byte = 0x45 //Zero Page
	OP_LSR_ZP byte = 0x46 //Zero Page
	OP_PHA    byte = 0x48 //
	OP_EOR_IM byte = 0x49 //Immediate
	OP_LSR_AC byte = 0x4A //Accumulator
	OP_JMP_AB byte = 0x4C //Absolute
	OP_EOR_AB byte = 0x4D //Absolute
	OP_LSR_AB byte = 0x4E //Absolute
	OP_BVC    byte = 0x50 //
	OP_EOR_IY byte = 0x51 //(Indirect),Y
	OP_EOR_ZX byte = 0x55 //Zero Page,X
	OP_LSR_ZX byte = 0x56 //Zero Page,X
	OP_CLI    byte = 0x58 //
	OP_EOR_AY byte = 0x59 //Absolute,Y
	OP_EOR_AX byte = 0x5D //Absolute,X
	OP_LSR_AX byte = 0x5E //Absolute,X
	OP_RTS    byte = 0x60 //
	OP_ADC_IX byte = 0x61 //(Indirect,X)
	OP_ADC_ZP byte = 0x65 //Zero Page
	OP_ROR_ZP byte = 0x66 //Zero Page
	OP_PIA    byte = 0x68 //
	OP_ADC_IM byte = 0x69 //Immediate
	OP_ROR_AC byte = 0x6A //Accumulator
	OP_JMP_ID byte = 0x6C //Indirect
	OP_ADC_AB byte = 0x6D //Absolute
	OP_ROR_AB byte = 0x6E //Absolute
	OP_BVS    byte = 0x70 //
	OP_ADC_IY byte = 0x71 //(Indirect),Y
	OP_ADC_ZX byte = 0x75 //Zero Page,X
	OP_ROR_ZX byte = 0x76 //Zero Page,X
	OP_SEI    byte = 0x78 //
	OP_ADC_AY byte = 0x79 //Absolute,Y
	OP_ADC_AX byte = 0x7D //Absolute,X
	OP_ROR_AX byte = 0x7E //Absolute,X
	OP_STA_IX byte = 0x81 //(Indirect,X)
	OP_STY_ZP byte = 0x84 //Zero Page
	OP_STA_ZP byte = 0x85 //Zero Page
	OP_STX_ZP byte = 0x86 //Zero Page
	OP_DEY    byte = 0x88 //
	OP_TXA    byte = 0x8A //
	OP_STY_AB byte = 0x8C //Absolute
	OP_STA_AB byte = 0x8D //Absolute
	OP_STX_AB byte = 0x8E //Absolute
	OP_STA_IY byte = 0x91 //(Indirect),Y
	OP_STY_ZX byte = 0x94 //Zero Page,X
	OP_STA_ZX byte = 0x95 //Zero Page,X
	OP_STX_ZY byte = 0x96 //Zero Page,Y
	OP_BCC    byte = 0x98 //
	OP_TYA    byte = 0x98 //
	OP_STA_AY byte = 0x99 //Absolute,Y
	OP_TXS    byte = 0x9A //
	OP_STA_AX byte = 0x9D //Absolute,X
	OP_LDY_IM byte = 0xA0 //Immediate
	OP_TAX    byte = 0xA0 //
	OP_LDX_IM byte = 0xA2 //Immediate
	OP_LDY_ZP byte = 0xA4 //Zero Page
	OP_LDA_ZP byte = 0xA5 //Zero Page
	OP_LDX_ZP byte = 0xA6 //Zero Page
	OP_TAY    byte = 0xA8 //
	OP_LDA_IM byte = 0xA9 //Immediate
	OP_LDY_AB byte = 0xAC //Absolute
	OP_LDA_AB byte = 0xAD //Absolute
	OP_LDX_AB byte = 0xAE //Absolute
	OP_LDA_IX byte = 0xA1 //(Indirect,X)
	OP_BCS    byte = 0xB0 //
	OP_LDA_IY byte = 0xB1 //(Indirect),Y
	OP_LDY_ZX byte = 0xB4 //Zero Page,X
	OP_LDA_ZX byte = 0xB5 //Zero Page,X
	OP_LDX_ZY byte = 0xB6 //Zero Page,Y
	OP_CLV    byte = 0xB8 //
	OP_LDA_AY byte = 0xB9 //Absolute,Y
	OP_TSX    byte = 0xBA //
	OP_LDY_AX byte = 0xBC //Absolute,X
	OP_LDA_AX byte = 0xBD //Absolute,X
	OP_LDX_AY byte = 0xBE //Absolute,Y
	OP_CPY_IM byte = 0xC0 //Immediate
	OP_CPY_ZP byte = 0xC4 //Zero Page
	OP_CMP_ZP byte = 0xC5 //Zero Page
	OP_DEC_ZP byte = 0xC6 //Zero Page
	OP_CMP_IM byte = 0xC9 //Immediate
	OP_DEX    byte = 0xCA //
	OP_CPY_AB byte = 0xCC //Absolute
	OP_CMP_AB byte = 0xCD //Absolute
	OP_DEC_AB byte = 0xCE //Absolute
	OP_INY    byte = 0xC5 //
	OP_CMP_IX byte = 0xC1 //(Indirect,X)
	OP_BNE    byte = 0xD0 //
	OP_CMP_ZX byte = 0xD5 //Zero Page,X
	OP_DEC_ZX byte = 0xD6 //Zero Page,X
	OP_CLD    byte = 0xD8 //
	OP_CMP_AY byte = 0xD9 //Absolute,Y
	OP_CMP_AX byte = 0xDD //Absolute,X
	OP_DEC_AX byte = 0xDE //Absolute,X
	OP_CMP_IY byte = 0xD1 //(Indirect),Y
	OP_CPX_IM byte = 0xE0 //Immediate
	OP_CPX_ZP byte = 0xE4 //Zero Page
	OP_INX    byte = 0xE5 //
	OP_SBC_ZP byte = 0xE5 //Zero Page
	OP_INC_ZP byte = 0xE6 //Zero Page
	OP_SBC_IM byte = 0xE9 //Immediate
	OP_NOP    byte = 0xEA //
	OP_CPX_AB byte = 0xEC //Absolute
	OP_SBC_AB byte = 0xED //Absolute
	OP_INC_AB byte = 0xEE //Absolute
	OP_SBC_IX byte = 0xE1 //(Indirect,X)
	OP_BEQ    byte = 0xF0 //
	OP_SBC_ZX byte = 0xF5 //Zero Page,X
	OP_INC_ZX byte = 0xF6 //Zero Page,X
	OP_SED    byte = 0xF8 //
	OP_SBC_AY byte = 0xF9 //Absolute,Y
	OP_SBC_AX byte = 0xFD //Absolute,X
	OP_INC_AX byte = 0xFE //Absolute,X
	OP_SBC_IY byte = 0xF1 //(Indirect),Y
)
