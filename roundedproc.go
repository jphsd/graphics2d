package graphics2d

// RoundedProc replaces adjacent line segments in a path with line-arc-line where the radius of the
// arc is the minimum of Radius or the maximum allowable for the length of the shorter line segment.
// This ensures that the rounded corner doesn't end beyond the mid point of either line.
type RoundedProc struct {
	Radius float64
}

// Process implements the PathProcessor interface.
func (rp RoundedProc) Process(p *Path) []*Path {
	parts := p.Parts()
	np := len(parts)
	if np < 2 {
		return []*Path{p}
	}

	res := [][][]float64{}
	for i, part := range parts {
		if len(part) != 2 {
			res = append(res, part)
			continue
		}
		if i < np-1 {
			if len(parts[i+1]) == 2 {
				mp1, mp2 := Lerp(0.5, part[0], part[1]), Lerp(0.5, part[1], parts[i+1][1])
				nparts := MakeRoundedParts(mp1, part[1], mp2, rp.Radius)
				res = append(res, [][]float64{part[0], nparts[0][0]})
				res = append(res, nparts...)
			} else {
				res = append(res, part)
			}
		} else {
			if !p.Closed() || len(parts[0]) != 2 {
				res = append(res, part)
				continue
			}
			// Path is closed and the first part is also a line
			mp1, mp2 := Lerp(0.5, part[0], part[1]), Lerp(0.5, part[1], parts[0][1])
			nparts := MakeRoundedParts(mp1, part[1], mp2, rp.Radius)
			res = append(res, [][]float64{part[0], nparts[0][0]})
			res = append(res, nparts...)
			lnp := len(nparts) - 1
			res[0][0] = nparts[lnp][len(nparts[lnp])-1]
		}
	}

	path := PartsToPath(res...)
	if p.Closed() {
		path.Close()
	}
	return []*Path{path}
}
