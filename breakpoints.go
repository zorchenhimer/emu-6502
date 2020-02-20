package emu

type BreakpointType int

const (
	READ BreakpointType = iota
	WRITE
	EXECUTE
)

type BreakpointCallback func(c *Core, eventType BreakpointType, value uint8)

type Breakpoint struct {
	Type BreakpointType
	Address uint16
	Callback BreakpointCallback
	Name string
}

type Breakpoints struct {
	registered map[uint16][]Breakpoint
}

func (b *Breakpoints) Register(t BreakpointType, name string, address uint16, fn BreakpointCallback) {
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
		lst = []Breakpoint{nbp}
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

func (b *Breakpoints) runBreakpoints(c *Core, t BreakpointType, address uint16, value uint8) {
	bplst, ok := b.registered[address]
	if !ok {
		return
	}

	for _, bp := range bplst {
		if bp.Type != t {
			continue
		}
		bp.Callback(c, t, value)
	}
}

func (b *Breakpoints) Clear() {
	b.registered = make(map[uint16][]Breakpoint)
}
