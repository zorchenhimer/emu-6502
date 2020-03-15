package mappers

import (
	"errors"
	"fmt"
)

type Mapper interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)

	// Debugging/Info
	Name() string
	State() string
}

var (
	ErrRomSize = errors.New("ROM data is incorrect size")
	//ErrNoMapper = errors.New("")
)

// Read a NES ROM from the given file.  This will read the iNES header
// and return the correct mapper.
func LoadFromFile(filename string) (Mapper, error) {
	return nil, fmt.Errorf("LoadFromFile() not implemented")
}
