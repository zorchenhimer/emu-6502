package emu

func TwosCompInv(value uint8) (uint8, bool) {
	if value&0x80 != 0 {
		return (value ^ 0xFF) + 1, true
	}
	return value, false
}

func PadWithVectors(rom []byte, nmi, reset, irq uint16) []byte {
	for len(rom)%256 != 0 {
		rom = append(rom, 0xFF)
	}

	addr := len(rom) - 6

	rom[addr] = byte(nmi & 0x00FF)
	rom[addr+1] = byte(nmi >> 8)

	rom[addr+2] = byte(reset & 0x00FF)
	rom[addr+3] = byte(reset >> 8)

	rom[addr+4] = byte(irq & 0x00FF)
	rom[addr+5] = byte(irq >> 8)

	return rom
}

func FlagsToString(ph uint8) string {
	sc := "-"
	sz := "-"
	si := "-"
	sd := "-"
	sv := "-"
	sn := "-"

	if ph&FLAG_CARRY != 0 {
		sc = "C"
	}

	if ph&FLAG_ZERO != 0 {
		sz = "Z"
	}

	if ph&FLAG_INTERRUPT != 0 {
		si = "I"
	}

	if ph&FLAG_DECIMAL != 0 {
		sd = "D"
	}

	if ph&FLAG_OVERFLOW != 0 {
		sv = "V"
	}

	if ph&FLAG_NEGATIVE != 0 {
		sn = "N"
	}

	return fmt.Sprintf("%s%s--%s%s%s%s", sn, sv, sd, si, sz, sc)
}

