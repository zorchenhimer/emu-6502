package dnasm

import (
	"fmt"
)

const (
	VECTOR_NMI   uint16 = 0xFFFA
	VECTOR_RESET uint16 = 0xFFFC
	VECTOR_IRQ   uint16 = 0xFFFE
)

func Disassemble(rom []byte) (*Disassembly, error) {

	d, err := New(rom)
	if err != nil {
		return nil, err
	}

	fmt.Println("[irq]")
	vIrq := d.m.ReadWord(VECTOR_IRQ)
	fmt.Println("[nmi]")
	vNmi := d.m.ReadWord(VECTOR_NMI)
	fmt.Println("[reset]")
	vReset := d.m.ReadWord(VECTOR_RESET)

	//fmt.Println("")

	//fmt.Printf("%s\nIRQ: %04X (%08X)\nNMI: %04X (%08X)\nRESET: %04X (%08X)\n",
	//	mapper.State(),
	//	vIrq, mapper.Offset(VECTOR_IRQ),
	//	vNmi, mapper.Offset(VECTOR_NMI),
	//	vReset, mapper.Offset(VECTOR_RESET),
	//)

	irq := true
	nmi := true

	if vIrq == 0x0000 || vIrq < 0x8000 || vIrq == vReset || vIrq == vNmi {
		irq = false
	}

	if vNmi == 0x0000 || vNmi < 0x8000 || vNmi == vReset {
		nmi = false
	}

	//fmt.Println("\n[NMI]")
	//if nmi {
	//	err = d.process(vNmi)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	//fmt.Println("\n[IRQ]")
	//if irq {
	//	if err = d.process(vIrq); err != nil {
	//		return nil, err
	//	}
	//}

	//return nil, fmt.Errorf("Stopping before processing RESET vector")

	//fmt.Println("\n[RESET]")
	//err = d.process(vReset)
	//fmt.Printf("process error: %v\n", err)

	d.AddVector(vReset)

	if nmi {
		d.AddVector(vNmi)
	}

	if irq {
		d.AddVector(vIrq)
	}

	//for offset, _ := range d.tokens {
	//	fmt.Printf("  %08X\n", offset)
	//}

	//for chunk, _ := range d.chunks {
	//	fmt.Printf("> %08X\n", chunk)
	//}

	//fmt.Printf("chunks processed: %d\n", d.processed)

	return d, d.process()
}


