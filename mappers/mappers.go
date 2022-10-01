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
	ReadWord(address uint16) uint16

	// Returns the offset in the ROM file given
	// the current bank configuration.
	// If the address is not in ROM space, the original address is returned, as
	// well as False.  If the address is in ROM space, the real offset is
	// returned with True.
	Offset(address uint16) (uint32, bool)

	// Given the current mapper state, is the address in question ROM or RAM?
	IsRom(address uint16) bool

	// GetState returns a mapper-specific snapshot of the internals of its state.
	GetState() any
	// SetState clobbers all current mapper settings with the provided state.
	SetState(data any) error

	// Debugging/Info
	Name() string
	State() string

	DumpFullStack() string

	ClearRam()
}

type CallbackType uint8

const (
	READ CallbackType = 1
	WRITE CallbackType = 2
)

type CallbackFunction func(address uint16, data uint8)

type CallbackMapper interface {
	Mapper

	// Both the individual address and the range callbacks use the same
	// underlying lookup table and will clobber eachother if addresses conflict.
	// In such an event, the last registered callback will be used.
	RegisterReadCallback(address uint16, f CallbackFunction)
	RegisterWriteCallback(address uint16, f CallbackFunction)

	RegisterReadCallbackRange(addressStart, addressEnd uint16, f CallbackFunction)
	RegisterWriteCallbackRange(addressStart, addressEnd uint16, f CallbackFunction)

	// These callbacks will always fire, regardless of address.
	CallbackRead(f CallbackFunction)
	CallbackWrite(f CallbackFunction)

	// Fires when a mapper address is written.  Exact addresses are mapper
	// specific.  This callback does not utilize the same lookup table as the
	// address based callbacks.  The only thing that will overwrite a mapper
	// callback is another mapper callback.
	CallbackMapperWrite(f CallbackFunction)
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
