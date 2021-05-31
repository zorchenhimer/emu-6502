package emu

// Parse ca65's debugging symbols file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type SymbolRecord struct {
	Name string
	Size uint16 // size of the data (the '.res #' value)
	AddrSize int // size of the address itself.  either zero page or absolute
	Value uint16
	Type string

	Defined *lineRecord

	References []*lineRecord
}

type fileRecord struct {
	Id int
	ModuleId int

	Name string
	Size int
	Mtime int // what is the format of this?
}

type Symbols struct {
	Version string

	files map[int]*fileRecord
	lines map[int]*lineRecord

	// indexed by scope name.  "" is default
	// "SomeLabel"
	// "SomeScope::SomeLabel"
	// "ParentLabel@localLabel"
	sym map[string]*SymbolRecord

	// Each address may have more than one symbol attached to it
	symAddr map[uint16][]*SymbolRecord

	segments map[int]*segmentRecord
}

type segmentRecord struct {
	Id int
	Name string
	Start int
	Size uint16
	Type string
	OutputName string
	OutputOffset int
}

func (s *Symbols) AllLabels() []string {
	l := []string{}
	for key, _ := range s.sym {
		l = append(l, key)
	}

	return l
}

// Returns a list of labels that occupy a given address.
// Returns nil if none found.
func (s *Symbols) LabelsAt(address uint16) ([]string) {
	if list, ok := s.symAddr[address]; ok {
		ret := []string{}
		for _, rec := range list {
			ret = append(ret, rec.Name)
		}
		return ret
	}

	return nil
}

// GetAddress("SomeLabel")
// GetAddress("SomeScope::SomeLabel")
// GetAddress("ParentLabel@localLabel")
func (s *Symbols) GetAddress(name string) (uint16, error) {
	sym, err := s.GetSymbol(name)
	if err != nil {
		return 0, err
	}
	return sym.Value, nil
}

// Returns the full symbol object
func (s *Symbols) GetSymbol(name string) (*SymbolRecord, error) {
	val, ok := s.sym[name]
	if ok {
		return val, nil
	}

	return nil, fmt.Errorf("Address symbol %q does not exist", name)
}

type lineRecord struct {
	Id int
	File *fileRecord
	Line int
	Span *spanRecord
}

type spanRecord struct {
}

