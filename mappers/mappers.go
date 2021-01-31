package mappers

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/zorchenhimer/go-nes/ines"
)

type mapperNewFunc func(raw []byte, hasRam bool) (Mapper, error)
var availableMappers map[int]mapperNewFunc = make(map[int]mapperNewFunc)

func registerMapper(id int, f mapperNewFunc) {
	if _, exists := availableMappers[id]; exists {
		panic(fmt.Sprintf("Mapper implementation with ID %d already exists", id))
	}

	availableMappers[id] = f
}

type Mapper interface {
	// ReadByte reads a single byte at the given address with the current
	// mapper configuration.  To read from unmapped banks, the mapper needs
	// to be written to to map those banks into CPU address space.
	ReadByte(address uint16) uint8

	// ReadWord does the same thing as ReadByte, but instead returns two
	// consecutive bytes as a uint16.
	ReadWord(address uint16) uint16

	// WriteByte writes a single byte to the mapper.  This is the method
	// in which mapped banks can be changed.  Note that each mapper has
	// their own specific "API" for this.
	WriteByte(address uint16, value uint8)

	// Returns the offset in the PRG ROM given the current bank
	// configuration.  Note that 16 needs to be added for raw
	// file offset to account for the header.
	Offset(address uint16) uint32

	// GetState returns a mapper-specific snapshot of the internals of its state.
	GetState() interface{}

	// SetState clobbers all current mapper settings with the provided state.
	SetState(data interface{}) error

	// Name returns the name of the mapper.  The mapper number will not be included here.
	Name() string

	// State returns the current state of the mapper.  The format is mapper-specific
	// and will differ between mappers.
	State() string

	// Wipes all CPU RAM (0x0000-0x07FFF inclusive) and PRG
	// Work RAM (typically 0x6000-0x7FFF).  Addresses cleared
	// in cartridge space (ie 0x4020-0xFFFF) are mapper
	// dependent may vary depending on the mapper used.
	ClearRam()
}

var (
	ErrRomSize = errors.New("ROM data is incorrect size")
	//ErrNoMapper = errors.New("")
)

// Read a NES ROM from the given file.  This will read the iNES header
// and return the correct mapper.
func LoadFromFile(filename string) (Mapper, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadFromBytes(raw)
}

// LoadFromBytes loads a rom and returns it as a mapper.  The Mapper type is auto-detected.
func LoadFromBytes(raw []byte) (Mapper, error) {
	header, err := ines.ParseHeader(raw)
	if err != nil {
		return nil, err
	}

	// TODO: NES2 submapper IDs
	init, exists := availableMappers[int(header.Mapper)]
	if !exists {
		return nil, fmt.Errorf("Mapper with ID %d not implemented")
	}

	// Assume all mappers have PRGRAM until parsing this info from a
	// NES2 header is implemented.
	return init(raw[16:int(header.PrgSize+16)], true)
}

func wramCopy(dst, src *[0x2000]byte) {
	for i := 0; i < 0x2000; i+=4 {
		dst[i], dst[i+1], dst[i+2], dst[i+3] = src[i], src[i+1], src[i+2], src[i+3]
	}
}

func ramCopy(dst, src *[0x0800]byte) {
	for i := 0; i < 0x0800; i+=4 {
		dst[i], dst[i+1], dst[i+2], dst[i+3] = src[i], src[i+1], src[i+2], src[i+3]
	}
}
