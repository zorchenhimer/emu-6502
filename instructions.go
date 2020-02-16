package emu

import (
	//"fmt"
)

type ExecFunc func(c *Core, address uint16)

type Instruction interface {
	Execute(c *Core)
	Name() string
	InstrLength(c *Core) uint8
	AddressMeta() AddressModeMeta
}

var instructionList = map[byte]Instruction{

	OP_ADC_AB: StandardInstruction{
		OpCode:         OP_ADC_AB,
		Instruction:    "ADC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_ADC},
	OP_ADC_AX: StandardInstruction{
		OpCode:         OP_ADC_AX,
		Instruction:    "ADC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_ADC},
	OP_ADC_AY: StandardInstruction{
		OpCode:         OP_ADC_AY,
		Instruction:    "ADC",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_ADC},
	OP_ADC_IM: StandardInstruction{
		OpCode:         OP_ADC_IM,
		Instruction:    "ADC",
		AddressMode: ADDR_Immediate,
		Exec:           instr_ADC},
	OP_ADC_IX: StandardInstruction{
		OpCode:         OP_ADC_IX,
		Instruction:    "ADC",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_ADC},
	OP_ADC_IY: StandardInstruction{
		OpCode:         OP_ADC_IY,
		Instruction:    "ADC",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_ADC},
	OP_ADC_ZP: StandardInstruction{
		OpCode:         OP_ADC_ZP,
		Instruction:    "ADC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_ADC},
	OP_ADC_ZX: StandardInstruction{
		OpCode:         OP_ADC_ZX,
		Instruction:    "ADC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_ADC},

	OP_AND_AB: StandardInstruction{
		OpCode:         OP_AND_AB,
		Instruction:    "AND",
		AddressMode: ADDR_Absolute,
		Exec:           instr_AND},
	OP_AND_AX: StandardInstruction{
		OpCode:         OP_AND_AX,
		Instruction:    "AND",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_AND},
	OP_AND_AY: StandardInstruction{
		OpCode:         OP_AND_AY,
		Instruction:    "AND",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_AND},
	OP_AND_IM: StandardInstruction{
		OpCode:         OP_AND_IM,
		Instruction:    "AND",
		AddressMode: ADDR_Immediate,
		Exec:           instr_AND},
	OP_AND_IX: StandardInstruction{
		OpCode:         OP_AND_IX,
		Instruction:    "AND",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_AND},
	OP_AND_IY: StandardInstruction{
		OpCode:         OP_AND_IY,
		Instruction:    "AND",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_AND},
	OP_AND_ZP: StandardInstruction{
		OpCode:         OP_AND_ZP,
		Instruction:    "AND",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_AND},
	OP_AND_ZX: StandardInstruction{
		OpCode:         OP_AND_ZX,
		Instruction:    "AND",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_AND},

	OP_ASL_AB: ReadModifyWrite{
		OpCode:         OP_ASL_AB,
		Instruction:    "ASL",
		AddressMode: ADDR_Absolute,
		Exec:           instr_ASL},
	OP_ASL_AX: ReadModifyWrite{
		OpCode:         OP_ASL_AX,
		Instruction:    "ASL",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_ASL},
	OP_ASL_ZP: ReadModifyWrite{
		OpCode:         OP_ASL_ZP,
		Instruction:    "ASL",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_ASL},
	OP_ASL_ZX: ReadModifyWrite{
		OpCode:         OP_ASL_ZX,
		Instruction:    "ASL",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_ASL},

	OP_BIT_AB: StandardInstruction{
		OpCode:         OP_BIT_AB,
		Instruction:    "BIT",
		AddressMode: ADDR_Absolute,
		Exec:           instr_BIT},
	OP_BIT_ZP: StandardInstruction{
		OpCode:         OP_BIT_ZP,
		Instruction:    "BIT",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_BIT},

	OP_BCC: Branch{
		OpCode: OP_BCC,
		Instruction: "BCC",
		Flag: FLAG_CARRY,
		Set: false},
	OP_BCS: Branch{
		OpCode: OP_BCS,
		Instruction: "BCS",
		Flag: FLAG_CARRY,
		Set: true},
	OP_BEQ: Branch{
		OpCode: OP_BEQ,
		Instruction: "BEQ",
		Flag: FLAG_ZERO,
		Set: true},
	OP_BMI: Branch{
		OpCode: OP_BMI,
		Instruction: "BMI",
		Flag: FLAG_NEGATIVE,
		Set: true},
	OP_BNE: Branch{
		OpCode: OP_BNE,
		Instruction: "BNE",
		Flag: FLAG_ZERO,
		Set: false},
	OP_BPL: Branch{
		OpCode: OP_BPL,
		Instruction: "BPL",
		Flag: FLAG_NEGATIVE,
		Set: false},
	OP_BVC: Branch{
		OpCode: OP_BVC,
		Instruction: "BVC",
		Flag: FLAG_OVERFLOW,
		Set: false},
	OP_BVS: Branch{
		OpCode: OP_BVS,
		Instruction: "BVS",
		Flag: FLAG_OVERFLOW,
		Set: true},

	OP_BRK: Jump{
		OpCode:         OP_BRK,
		Instruction:    "BRK",
		AddressMode: ADDR_Implied,
		Exec:           instr_BRK},

	OP_CLC: StandardInstruction{
		OpCode:         OP_CLC,
		Instruction:    "CLC",
		AddressMode: ADDR_Implied,
		Exec:           instr_CLC},
	OP_CLD: StandardInstruction{
		OpCode:         OP_CLD,
		Instruction:    "CLD",
		AddressMode: ADDR_Implied,
		Exec:           instr_CLD},
	OP_CLI: StandardInstruction{
		OpCode:         OP_CLI,
		Instruction:    "CLI",
		AddressMode: ADDR_Implied,
		Exec:           instr_CLI},
	OP_CLV: StandardInstruction{
		OpCode:         OP_CLV,
		Instruction:    "CLV",
		AddressMode: ADDR_Implied,
		Exec:           instr_CLV},

	OP_CMP_AB: StandardInstruction{
		OpCode:         OP_CMP_AB,
		Instruction:    "CMP",
		AddressMode: ADDR_Absolute,
		Exec:           instr_CMP},
	OP_CMP_AX: StandardInstruction{
		OpCode:         OP_CMP_AX,
		Instruction:    "CMP",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_CMP},
	OP_CMP_AY: StandardInstruction{
		OpCode:         OP_CMP_AY,
		Instruction:    "CMP",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_CMP},
	OP_CMP_IM: StandardInstruction{
		OpCode:         OP_CMP_IM,
		Instruction:    "CMP",
		AddressMode: ADDR_Immediate,
		Exec:           instr_CMP},
	OP_CMP_IX: StandardInstruction{
		OpCode:         OP_CMP_IX,
		Instruction:    "CMP",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_CMP},
	OP_CMP_IY: StandardInstruction{
		OpCode:         OP_CMP_IY,
		Instruction:    "CMP",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_CMP},
	OP_CMP_ZP: StandardInstruction{
		OpCode:         OP_CMP_ZP,
		Instruction:    "CMP",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_CMP},
	OP_CMP_ZX: StandardInstruction{
		OpCode:         OP_CMP_ZX,
		Instruction:    "CMP",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_CMP},

	OP_CPX_AB: StandardInstruction{
		OpCode:         OP_CPX_AB,
		Instruction:    "CPX",
		AddressMode: ADDR_Absolute,
		Exec:           instr_CPX},
	OP_CPX_IM: StandardInstruction{
		OpCode:         OP_CPX_IM,
		Instruction:    "CPX",
		AddressMode: ADDR_Immediate,
		Exec:           instr_CPX},
	OP_CPX_ZP: StandardInstruction{
		OpCode:         OP_CPX_ZP,
		Instruction:    "CPX",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_CPX},

	OP_CPY_AB: StandardInstruction{
		OpCode:         OP_CPY_AB,
		Instruction:    "CPY",
		AddressMode: ADDR_Absolute,
		Exec:           instr_CPY},
	OP_CPY_IM: StandardInstruction{
		OpCode:         OP_CPY_IM,
		Instruction:    "CPY",
		AddressMode: ADDR_Immediate,
		Exec:           instr_CPY},
	OP_CPY_ZP: StandardInstruction{
		OpCode:         OP_CPY_ZP,
		Instruction:    "CPY",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_CPY},

	OP_DEC_AB: ReadModifyWrite{
		OpCode:         OP_DEC_AB,
		Instruction:    "DEC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_DEC},
	OP_DEC_AX: ReadModifyWrite{
		OpCode:         OP_DEC_AX,
		Instruction:    "DEC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_DEC},
	OP_DEC_ZP: ReadModifyWrite{
		OpCode:         OP_DEC_ZP,
		Instruction:    "DEC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_DEC},
	OP_DEC_ZX: ReadModifyWrite{
		OpCode:         OP_DEC_ZX,
		Instruction:    "DEC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_DEC},

	OP_DEX: StandardInstruction{
		OpCode:         OP_DEX,
		Instruction:    "DEX",
		AddressMode: ADDR_Implied,
		Exec:           instr_DEX},
	OP_DEY: StandardInstruction{
		OpCode:         OP_DEY,
		Instruction:    "DEY",
		AddressMode: ADDR_Implied,
		Exec:           instr_DEY},

	OP_EOR_AB: StandardInstruction{
		OpCode:         OP_EOR_AB,
		Instruction:    "EOR",
		AddressMode: ADDR_Absolute,
		Exec:           instr_EOR},
	OP_EOR_AX: StandardInstruction{
		OpCode:         OP_EOR_AX,
		Instruction:    "EOR",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_EOR},
	OP_EOR_AY: StandardInstruction{
		OpCode:         OP_EOR_AY,
		Instruction:    "EOR",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_EOR},
	OP_EOR_IM: StandardInstruction{
		OpCode:         OP_EOR_IM,
		Instruction:    "EOR",
		AddressMode: ADDR_Immediate,
		Exec:           instr_EOR},
	OP_EOR_IX: StandardInstruction{
		OpCode:         OP_EOR_IX,
		Instruction:    "EOR",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_EOR},
	OP_EOR_IY: StandardInstruction{
		OpCode:         OP_EOR_IY,
		Instruction:    "EOR",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_EOR},
	OP_EOR_ZP: StandardInstruction{
		OpCode:         OP_EOR_ZP,
		Instruction:    "EOR",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_EOR},
	OP_EOR_ZX: StandardInstruction{
		OpCode:         OP_EOR_ZX,
		Instruction:    "EOR",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_EOR},

	OP_LSR_AB: ReadModifyWrite{
		OpCode:         OP_LSR_AB,
		Instruction:    "LSR",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LSR},
	OP_LSR_AX: ReadModifyWrite{
		OpCode:         OP_LSR_AX,
		Instruction:    "LSR",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_LSR},
	OP_LSR_ZP: ReadModifyWrite{
		OpCode:         OP_LSR_ZP,
		Instruction:    "LSR",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LSR},
	OP_LSR_ZX: ReadModifyWrite{
		OpCode:         OP_LSR_ZX,
		Instruction:    "LSR",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_LSR},

	OP_JMP_AB: Jump{
		OpCode: OP_JMP_AB,
		Instruction: "JMP",
		AddressMode: ADDR_Absolute,
		Exec: instr_JMP},
	OP_JMP_ID: Jump{
		OpCode: OP_JMP_ID,
		Instruction: "JMP",
		AddressMode: ADDR_Indirect,
		Exec: instr_JMP},
	OP_JSR: Jump{
		OpCode: OP_JSR,
		Instruction: "JSR",
		AddressMode: ADDR_Absolute,
		Exec: instr_JSR},

	OP_LDA_AB: StandardInstruction{
		OpCode:         OP_LDA_AB,
		Instruction:    "LDA",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDA},
	OP_LDA_AX: StandardInstruction{
		OpCode:         OP_LDA_AX,
		Instruction:    "LDA",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_LDA},
	OP_LDA_AY: StandardInstruction{
		OpCode:         OP_LDA_AY,
		Instruction:    "LDA",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_LDA},
	OP_LDA_IM: StandardInstruction{
		OpCode:         OP_LDA_IM,
		Instruction:    "LDA",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDA},
	OP_LDA_IX: StandardInstruction{
		OpCode:         OP_LDA_IX,
		Instruction:    "LDA",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_LDA},
	OP_LDA_IY: StandardInstruction{
		OpCode:         OP_LDA_IY,
		Instruction:    "LDA",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_LDA},
	OP_LDA_ZP: StandardInstruction{
		OpCode:         OP_LDA_ZP,
		Instruction:    "LDA",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDA},
	OP_LDA_ZX: StandardInstruction{
		OpCode:         OP_LDA_ZX,
		Instruction:    "LDA",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_LDA},

	OP_LDX_AB: StandardInstruction{
		OpCode:         OP_LDX_AB,
		Instruction:    "LDX",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDX},
	OP_LDX_AY: StandardInstruction{
		OpCode:         OP_LDX_AY,
		Instruction:    "LDX",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_LDX},
	OP_LDX_IM: StandardInstruction{
		OpCode:         OP_LDX_IM,
		Instruction:    "LDX",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDX},
	OP_LDX_ZP: StandardInstruction{
		OpCode:         OP_LDX_ZP,
		Instruction:    "LDX",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDX},
	OP_LDX_ZY: StandardInstruction{
		OpCode:         OP_LDX_ZY,
		Instruction:    "LDX",
		AddressMode: ADDR_ZeroPageY,
		Exec:           instr_LDX},

	OP_LDY_AB: StandardInstruction{
		OpCode:         OP_LDY_AB,
		Instruction:    "LDY",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDY},
	OP_LDY_AX: StandardInstruction{
		OpCode:         OP_LDY_AX,
		Instruction:    "LDY",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_LDY},
	OP_LDY_IM: StandardInstruction{
		OpCode:         OP_LDY_IM,
		Instruction:    "LDY",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDY},
	OP_LDY_ZP: StandardInstruction{
		OpCode:         OP_LDY_ZP,
		Instruction:    "LDY",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDY},
	OP_LDY_ZX: StandardInstruction{
		OpCode:         OP_LDY_ZX,
		Instruction:    "LDY",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_LDY},

	OP_INC_AB: ReadModifyWrite{
		OpCode:         OP_INC_AB,
		Instruction:    "INC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_INC},
	OP_INC_AX: ReadModifyWrite{
		OpCode:         OP_INC_AX,
		Instruction:    "INC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_INC},
	OP_INC_ZP: ReadModifyWrite{
		OpCode:         OP_INC_ZP,
		Instruction:    "INC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_INC},
	OP_INC_ZX: ReadModifyWrite{
		OpCode:         OP_INC_ZX,
		Instruction:    "INC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_INC},

	OP_INX: StandardInstruction{
		OpCode:         OP_INX,
		Instruction:    "INX",
		AddressMode: ADDR_Implied,
		Exec:           instr_INX},
	OP_INY: StandardInstruction{
		OpCode:         OP_INY,
		Instruction:    "INY",
		AddressMode: ADDR_Implied,
		Exec:           instr_INY},

	OP_NOP: StandardInstruction{
		OpCode:         OP_NOP,
		Instruction:    "NOP",
		AddressMode: ADDR_Implied,
		Exec:           instr_NOP},

	OP_ORA_AB: StandardInstruction{
		OpCode:         OP_ORA_AB,
		Instruction:    "ORA",
		AddressMode: ADDR_Absolute,
		Exec:           instr_ORA},
	OP_ORA_AX: StandardInstruction{
		OpCode:         OP_ORA_AX,
		Instruction:    "ORA",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_ORA},
	OP_ORA_AY: StandardInstruction{
		OpCode:         OP_ORA_AY,
		Instruction:    "ORA",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_ORA},
	OP_ORA_IM: StandardInstruction{
		OpCode:         OP_ORA_IM,
		Instruction:    "ORA",
		AddressMode: ADDR_Immediate,
		Exec:           instr_ORA},
	OP_ORA_IX: StandardInstruction{
		OpCode:         OP_ORA_IX,
		Instruction:    "ORA",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_ORA},
	OP_ORA_IY: StandardInstruction{
		OpCode:         OP_ORA_IY,
		Instruction:    "ORA",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_ORA},
	OP_ORA_ZP: StandardInstruction{
		OpCode:         OP_ORA_ZP,
		Instruction:    "ORA",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_ORA},
	OP_ORA_ZX: StandardInstruction{
		OpCode:         OP_ORA_ZX,
		Instruction:    "ORA",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_ORA},

	OP_PHA: StandardInstruction{
		OpCode:         OP_PHA,
		Instruction:    "PHA",
		AddressMode: ADDR_Implied,
		Exec:           instr_PHA},
	OP_PHP: StandardInstruction{
		OpCode:         OP_PHP,
		Instruction:    "PHP",
		AddressMode: ADDR_Implied,
		Exec:           instr_PHP},
	OP_PLA: StandardInstruction{
		OpCode:         OP_PLA,
		Instruction:    "PLA",
		AddressMode: ADDR_Implied,
		Exec:           instr_PLA},
	OP_PLP: StandardInstruction{
		OpCode:         OP_PLP,
		Instruction:    "PLP",
		AddressMode: ADDR_Implied,
		Exec:           instr_PLP},

	OP_ROL_AB: ReadModifyWrite{
		OpCode:         OP_ROL_AB,
		Instruction:    "ROL",
		AddressMode: ADDR_Absolute,
		Exec:           instr_ROL},
	OP_ROL_AX: ReadModifyWrite{
		OpCode:         OP_ROL_AX,
		Instruction:    "ROL",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_ROL},
	OP_ROL_ZP: ReadModifyWrite{
		OpCode:         OP_ROL_ZP,
		Instruction:    "ROL",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_ROL},
	OP_ROL_ZX: ReadModifyWrite{
		OpCode:         OP_ROL_ZX,
		Instruction:    "ROL",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_ROL},

	OP_ROR_AB: ReadModifyWrite{
		OpCode:         OP_ROR_AB,
		Instruction:    "ROR",
		AddressMode: ADDR_Absolute,
		Exec:           instr_ROR},
	OP_ROR_AX: ReadModifyWrite{
		OpCode:         OP_ROR_AX,
		Instruction:    "ROR",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_ROR},
	OP_ROR_ZP: ReadModifyWrite{
		OpCode:         OP_ROR_ZP,
		Instruction:    "ROR",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_ROR},
	OP_ROR_ZX: ReadModifyWrite{
		OpCode:         OP_ROR_ZX,
		Instruction:    "ROR",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_ROR},

	OP_RTI: Jump{
		OpCode: OP_RTI,
		Instruction: "RTI",
		AddressMode: ADDR_Implied,
		Exec: instr_RTI},
	OP_RTS: Jump{
		OpCode: OP_RTS,
		Instruction: "RTS",
		AddressMode: ADDR_Implied,
		Exec: instr_RTS},

	OP_SBC_AB: StandardInstruction{
		OpCode:         OP_SBC_AB,
		Instruction:    "SBC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_SBC},
	OP_SBC_AX: StandardInstruction{
		OpCode:         OP_SBC_AX,
		Instruction:    "SBC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_SBC},
	OP_SBC_AY: StandardInstruction{
		OpCode:         OP_SBC_AY,
		Instruction:    "SBC",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_SBC},
	OP_SBC_IM: StandardInstruction{
		OpCode:         OP_SBC_IM,
		Instruction:    "SBC",
		AddressMode: ADDR_Immediate,
		Exec:           instr_SBC},
	OP_SBC_IX: StandardInstruction{
		OpCode:         OP_SBC_IX,
		Instruction:    "SBC",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_SBC},
	OP_SBC_IY: StandardInstruction{
		OpCode:         OP_SBC_IY,
		Instruction:    "SBC",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_SBC},
	OP_SBC_ZP: StandardInstruction{
		OpCode:         OP_SBC_ZP,
		Instruction:    "SBC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_SBC},
	OP_SBC_ZX: StandardInstruction{
		OpCode:         OP_SBC_ZX,
		Instruction:    "SBC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_SBC},

	OP_SEC: StandardInstruction{
		OpCode:         OP_SEC,
		Instruction:    "SEC",
		AddressMode: ADDR_Implied,
		Exec:           instr_SEC},
	OP_SED: StandardInstruction{
		OpCode:         OP_SED,
		Instruction:    "SED",
		AddressMode: ADDR_Implied,
		Exec:           instr_SED},
	OP_SEI: StandardInstruction{
		OpCode:         OP_SEI,
		Instruction:    "SEI",
		AddressMode: ADDR_Implied,
		Exec:           instr_SEI},

	OP_STA_AB: StandardInstruction{
		OpCode:         OP_STA_AB,
		Instruction:    "STA",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STA},
	OP_STA_AX: StandardInstruction{
		OpCode:         OP_STA_AX,
		Instruction:    "STA",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_STA},
	OP_STA_AY: StandardInstruction{
		OpCode:         OP_STA_AY,
		Instruction:    "STA",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_STA},
	OP_STA_IX: StandardInstruction{
		OpCode:         OP_STA_IX,
		Instruction:    "STA",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_STA},
	OP_STA_IY: StandardInstruction{
		OpCode:         OP_STA_IY,
		Instruction:    "STA",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_STA},
	OP_STA_ZP: StandardInstruction{
		OpCode:         OP_STA_ZP,
		Instruction:    "STA",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STA},
	OP_STA_ZX: StandardInstruction{
		OpCode:         OP_STA_ZX,
		Instruction:    "STA",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_STA},

	OP_STX_AB: StandardInstruction{
		OpCode:         OP_STX_AB,
		Instruction:    "STX",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STX},
	OP_STX_ZP: StandardInstruction{
		OpCode:         OP_STX_ZP,
		Instruction:    "STX",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STX},
	OP_STX_ZY: StandardInstruction{
		OpCode:         OP_STX_ZY,
		Instruction:    "STX",
		AddressMode: ADDR_ZeroPageY,
		Exec:           instr_STX},

	OP_STY_AB: StandardInstruction{
		OpCode:         OP_STY_AB,
		Instruction:    "STY",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STY},
	OP_STY_ZP: StandardInstruction{
		OpCode:         OP_STY_ZP,
		Instruction:    "STY",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STY},
	OP_STY_ZX: StandardInstruction{
		OpCode:         OP_STY_ZX,
		Instruction:    "STY",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_STY},

	OP_TAX: StandardInstruction{
		OpCode:         OP_TAX,
		Instruction:    "TAX",
		AddressMode: ADDR_Implied,
		Exec:           instr_TAX},
	OP_TAY: StandardInstruction{
		OpCode:         OP_TAY,
		Instruction:    "TAY",
		AddressMode: ADDR_Implied,
		Exec:           instr_TAY},
	OP_TSX: StandardInstruction{
		OpCode:         OP_TSX,
		Instruction:    "TSX",
		AddressMode: ADDR_Implied,
		Exec:           instr_TSX},
	OP_TXA: StandardInstruction{
		OpCode:         OP_TXA,
		Instruction:    "TXA",
		AddressMode: ADDR_Implied,
		Exec:           instr_TXA},
	OP_TXS: StandardInstruction{
		OpCode:         OP_TXS,
		Instruction:    "TXS",
		AddressMode: ADDR_Implied,
		Exec:           instr_TXS},
	OP_TYA: StandardInstruction{
		OpCode:         OP_TYA,
		Instruction:    "TYA",
		AddressMode: ADDR_Implied,
		Exec:           instr_TYA},
}

type StandardInstruction struct {
	AddressMode AddressModeMeta
	OpCode      byte
	Instruction string
	Exec        ExecFunc
}

func (i StandardInstruction) AddressMeta() AddressModeMeta {
	return i.AddressMode
}

func (i StandardInstruction) Execute(c *Core) {
	address, size := i.AddressMode.Address(c)
	i.Exec(c, address)
	c.PC += uint16(size)
}

func (i StandardInstruction) InstrLength(c *Core) uint8 {
	_, size := i.AddressMode.Address(c)
	return size
}

func (i StandardInstruction) Name() string {
	return i.Instruction
}

func instr_CLC(c *Core, address uint16) {
	c.Phlags &^= FLAG_CARRY
}

func instr_CLD(c *Core, address uint16) {
	c.Phlags &^= FLAG_DECIMAL
}

func instr_CLI(c *Core, address uint16) {
	c.Phlags &^= FLAG_INTERRUPT
}

func instr_CLV(c *Core, address uint16) {
	c.Phlags &^= FLAG_OVERFLOW
}

func (c *Core) compare(a, b uint8) {
	overflow := c.Phlags & FLAG_OVERFLOW
        c.Phlags |= FLAG_CARRY

	c.twosCompAdd(a, b ^ 0xff)
	c.Phlags = (c.Phlags &^ FLAG_OVERFLOW) | overflow
}

func instr_CMP(c *Core, address uint16) {
	c.compare(c.A, c.ReadByte(address))
}

func instr_CPX(c *Core, address uint16) {
	c.compare(c.X, c.ReadByte(address))
}

func instr_CPY(c *Core, address uint16) {
	c.compare(c.Y, c.ReadByte(address))
}

func instr_ADC(c *Core, address uint16) {
	c.A = c.twosCompAdd(c.A, c.ReadByte(address))
}

func instr_AND(c *Core, address uint16) {
	c.A &= c.ReadByte(address)
	c.setZeroNegative(c.A)
}

func instr_BIT(c *Core, address uint16) {
	val := c.A & c.ReadByte(address)
	c.setZeroNegative(val)
	c.Phlags = (c.Phlags &^ FLAG_OVERFLOW) | (val & FLAG_OVERFLOW)
}

func instr_DEX(c *Core, address uint16) {
	c.X -= 1
	c.setZeroNegative(c.X)
}

func instr_DEY(c *Core, address uint16) {
	c.Y -= 1
	c.setZeroNegative(c.Y)
}

func instr_EOR(c *Core, address uint16) {
	c.A = c.A ^ c.ReadByte(address)
	c.setZeroNegative(c.A)
}

func instr_INX(c *Core, address uint16) {
	c.X += 1
	c.setZeroNegative(c.X)
}

func instr_INY(c *Core, address uint16) {
	c.Y += 1
	c.setZeroNegative(c.Y)
}

func instr_LDA(c *Core, address uint16) {
	c.A = c.ReadByte(address)
	c.setZeroNegative(c.A)
}

func instr_LDX(c *Core, address uint16) {
	c.X = c.ReadByte(address)
	c.setZeroNegative(c.X)
}

func instr_LDY(c *Core, address uint16) {
	c.Y = c.ReadByte(address)
	c.setZeroNegative(c.Y)
}

func instr_NOP(c *Core, address uint16) {
	return
}

func instr_ORA(c *Core, address uint16) {
	c.A |= c.ReadByte(address)
}

func instr_PHA(c *Core, address uint16) {
	c.pushByte(c.A)
}

func instr_PHP(c *Core, address uint16) {
	c.pushByte(c.Phlags | FLAG_BREAK)
}

func instr_PLA(c *Core, address uint16) {
	c.A = c.pullByte()
	c.setZeroNegative(c.A)
}

func instr_PLP(c *Core, address uint16) {
	c.Phlags = c.pullByte() & (0xCF) // ignore bits 4 and 5
}

func instr_SBC(c *Core, address uint16) {
	c.A = c.twosCompSubtract(c.A, c.ReadByte(address))
}

func instr_SEC(c *Core, address uint16) {
	c.Phlags |= FLAG_CARRY
}

func instr_SED(c *Core, address uint16) {
	c.Phlags |= FLAG_DECIMAL
}

func instr_SEI(c *Core, address uint16) {
	c.Phlags |= FLAG_INTERRUPT
}

func instr_STA(c *Core, address uint16) {
	c.WriteByte(address, c.A)
}

func instr_STX(c *Core, address uint16) {
	c.WriteByte(address, c.X)
}

func instr_STY(c *Core, address uint16) {
	c.WriteByte(address, c.Y)
}

func instr_TAX(c *Core, address uint16) {
	c.X = c.A
	c.setZeroNegative(c.X)
}

func instr_TAY(c *Core, address uint16) {
	c.Y = c.A
	c.setZeroNegative(c.Y)
}

func instr_TSX(c *Core, address uint16) {
	c.X = c.SP
	c.setZeroNegative(c.X)
}

func instr_TXA(c *Core, address uint16) {
	c.A = c.X
	c.setZeroNegative(c.A)
}

func instr_TXS(c *Core, address uint16) {
	c.SP = c.X
}

func instr_TYA(c *Core, address uint16) {
	c.A = c.Y
	c.setZeroNegative(c.A)
}

type ReadModifyWrite struct {
	OpCode         byte
	Instruction    string
	AddressMode AddressModeMeta
	Exec           func(c *Core, value uint8) uint8
}

func (rmw ReadModifyWrite) AddressMeta() AddressModeMeta {
	return rmw.AddressMode
}

func (rmw ReadModifyWrite) Execute(c *Core) {
	address, size := rmw.AddressMode.Address(c)
	c.WriteByte(address, rmw.Exec(c, c.ReadByte(address)))
	c.PC += uint16(size)
}

func (rmw ReadModifyWrite) Name() string {
	return rmw.Instruction
}

func (rmw ReadModifyWrite) InstrLength(c *Core) uint8 {
	_, size := rmw.AddressMode.Address(c)
	return size
}

func instr_DEC(c *Core, value uint8) uint8 {
	value -= 1
	c.setZeroNegative(value)
	return value
}

func instr_INC(c *Core, value uint8) uint8 {
	value += 1
	c.setZeroNegative(value)
	return value
}

func instr_LSR(c *Core, value uint8) uint8 {
	c.Phlags = (c.Phlags &^ FLAG_CARRY) | (value & FLAG_CARRY)
	return c.setZeroNegative(value >> 1)
}

func instr_ASL(c *Core, value uint8) uint8 {
	c.Phlags = (c.Phlags &^ FLAG_CARRY) | ((value & 0x80) >> 7)
	return c.setZeroNegative(value << 1)
}

func instr_ROL(c *Core, value uint8) uint8 {
	carry := c.Phlags & FLAG_CARRY
	c.Phlags = (c.Phlags &^ FLAG_CARRY) | ((value & 0x80) >> 7)
	return c.setZeroNegative((value << 1) | carry)
}

func instr_ROR(c *Core, value uint8) uint8 {
	carry := (c.Phlags & FLAG_CARRY) << 7
	c.Phlags = (c.Phlags &^ FLAG_CARRY) | (value & FLAG_CARRY)
	return c.setZeroNegative((value >> 1) | carry)
}

type Branch struct {
	OpCode byte
	Instruction string
	Flag uint8
	Set bool
}

func (b Branch) AddressMeta() AddressModeMeta {
	return ADDR_Relative
}

func (b Branch) Name() string {
	return b.Instruction
}

func (b Branch) Execute(c *Core) {
	var v uint8 = 0
	if b.Set {
		v = b.Flag
	}

	if (c.Phlags & b.Flag) == v {
		c.PC = c.addrRelative(c.PC, c.ReadByte(c.PC + 1))
	} else {
		c.PC += 2
	}
}

func (b Branch) InstrLength(c *Core) uint8 {
	return 2
}

// anything that modifies the PC directly, aside form branches
type Jump struct {
	OpCode byte
	Instruction string
	AddressMode AddressModeMeta
	Exec func(c *Core, address uint16) uint16
}

func (j Jump) Name() string {
	return j.Instruction
}

func (j Jump) Execute(c *Core) {
	address, _ := j.AddressMode.Address(c)
	c.PC = j.Exec(c, address)
}

func (j Jump) InstrLength(c *Core) uint8 {
	_, size := j.AddressMode.Address(c)
	return size
}

func (j Jump) AddressMeta() AddressModeMeta {
	return j.AddressMode
}

func instr_JMP(c *Core, address uint16) uint16 {
	return address
}

func instr_JSR(c *Core, address uint16) uint16 {
	c.pushAddress(c.PC + 2)
	return address
}

func instr_RTS(c *Core, address uint16) uint16 {
	return c.pullAddress() + 1
}

func instr_RTI(c *Core, address uint16) uint16 {
	c.Phlags = c.pullByte()
	return c.pullAddress()
}

func instr_BRK(c *Core, address uint16) uint16 {
	c.pushAddress(c.PC + 2)
	c.pushByte(c.Phlags | FLAG_BREAK)
	c.Phlags = c.Phlags | FLAG_INTERRUPT
	return c.ReadWord(0xFFFE)
}