func NewSymbols(filename string) (*Symbols, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scopes := map[int]map[string]string{}
	symbols := map[int]map[string]string{}
	lines := map[int]map[string]string{}

	sym := &Symbols{
		sym: map[string]*SymbolRecord{},
		symAddr: map[uint16][]*SymbolRecord{},
		Version: "",
		files: map[int]*fileRecord{},
		lines: map[int]*lineRecord{},
	}

	// pass one
	reader := bufio.NewReader(file)
	var readerr error
	for readerr == nil {
		var line string
		line, readerr = reader.ReadString('\n')
		if readerr != nil && readerr != io.EOF {
			return nil, readerr
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		idx := strings.Index(line, "\t")
		if idx == -1 {
			return nil, fmt.Errorf("Invalid line: %q", line)
		}

		t := line[:idx]
		vals := strings.Split(line[idx+1:], ",")

		m := map[string]string{}
		for _, kv := range vals {
			data := strings.Split(kv, "=")
			if len(data) != 2 {
				return nil, fmt.Errorf("Invalid key/value pair: %q", kv)
			}

			m[data[0]] = strings.Trim(data[1], `"`)
		}

		idStr, ok := m["id"]
		if !ok {
			continue
		}

		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			return nil, err
		}

		switch t {
		case "scope":
			scopes[int(id)] = m

		case "sym":
			symbols[int(id)] = m

		case "line":
			lines[int(id)] = m

		case "version":
			sym.Version = m["major"] + "." + m["minor"]

		case "seg":
			id, err := strconv.Atoi(m["id"])
			if err != nil {
				return nil, fmt.Errorf("Unable to parse id value for segment %q: %w", m["id"], err)
			}

			start, err := strconv.ParseInt(m["start"], 0, 32)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse start value for segment %q: %w", m["start"], err)
			}

			size, err := strconv.ParseUint(m["size"], 0, 16)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse size value for segment %q: %w", m["size"], err)
			}

			seg := &segmentRecord{
				Id: id,
				Name: m["name"],
				Start: int(start),
				Size: uint16(size),
				Type: m["type"],
			}

			var oname string
			if val, ok := m["oname"]; ok {
				oname = val
				offset, err := strconv.Atoi(m["ooffs"])
				if err != nil {
					return nil, fmt.Errorf("Unable to parse ooffs value for segment %q: %w", m["ooffs"], err)
				}

				seg.OutputName = oname
				seg.OutputOffset = offset
			}

		case "file":
			mtime, err := strconv.ParseInt(m["mtime"], 0, 32)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse mtime value of %q for file record: %w", m["mtime"], err)
			}

			size , err := strconv.ParseInt(m["size"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse size value of %q for file record: %w", m["size"], err)
			}

			sym.files[int(id)] =  &fileRecord{
				Id: int(id),
				Name: m["name"],
				Mtime: int(mtime),
				Size: int(size),
			}


		default:
			// all others not implemented yet
			continue
		}
	}

	// second passes
	for _, line := range lines {
		id, err := strconv.Atoi(line["id"])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse line ID %q: %w", err)
		}

		fileId, err := strconv.Atoi(line["file"])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse line fileId %q: %w", err)
		}

		lineNum, err := strconv.Atoi(line["line"])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse line number %q: %w", err)
		}

		f, ok := sym.files[fileId]
		if !ok {
			return nil, fmt.Errorf("Cannot find file with ID %d", fileId)
		}

		sym.lines[id] = &lineRecord{
			Id: id,
			Line: lineNum,
			File: f,
			// TODO: span
		}
	}

	for _, symbol := range symbols {
		prefix := ""

		var parent map[string]string
		var parentOk bool = false

		parentIdStr, ok := symbol["parent"]
		if ok {
			parentId, err := strconv.ParseInt(parentIdStr, 10, 32)
			if err != nil {
				return nil, err
			}

			parent, parentOk = symbols[int(parentId)]
		}

		var scopeIdStr string
		if parentOk {
			scopeIdStr = parent["scope"]
			prefix = parent["name"]
		} else {
			scopeIdStr = symbol["scope"]
		}

		scopeId, err := strconv.ParseInt(scopeIdStr, 10, 32)
		if err != nil {
			return nil, err
		}

		scope, ok := scopes[int(scopeId)]
		if !ok {
			fmt.Printf("Scope with ID %d not found", scopeId)
			continue
		}

		if scope["name"] != "" {
			prefix = scope["name"] + "::" + prefix
		}

		var addrSize int = 0
		switch symbol["addrsize"] {
		case "absolute":
			addrSize = 2
		case "zeropage":
			addrSize = 1
		default:
			fmt.Println("unknown addrsize:", symbol["addrsize"])
			continue
		}

		size := uint16(1)
		if val, ok := symbol["size"]; ok {
			s, err := strconv.ParseUint(val, 0, 16)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse size value %q: %w", val, err)
			}
			size = uint16(s)
		}

		addrValue, err := strconv.ParseUint(symbol["val"], 0, 16)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse address value %q: %w", symbol["val"], err)
		}

		name := symbol["name"]
		if prefix != "" {
			name = prefix+name
		}

		def := symbol["def"]
		plusIdx := strings.LastIndex(def, "+")
		if plusIdx > -1 {
			def = def[plusIdx:]
		}

		definedId, err := strconv.Atoi(def)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse def value for symbol %q: %w", symbol["def"], err)
		}

		defined, ok := sym.lines[definedId]
		if !ok {
			return nil, fmt.Errorf("No line with id %d", definedId)
		}

		record := &SymbolRecord{
			Name: name,
			AddrSize: addrSize,
			Size: size,
			Value: uint16(addrValue),
			Type: scope["type"],
			Defined: defined,
			References: []*lineRecord{},
		}
		sym.sym[name] = record

		if references, ok := symbol["ref"]; ok {
			refList := strings.Split(references, "+")
			for _, lineIdStr := range refList {
				lineId, err := strconv.Atoi(lineIdStr)
				if err != nil {
					return nil, fmt.Errorf("Unable to parse reference id %q: %w", lineIdStr, err)
				}

				if line, lok := sym.lines[lineId]; lok {
					record.References = append(record.References, line)
				}
			}
		}

		// Add a reference for every address this label occupies
		for i := uint16(0); i < size; i++ {
			if _, ok := sym.symAddr[uint16(addrValue)+i]; !ok {
				sym.symAddr[uint16(addrValue)+i] = []*SymbolRecord{}
			}

			sym.symAddr[uint16(addrValue)+i] = append(sym.symAddr[uint16(addrValue)+i], record)
		}

	}

	//for name, _ := range sym.sym {
	//	fmt.Println(name)
	//}

	return sym, nil
}
