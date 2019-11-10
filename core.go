package emu

import (
    "fmt"
    "testing"
)

const (
    VECTOR_NMI uint16 = 0xFFFA
    VECTOR_RESET uint16 = 0xFFFC
    VECTOR_IRQ uint16 = 0xFFFE
)

type Core struct {
    // Main registers
    A   uint8
    X   uint8
    Y   uint8

    // Other registers
    PC  uint16   // Program counter
    Phlags uint8 // Status flags
    SP  uint8    // Stack pointer

    memory []byte    // Slice of loaded memory.  This is only main RAM.
    rom []byte      // ROM image.  Needs to be a multiple of 256.
    wram []byte

    InstructionLimit uint64 // number of instructions to run
    testing bool
    testDone bool
    t *testing.T
}

func NewCore(rom []byte, wram bool, instrLimit uint64) (*Core, error) {
    if len(rom) % 256 != 0 {
        return nil, fmt.Errorf("ROM is not divisible by 256: %d", len(rom))
    }

    c := &Core{
        A: 0,
        X: 0,
        Y: 0,
        PC: 0,
        Phlags: 0,
        SP: 0,

        memory: make([]byte, 0x1000), // no registers, no WRAM, no ROM
        rom: rom,

        InstructionLimit: instrLimit,
    }

    if wram {
        c.wram = make([]byte, 0x2000)
    }

    c.PC = c.ReadWord(VECTOR_RESET)

    return c, nil
}

// Read address.  This will read from API registers if needed.
func (c *Core) ReadByte(addr uint16) uint8 {
    if addr < 0x1000 {
        return c.memory[addr]
    }

    if addr >= 0x6000 && addr < 0x8000 {
        if c.wram != nil {
            // TODO: make sure this works with variable WRAM sizes (paging?)
            return c.wram[addr % uint16(len(c.wram))]
        }
        return 0
    }

    if addr >= 0x8000 {
        return c.rom[addr % uint16(len(c.rom))]
    }

    // "Open bus"  always return zero.
    return 0
}

func (c *Core) ReadWord(addr uint16) uint16 {
    return uint16(c.ReadByte(addr)) | (uint16(c.ReadByte(addr+1)) << 8)
}

// Write to an address.  This will delegate to API if needed.
func (c *Core) WriteByte(addr uint16, value byte) {
    if addr < 0x1000 {
        c.memory[addr] = value
    } else if addr < 0x6000 {
        // TODO: software registers
    } else if addr >= 0x6000 && addr < 0x8000 && c.wram != nil {
        c.wram[addr] = value
    }
}

func (c *Core) WriteInt(addr uint16, value uint8) {
    c.WriteByte(addr, byte(value))
}

func (c *Core) Run() error {

    limit := false
    if c.InstructionLimit > 0 {
        //fmt.Printf("Setting instruction limit to %d\n", c.InstructionLimit)
        limit = true
    }

    done := false
    var err error
    for !done {
        err = c.tick()
        if err != nil {
            return err
        }

        if limit {
            c.InstructionLimit -= 1
            if c.InstructionLimit <= 0 {
                if c.testing {
                    return fmt.Errorf("Instruction limit hit")
                }
                done = true
            }
        }

        if c.testing {
            done = c.testDone
        }
    }

    return nil
}

func (c *Core) tick() error {
    opcode := c.ReadByte(c.PC)

    switch opcode {

    case OP_JMP_AB:
        c.PC = c.ReadWord(c.PC + 1)
    case OP_JMP_ID:
        c.PC = c.ReadWord(c.ReadWord(c.PC + 1))

        //LDA
    case OP_LDA_IM:
        c.A = c.ReadByte(c.PC + 1)
        c.PC += 2
    case OP_LDA_AB:
        c.A = c.ReadByte(c.ReadWord(c.PC + 1))
        c.PC += 3
    case OP_LDX_IM:
        c.X = c.ReadByte(c.PC + 1)
        c.PC += 2
    case OP_LDA_AX:
        c.A = c.ReadByte(c.ReadWord(c.PC + 1) + uint16(c.X))
        c.PC += 3
    case OP_LDA_AY:
        c.A = c.ReadByte(c.ReadWord(c.PC + 1) + uint16(c.Y))
        c.PC += 3
    case OP_LDA_ZP:
        c.A = c.ReadByte(uint16(c.ReadByte(c.PC + 1)))
        c.PC += 2

    case OP_NOP:
        c.PC += 1

    case OP_STA_AB:
        addr := c.ReadWord(c.PC + 1)
        c.memory[addr] = c.A
        c.PC += 3
    case OP_STA_ZP:
        addr := uint16(c.ReadByte(c.PC + 1))
        c.memory[addr] = c.A
        c.PC += 2
    case OP_STA_ZX:
        addr8 := c.ReadByte(c.PC + 1) + c.X
        c.memory[addr8] = c.A
        c.PC += 2

    default:
        if opcode == 0xFF && c.testing {
            c.testDone = true
            return nil // 0xFF means end of test
        }
        return fmt.Errorf("OpCode $%02X not implemented", opcode)
    }

    return nil
}

func (c *Core) DumpMemoryRange(start, end uint16) {
    if end < start {
        fmt.Println("Invalid dump range given")
        return
    }

    fmt.Printf("start: $%02X end: $%02X\n", start, end)

    vals := []byte{}
    current := start

    for current <= end {
        vals = append(vals, c.ReadByte(current))
        current++
    }

    for i, b := range vals {
        fmt.Printf("$%02X: $%02X (%d)\n", i + int(start), b, b)
    }
}

func (c *Core) tlog(msg string) {
    if c.t != nil {
        c.t.Log(msg)
    }
}

func (c *Core) tlogf(fmt string, args... interface{}) {
    if c.t != nil {
        c.t.Logf(fmt, args...)
    }
}

func testCore(rom []byte, mem []byte, wram []byte) (*Core, error) {
    core, err := NewCore(rom, false, 1000)
    if err != nil {
        return nil, err
    }
    core.testing = true

    if mem != nil {
        for len(mem) < 0x1000 {
            mem = append(mem, 0x00)
        }
        core.memory = mem
    }

    if wram != nil {
        core.wram = wram
    }

    return core, core.Run()
}

func padWithVectors(rom []byte, nmi, reset, irq uint16) []byte {
    for len(rom) % 256 != 0 {
        rom = append(rom, 0xFF)
    }

    addr := len(rom) - 6

    rom[addr] = byte(nmi & 0x00FF)
    rom[addr+1] = byte(nmi >> 8 )

    rom[addr+2] = byte(reset & 0x00FF)
    rom[addr+3] = byte(reset >> 8 )

    rom[addr+4] = byte(irq & 0x00FF)
    rom[addr+5] = byte(irq >> 8 )

    return rom
}
