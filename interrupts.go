package emu

type Interrupt struct {
	Name string
	vector uint16
	phlags uint8
}

func (i Interrupt) Execute(c *Core) {
	c.pushAddress(c.PC)
	c.pushByte(i.phlags | c.Phlags)
	c.PC = c.ReadWord(i.vector)
}

var interruptList = map[uint16]Interrupt{
	VECTOR_NMI: Interrupt{
		Name:   "NMI",
		vector: VECTOR_NMI,
		phlags: FLAG_IRQ},

	VECTOR_RESET: Interrupt{
		Name:   "RESET",
		vector: VECTOR_RESET,
		phlags: FLAG_IRQ},

	VECTOR_IRQ: Interrupt{
		Name:   "IRQ",
		vector: VECTOR_IRQ,
		phlags: FLAG_IRQ},
}
