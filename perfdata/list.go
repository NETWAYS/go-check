package perfdata

import (
	"strings"
)

// PerfdataList can store multiple perfdata and brings a simple fmt.Stringer interface
type PerfdataList []*Perfdata

// String returns string representations of all Perfdata
func (l PerfdataList) String() string {
	var out strings.Builder

	for _, p := range l {
		if len(out.String()) > 0 {
			out.WriteString(" ")
		}

		out.WriteString(p.String())
	}

	return out.String()
}

// Add a Perfdata to the list
func (l *PerfdataList) Add(p *Perfdata) {
	*l = append(*l, p)
}
