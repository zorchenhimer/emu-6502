package mmu

import (
	"io"
	"fmt"

	fdslib "git.zorchenhimer.com/zorchenhimer/go-fds"

	"github.com/zorchenhimer/emu-6502/labels"
)

type FDS struct {
	bios   [fdsBiosSize]byte
	fdsRam [fdsRamSize]byte
	cpuRam [0x800]byte
	ppuRam [fdsPpuSize]byte

	// currently loaded disk
	//disk io.ReadSeeker
	disk *fdslib.Rom
	diskSide int

	irqCounter uint16
	irqEnabled bool
	irqReload  uint16
	irqRepeat  bool

	// bit 7 on $4025.  generate IRQ on byte transfer flag raising
	irqTransferEnabled bool

	diskIoEnabled  bool
	soundIoEnabled bool

	transferMode fdsTransferMode
	mirrorMode   fdsMirrorMode
	motorControl fdsMotorControl
	crcEnabled   bool
}

type fdsTransferMode int

const (
	fdsWrite fdsTransferMode = 0
	fdsRead  fdsTransferMode = 1
)

type fdsMirrorMode int

const (
	fdsHorizontal fdsMirrorMode = 0
	fdsVertical   fdsMirrorMode = 1
)

type fdsMotorControl int

const (
	fdsStart fdsMirrorMode = 0
	fdsStop  fdsMirrorMode = 1
)

const (
	fdsRamSize  uint16 = 32*1024
	fdsBiosSize uint16 = 8*1024
	fdsPpuSize  uint16 = 8*1024
)

func NewFDS(bios io.Reader) (*FDS, error) {
	rawBios, err := io.ReadAll(bios)
	if err != nil {
		return nil, fmt.Errorf("Unable to read BIOS: %w", err)
	}

	if len(rawBios) != 8*1024 {
		return nil, fmt.Errorf("Invald BIOS size of %d", len(rawBios))
	}

	return &FDS{
		bios: [8*1024]byte(rawBios),
	}, nil
}

func (fds *FDS) LoadDisk(filename string) error {
	var err error
	fds.disk, err = fdslib.ReadRomFile(filename)
	if err != nil {
		return err
	}
	return nil
}

func (fds *FDS) SetSide(side int) error {
	if side > len(fds.disk.Sides) || side < 0 {
		return fmt.Errorf("Invalid side number: %d (available: %d)", side, len(fds.disk.Sides))
	}

	fds.diskSide = side
	return nil
}

func (fds *FDS) ClearRam() {
	fds.fdsRam = [32*1024]byte{}
	fds.cpuRam = [0x800]byte{}
}

func (fds *FDS) ReadByte(address uint16) uint8 {
	if address < 0x800 {
		return fds.cpuRam[address]
	}
	if address >= 0x6000 && address < 0xE000 {
		return fds.fdsRam[address%fdsRamSize]
	}
	if address >= 0xE000 {
		return fds.bios[address%fdsBiosSize]
	}

	// TODO: audio registers
	switch address {
	case 0x4020: // Timer IRQ reload value low
	case 0x4021: // Timer IRQ reload value high
	case 0x4022: // Timer IRQ control
	case 0x4023: // Master I/O enable
	case 0x4024: // Write data register
	case 0x4025: // FDS Control
	case 0x4026: // External connector
	case 0x4030: // Disk Status register
	case 0x4031: // Read data register
	case 0x4032: // Disk drive status register
	case 0x4033: // External connector read
	}

	return 0
}

func (fds *FDS) WriteByte(address uint16, value uint8) {
	if address < 0x800 {
		fds.cpuRam[address] = value
	}
	if address >= 0x6000 && address < 0xE000 {
		fds.fdsRam[address%fdsRamSize] = value
	}
	if address >= 0xE000 {
		fds.bios[address%fdsBiosSize] = value
	}

	// TODO: audio registers
	switch address {
	case 0x4020: // Timer IRQ reload value low
		fds.irqReload = (fds.irqReload & 0xFF00) | uint16(value)

	case 0x4021: // Timer IRQ reload value high
		fds.irqReload = (fds.irqReload & 0x00FF) | (uint16(value) << 8)

	case 0x4022: // Timer IRQ control
		if !fds.diskIoEnabled {
			return
		}

		fds.irqEnabled = (value & 0x02) == 0x02
		if (value & 0x02) == 0x02 {
			fds.irqCounter = fds.irqReload
		}

		if (value & 0x01) == 0x01 {
			fds.irqRepeat = true
		}

	case 0x4023: // Master I/O enable
		fds.diskIoEnabled  = (value % 0x01) == 0x01
		fds.soundIoEnabled = (value % 0x02) == 0x02

	case 0x4024: // Write data register

	case 0x4025: // FDS Control
		if (value % 0x04) == 0x04 {
			fds.transferMode = fdsRead
		} else {
			fds.transferMode = fdsWrite
		}

		if (value % 0x08) == 0x08 {
			fds.mirrorMode = fdsVertical
		} else {
			fds.mirrorMode = fdsHorizontal
		}

		fds.crcEnabled = (value % 0x10) == 0x10
		fds.irqTransferEnabled = (value % 0x80) == 0x80

	case 0x4026: // External connector

	case 0x4030: // Disk Status register
		// read only

	case 0x4031: // Read data register
		// read only

	case 0x4032: // Disk drive status register
		// read only
	case 0x4033: // External connector read
		// read only
	}
}

func (fds *FDS) FindLabel(name string) (uint, labels.MemoryType) {
	panic("FDS.FindLabel() is not implemented")
}

func (fds *FDS) GetLabel(address uint) string {
	panic("FDS.GetLabel() is not implemented")
}

func (fds *FDS) GetZpLabel(address uint) string {
	panic("FDS.GetZpLabel() is not implemented")
}

func (fds *FDS) AddAutoLabel(address uint16) {
	panic("FDS.AddAutoLabel() is not implemented")
}

func (fds *FDS) Labels(t labels.MemoryType) labels.LabelMap {
	panic("FDS.Labels() is not implemented")
}

func (fds *FDS) Visited(address uint16) bool {
	panic("FDS.Visited() is not implemented")
}

func (fds *FDS) Visit(address uint16) {
	panic("FDS.Visit() is not implemented")
}

func (fds *FDS) AddDasm(address uint16, instruction Instruction, arg uint) {
	panic("FDS.AddDasm() is not implemented")
}

func (fds *FDS) WriteDasm(writer io.Writer) error {
	panic("FDS.WriteDasm() is not implemented")
}
