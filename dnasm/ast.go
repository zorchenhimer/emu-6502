package dnasm

type NodeType int
const (
	NT_Read NodeType = iota
	NT_Write
	NT_ReadModifyWrite
	NT_Modify	// INX, ADC, JMP, etc
	NT_Flags	// check flags

	NT_Branch
	NT_JSR		// JSR
)

type Node interface {
	Type() NodeType
	OpCode() uint8
	Label() Label
}

type LinearNode interface {
	Next() Node
}

type BranchingNode interface {
	NextA() Node
	NextB() Node
}

type Branch struct {
	opCode uint8
	label Label

	// branch not taken
	nextA Node
	// branch taken
	nextB Node
}

func (b *Branch) Label() Label  { return b.label }
func (b *Branch) NextA() Node   { return b.nextA }
func (b *Branch) NextB() Node   { return b.nextB }
func (b *Branch) OpCode() uint8 { return b.opCode }

type Linear struct {
	opCode uint8
	label Label
	next Node
}

func (l *Linear) Label() Label  { return l.label }
func (l *Linear) Next() Node    { return l.next }
func (l *Linear) OpCode() uint8 { return l.opCode }
