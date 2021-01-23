package emu

import (
	"fmt"
)

type TokenType int
const (
	TT_Unknown TokenType = iota
	TT_InstructionImplied
	TT_InstructionByte
	TT_InstructionWord
	TT_InstructionBranch
	TT_Data
)

type Token interface {
	Type() TokenType
	String() string
}

type InstructionAny interface {
	Instruction
	String() string
	OpCode() uint8
}

type InstructionImplied struct {
	Instruction
	opCode uint8
}

type InstructionByte struct {
	Instruction
	opCode uint8
	arg uint8
}

type InstructionWord struct {
	Instruction
	opCode uint8
	arg uint16
}

type InstructionBranch struct {
	Instruction
	opCode uint8
	dest uint16
	arg uint8
}

func (i InstructionImplied) Type() TokenType { return TT_InstructionImplied }
func (i InstructionByte) Type() TokenType { return TT_InstructionByte }
func (i InstructionWord) Type() TokenType { return TT_InstructionWord }
func (i InstructionBranch) Type() TokenType { return TT_InstructionBranch }

func (i InstructionImplied) OpCode() uint8 { return i.opCode }
func (i InstructionByte) OpCode() uint8 { return i.opCode }
func (i InstructionWord) OpCode() uint8 { return i.opCode }
func (i InstructionBranch) OpCode() uint8 { return i.opCode }

func (i InstructionByte) Arg() uint8 { return i.arg }
func (i InstructionWord) Arg() uint16 { return i.arg }
func (i InstructionBranch) Arg() uint8 { return i.arg }

func (i InstructionBranch) Destination() uint16 { return i.dest }

func (i InstructionImplied) String() string {
	return instructionList[i.opCode].Name()
}

func (i InstructionByte) String() string {
	return instructionList[i.opCode].Name() +" "+ instructionList[i.opCode].AddressMeta().TokenAsm(uint16(i.arg))
}

func (i InstructionWord) String() string {
	return instructionList[i.opCode].Name() +" "+ instructionList[i.opCode].AddressMeta().TokenAsm(i.arg)
}

func (i InstructionBranch) String() string {
	//return instructionList[i.opCode].Name() +" "+ instructionList[i.opCode].AddressMeta().TokenAsm(i.Dest)
	return fmt.Sprintf("%s $%04X (%d)", instructionList[i.opCode].Name(), i.dest, int8(i.arg))
}


type Data struct {
	raw byte
}

func (d *Data) Type() TokenType { return TT_Data }
func (d *Data) String() string { return fmt.Sprintf("$%02X", d.raw) }

