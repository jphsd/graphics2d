package graphics2d

// CurvesToLinesProc takes a path and converts all of the points to lines.
type CurvesToLinesProc struct {
	SkipCP bool
}

// Process implements the PathProcessor interface.
func (clp *CurvesToLinesProc) Process(p *Path) []*Path {
	parts := p.Parts()
	nparts := [][][]float64{}

	cp := parts[0][0]
	for _, part := range parts {
		if clp.SkipCP {
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
	path, _ := PartsToPath(nparts...)
	return []*Path{path}
}

// SplitProc breaks up a path into a collection of paths, one for each step in the original path.
type SplitProc struct{}

// Process implements the PathProcessor interface.
func (sp *SplitProc) Process(p *Path) []*Path {
	parts := p.Parts()
	n := len(parts)
	res := make([]*Path, n)
	for i := 0; i < n; i++ {
		res[i], _ = PartsToPath(parts[i])
	}
	return res
}
