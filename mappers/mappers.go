package mappers

import (
	"bytes"
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
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
	//ReadWord(address uint16) uint16

	// Returns the offset in the ROM given
	// the current bank configuration.  Offset does
	// not include the iNES header.
	Offset(address uint16) uint32

	// TODO
	// Given the current mapper configuration, is
	// the provided CPU address in RAM?
	//IsRam(address uint16) bool

	// GetState returns a mapper-specific snapshot of the internals of its state.
	GetState() interface{}
	// SetState clobbers all current mapper settings with the provided state.
	SetState(data interface{}) error

	Info() Info

	// Debugging/Info
	Name() string
	State() string

	//DumpFullStack() string

	ClearRam()
}

type Info struct {
	PrgSize uint
	PrgRamSize uint
	PrgBankSize uint

	ChrSize uint
	ChrRamSize uint
	ChrBankSize uint

	// Start addresses in CPU space
	PrgStartAddress uint16
	PrgRamStartAddress uint16
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
	// Handle studybox files
	if bytes.Equal(raw[:4], []byte("STBX")) {
		return NewStudyBox(raw, true)
	}

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
