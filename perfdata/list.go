package perfdata

type PerfdataList []*Perfdata

func (l PerfdataList) String() (s string) {
	for _, p := range l {
		if len(s) > 0 {
			s += " "
		}

		s += p.String()
	}

	return
}

func (l *PerfdataList) Add(p *Perfdata) {
	*l = append(*l, p)
}
