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

type opcodeExec func(core *Core)

//var opcodes map[byte]opcodeExec = map[byte]opcodeExec{
//	//OP_ADC_AB: exec_ADC_AB,
//	//OP_ADC_AX: exec_ADC_AX,
//	//OP_ADC_AY: exec_ADC_AY,
//	//OP_ADC_IM: exec_ADC_IM,
//	//OP_ADC_IX: exec_ADC_IX,
//	//OP_ADC_IY: exec_ADC_IY,
//	//OP_ADC_ZP: exec_ADC_ZP,
//	//OP_ADC_ZX: exec_ADC_ZX,
//
//	//OP_AND_AB: exec_AND_AB,
//	//OP_AND_AX: exec_AND_AX,
//	//OP_AND_AY: exec_AND_AY,
//	//OP_AND_IM: exec_AND_IM,
//	//OP_AND_IX: exec_AND_IX,
//	//OP_AND_IY: exec_AND_IY,
//	//OP_AND_ZP: exec_AND_ZP,
//	//OP_AND_ZX: exec_AND_ZX,
//
//	//OP_ASL_AB: exec_ASL_AB,
//	//OP_ASL_AC: exec_ASL_AC,
//	//OP_ASL_AX: exec_ASL_AX,
//	//OP_ASL_ZP: exec_ASL_ZP,
//	//OP_ASL_ZX: exec_ASL_ZX,
//
//	OP_BCC: exec_BCC,
//	OP_BCS: exec_BCS,
//	OP_BEQ: exec_BEQ,
//	OP_BMI: exec_BMI,
//	OP_BNE: exec_BNE,
//	OP_BPL: exec_BPL,
//	OP_BVC: exec_BVC,
//	OP_BVS: exec_BVS,
//
//	//OP_BRK: exec_BRK,
//
//	//OP_BIT_AB: exec_BIT_AB,
//	//OP_BIT_ZP: exec_BIT_ZP,
//
//	//OP_CLC: exec_CLC,
//	OP_CLD: exec_CLD,
//	//OP_CLI: exec_CLI,
//	//OP_CLV: exec_CLV,
//
//	//OP_CMP_AB: exec_CMP_AB,
//	//OP_CMP_AX: exec_CMP_AX,
//	//OP_CMP_AY: exec_CMP_AY,
//	//OP_CMP_IM: exec_CMP_IM,
//	//OP_CMP_IX: exec_CMP_IX,
//	//OP_CMP_IY: exec_CMP_IY,
//	//OP_CMP_ZP: exec_CMP_ZP,
//	//OP_CMP_ZX: exec_CMP_ZX,
//
//	//OP_CPX_AB: exec_CPX_AB,
//	//OP_CPX_IM: exec_CPX_IM,
//	//OP_CPX_ZP: exec_CPX_ZP,
//	//OP_CPY_AB: exec_CPY_AB,
//	//OP_CPY_IM: exec_CPY_IM,
//	//OP_CPY_ZP: exec_CPY_ZP,
//
//	OP_DEC_AB: exec_DEC_AB,
//	OP_DEC_AX: exec_DEC_AX,
//	OP_DEC_ZP: exec_DEC_ZP,
//	OP_DEC_ZX: exec_DEC_ZX,
//
//	OP_DEX: exec_DEX,
//	OP_DEY: exec_DEY,
//
//	//OP_EOR_AB: exec_EOR_AB,
//	//OP_EOR_AX: exec_EOR_AX,
//	//OP_EOR_AY: exec_EOR_AY,
//	//OP_EOR_IX: exec_EOR_IX,
//	//OP_EOR_IY: exec_EOR_IY,
//	//OP_EOR_ZP: exec_EOR_ZP,
//	//OP_EOR_ZX: exec_EOR_ZX,
//
//	OP_INC_AB: exec_INC_AB,
//	OP_INC_AX: exec_INC_AX,
//	OP_INC_ZP: exec_INC_ZP,
//	OP_INC_ZX: exec_INC_ZX,
//
//	OP_INX: exec_INX,
//	OP_INY: exec_INY,
//
//	OP_JMP_AB: exec_JMP_AB,
//	OP_JMP_ID: exec_JMP_ID,
//	//OP_JSR: exec_JSR,
//
//	OP_LDA_AB: exec_LDA_AB,
//	OP_LDA_AX: exec_LDA_AX,
//	OP_LDA_AY: exec_LDA_AY,
//	OP_LDA_IM: exec_LDA_IM,
//	OP_LDA_IX: exec_LDA_IX,
//	OP_LDA_IY: exec_LDA_IY,
//	OP_LDA_ZP: exec_LDA_ZP,
//	OP_LDA_ZX: exec_LDA_ZX,
//
//	OP_LDX_AB: exec_LDX_AB,
//	OP_LDX_AY: exec_LDX_AY,
//	OP_LDX_IM: exec_LDX_IM,
//	OP_LDX_ZP: exec_LDX_ZP,
//	OP_LDX_ZY: exec_LDX_ZY,
//
//	OP_LDY_AB: exec_LDY_AB,
//	OP_LDY_AX: exec_LDY_AX,
//	OP_LDY_IM: exec_LDY_IM,
//	OP_LDY_ZP: exec_LDY_ZP,
//	OP_LDY_ZX: exec_LDY_ZX,
//
//	//OP_LSR_AB: exec_LSR_AB,
//	//OP_LSR_AC: exec_LSR_AC,
//	//OP_LSR_AX: exec_LSR_AX,
//	//OP_LSR_ZP: exec_LSR_ZP,
//	//OP_LSR_ZX: exec_LSR_ZX,
//
//	OP_NOP: exec_NOP,
//
//	//OP_ORA_AB: exec_OPA_AB,
//	//OP_ORA_AX: exec_ORA_AX,
//	//OP_ORA_AY: exec_ORA_AY,
//	//OP_ORA_IM: exec_OPA_IM,
//	//OP_ORA_IX: exec_ORA_IX,
//	//OP_ORA_IY: exec_OPA_IY,
//	//OP_ORA_ZP: exec_ORA_ZP,
//	//OP_ORA_ZX: exec_ORA_ZX,
//
//	//OP_PHA: exec_PHA,
//	//OP_PHP: exec_PHP,
//	//OP_PIA: exec_PIA,
//	//OP_PLP: exec_PLP,
//
//	//OP_ROL_AB: exec_ROL_AB,
//	//OP_ROL_AC: exec_ROL_AC,
//	//OP_ROL_AX: exec_ROL_AX,
//	//OP_ROL_ZP: exec_ROL_ZP,
//	//OP_ROL_ZX: exec_ROL_ZX,
//
//	//OP_ROR_AB: exec_ROR_AB,
//	//OP_ROR_AC: exec_ROR_AC,
//	//OP_ROR_AX: exec_ROR_AX,
//	//OP_ROR_ZP: exec_ROR_ZP,
//	//OP_ROR_ZX: exec_ROR_ZX,
//
//	//OP_RTI: exec_RTI,
//	//OP_RTS: exec_RTS,
//
//	//OP_SBC_AB: exec_SBC_AB,
//	//OP_SBC_AX: exec_SBC_AX,
//	//OP_SBC_AY: exec_SBC_AY,
//	//OP_SBC_IM: exec_SBC_IM,
//	//OP_SBC_IX: exec_SBC_IX,
//	//OP_SBC_IY: exec_SBC_IY,
//	//OP_SBC_ZP: exec_SBC_ZP,
//	//OP_SBC_ZX: exec_SBC_ZX,
//
//	//OP_SEC: exec_SEC,
//	//OP_SED: exec_SED,
//	//OP_SEI: exec_SEI,
//
//	OP_STA_AB: exec_STA_AB,
//	OP_STA_AX: exec_STA_AX,
//	OP_STA_AY: exec_STA_AY,
//	OP_STA_IX: exec_STA_IX,
//	OP_STA_IY: exec_STA_IY,
//	OP_STA_ZP: exec_STA_ZP,
//	OP_STA_ZX: exec_STA_ZX,
//
//	OP_STX_AB: exec_STX_AB,
//	OP_STX_ZP: exec_STX_ZP,
//	OP_STX_ZY: exec_STX_ZY,
//
//	//OP_STY_AB: exec_STY_AB,
//	//OP_STY_ZP: exec_STY_ZP,
//	//OP_STY_ZX: exec_STY_ZX,
//
//	OP_TAX: exec_TAX,
//	OP_TAY: exec_TAY,
//	OP_TSX: exec_TSX,
//	OP_TXA: exec_TXA,
//	OP_TXS: exec_TXS,
//	OP_TYA: exec_TYA,
//}

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
	OP_ORA_IM byte = 0x09 //Immediate
	OP_ASL_AC byte = 0x0A //Accumulator
	OP_ORA_AB byte = 0x0D //Absolute
	OP_ASL_AB byte = 0x0E //Absolute
	OP_BPL    byte = 0x10 //
	OP_ORA_IY byte = 0x11 //(Indirect),Y
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
	OP_BCC    byte = 0x90 //
	OP_STA_IY byte = 0x91 //(Indirect),Y
	OP_STY_ZX byte = 0x94 //Zero Page,X
	OP_STA_ZX byte = 0x95 //Zero Page,X
	OP_STX_ZY byte = 0x96 //Zero Page,Y
	OP_TYA    byte = 0x98 //
	OP_STA_AY byte = 0x99 //Absolute,Y
	OP_TXS    byte = 0x9A //
	OP_STA_AX byte = 0x9D //Absolute,X
	OP_LDY_IM byte = 0xA0 //Immediate
	OP_LDX_IM byte = 0xA2 //Immediate
	OP_LDY_ZP byte = 0xA4 //Zero Page
	OP_LDA_ZP byte = 0xA5 //Zero Page
	OP_LDX_ZP byte = 0xA6 //Zero Page
	OP_TAY    byte = 0xA8 //
	OP_LDA_IM byte = 0xA9 //Immediate
	OP_TAX    byte = 0xAA //
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
