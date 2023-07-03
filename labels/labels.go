package labels

import (
	"fmt"
	"os"

	"github.com/zorchenhimer/go-nes/mesen"
)

type Label struct {
	Name string
	Comment string
	Size uint
}

type LabelMap map[uint]*Label

// Return the label for the given address.  If there is an exact match, return
// that label.  Otherwise, look for the closest previous label and return that
// label if the given address is within the Size of that label.
func (lm LabelMap) FindLabel(address uint) string {
	if lbl, ok := lm[address]; ok {
		if lbl.Size > 1 {
			return lbl.Name+"+0"
		}
		return lbl.Name
	}

	max := uint(0)
	var maxlbl *Label
	for addr, lbl := range lm {
		if addr > address || addr < max {
			continue
		}

		if addr >= max {
			max = addr
			maxlbl = lbl
		}
	}

	if maxlbl != nil && (max + maxlbl.Size) > address {
		return fmt.Sprintf("%s+%d", maxlbl.Name, (address - max))
	}

	return ""
}

func LoadMesen2(filename string) (map[mesen.MemoryType]LabelMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", filename, err)
	}
	defer file.Close()

	ws, err := mesen.LoadWorkspace(file)
	if err != nil {
		return nil, fmt.Errorf("unable to load workspace: %w", err)
	}

	ret := make(map[mesen.MemoryType]LabelMap)

	for _, lbl := range ws.Labels {
		if _, ok := ret[mesen.MemoryType(lbl.MemoryType)]; !ok {
			ret[mesen.MemoryType(lbl.MemoryType)] = make(LabelMap)
		}
		ret[mesen.MemoryType(lbl.MemoryType)][uint(lbl.Address)] = &Label{Name: lbl.Label, Comment: lbl.Comment, Size: uint(lbl.Length)}
	}

	return ret, nil
}
