package mappers

type Bank struct {
	// Possible start addresses in CPU or PPU space
	StartAddresses []uint16

	// Will probably need to be a multiple of 8
	Size uint

	Type BankType
}

type BankType string

const (
	PrgRom BankType = "PrgRom"
	PrgRam BankType = "PrgRam"
	ChrRom BankType = "ChrRom"
	ChrRam BankType = "ChrRam"
)

type BankConfiguration struct {
	Banks []*Bank
}

// MMC1 examples

var mmc1_banks = BankConfiguration{
	Banks: []*Bank{
		&Bank{
			StartAddresses: []uint16{0x6000},
			Size: 0x2000,
			Type: PrgRam,
		},

		&Bank{
			StartAddresses: []uint16{0x8000, 0xC000},
			Size: 0x4000,
			Type: PrgRom,
		},

		&Bank{
			StartAddresses: []uint16{0x8000},
			Size: 0x8000,
			Type: PrgRom,
		},

		&Bank{
			StartAddresses: []uint16{0x0000, 0x1000},
			Size: 0x1000,
			Type: ChrRom,
		},
	},
}
