package mmu

type Manager interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)

	ClearRam()
}
