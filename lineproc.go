package graphics2d

// StepsToLinesProc takes a path and converts all of the points to lines.
type StepsToLinesProc struct {
	IncCP bool
}

// Process implements the PathProcessor interface.
func (clp StepsToLinesProc) Process(p *Path) []*Path {
	parts := p.Parts()
	nparts := [][][]float64{}

	cp := parts[0][0]
	for _, part := range parts {
		if !clp.IncCP {
			lp := part[len(part)-1]
			nparts = append(nparts, [][]float64{cp, lp})
			cp = lp
		} else {
			for i := 1; i < len(part); i++ {
				nparts = append(nparts, [][]float64{cp, part[i]})
				cp = part[i]
			}
		}
	}
	path := PartsToPath(nparts...)
	if p.Closed() {
		path.Close()
	}
	return []*Path{path}
}
