package perfdata

// PerfdataList can store multiple perfdata and brings a simple fmt.Stringer interface
type PerfdataList []*Perfdata

// String returns string representations of all Perfdata
func (l PerfdataList) String() (s string) {
	for _, p := range l {
		if len(s) > 0 {
			s += " "
		}

		s += p.String()
	}

	return
}

// Add a Perfdata to the list
func (l *PerfdataList) Add(p *Perfdata) {
	*l = append(*l, p)
}
