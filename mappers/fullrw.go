package mappers

type FullRW struct {
	rom []byte
}

func NewFullRW(data []byte) (Mapper, error) {
	if len(data) != 0x10000 {
		return nil, ErrRomSize
	}

	return &FullRW{
		rom: data,
	}, nil
}

func (rw *FullRW) Name() string {
	return "FullRW"
}

func (rw *FullRW) State() string {
	return rw.Name()
}

func (rw *FullRW) ReadByte(address uint16) uint8 {
	return rw.rom[address]
}

func (rw *FullRW) WriteByte(address uint16, value uint8) {
	rw.rom[address] = value
}

func (rw *FullRW) ClearRam() {}
