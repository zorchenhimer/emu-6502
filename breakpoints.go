package emu

import (
	"bytes"
	"strings"
	"fmt"
)

const (
	READ uint8 = 1
	WRITE uint8 = 2
	EXECUTE uint8 = 4
)

func EventToString(event uint8) string {
	types := []string{}
	if event & READ != 0 {
		types = append(types, "READ")
	}
	if event & WRITE != 0 {
		types = append(types, "WRITE")
	}
	if event & EXECUTE != 0 {
		types = append(types, "EXECUTE")
	}

	return strings.Join(types, "|")
}

type BreakpointCallback func(c *Core, eventType uint8, value uint8)

type Breakpoint struct {
	Type uint8
	Address uint16
	Callback BreakpointCallback
	Name string
}

func (b Breakpoint) String() string {
	return fmt.Sprintf("$%04X [%s] %s", b.Address, EventToString(b.Type), b.Name)
}

type Breakpoints struct {
	registered map[uint16][]Breakpoint
}

func (b *Breakpoints) String() string {
	var out bytes.Buffer
	for _, lst := range b.registered {
		for _, brk := range lst {
			out.WriteString(brk.String())
			out.WriteString("\n")
		}
	}
	return out.String()
}

func (b *Breakpoints) Register(t uint8, name string, address uint16, fn BreakpointCallback) {
	fmt.Printf("Registering %s breakpoint at $%04X\n", EventToString(t), address)
	nbp := Breakpoint{
		Type: t,
		Address: address,
		Name: name,
		Callback: fn,
	}

	if b.registered == nil {
		b.registered = make(map[uint16][]Breakpoint)
	}

	lst, ok := b.registered[address]
	if !ok {
		b.registered[address] = []Breakpoint{nbp}
		return
	}

	nlst := []Breakpoint{}
	found := false

	// check for duplicate breakpoint at this address.
	// Overwrite if found.
	for _, bp := range lst {
		if name == bp.Name {
			nlst = append(nlst, nbp)
			found = true
		} else {
			nlst = append(nlst, bp)
		}
	}

	if !found {
		nlst = append(nlst, nbp)
	}

	b.registered[address] = nlst

	fmt.Println(b.String())
}

func (b *Breakpoints) Read(c *Core, address uint16, value uint8) {
	b.runBreakpoints(c, READ, address, value)
}

func (b *Breakpoints) Write(c *Core, address uint16, value uint8) {
	b.runBreakpoints(c, WRITE, address, value)
}

func (b *Breakpoints) Execute(c *Core, address uint16, value uint8) {
	b.runBreakpoints(c, EXECUTE, address, value)
}

func (b *Breakpoints) runBreakpoints(c *Core, t uint8, address uint16, value uint8) {
	bplst, ok := b.registered[address]
	if !ok {
		return
	}

	for _, bp := range bplst {
		if bp.Type & t == 0 {
			continue
		}
		bp.Callback(c, t, value)
	}
}

func (b *Breakpoints) Clear() {
	b.registered = make(map[uint16][]Breakpoint)
}
