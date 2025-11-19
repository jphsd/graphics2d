package graphics2d

import (
	"github.com/jphsd/graphics2d/util"
)

// PathSnipProc contains the snip path.
type PathSnipProc struct {
	Flatten float64
	Path    *Path
}

// NewPathSnipProc creates a new path snip path processor with the supplied path.
func NewPathSnipProc(path *Path) PathSnipProc {
	return PathSnipProc{RenderFlatten, path}
}

// Process implements the PathProcessor interface.
func (psp PathSnipProc) Process(p *Path) []*Path {
	// Flatten the paths and get the parts
	spparts := psp.Path.Flatten(psp.Flatten).Parts()
	if len(spparts) == 0 {
		return []*Path{p.Copy()}
	}
	pparts := p.Flatten(psp.Flatten).Parts()
	if len(pparts) == 0 {
		return []*Path{p.Copy()}
	}

	paths := [][]Part{pparts}
	for _, sppart := range spparts {
		npaths := [][]Part{}
		for _, pparts := range paths {
			npparts := []Part{}
			for _, ppart := range pparts {
				splitparts := partsplit(sppart, ppart)
				if splitparts == nil {
					npparts = append(npparts, ppart)
					continue
				}
				if splitparts[0] != nil {
					npparts = append(npparts, splitparts[0])
				}
				if len(npparts) != 0 {
					npaths = append(npaths, npparts)
				}
				if splitparts[1] != nil {
					npparts = []Part{splitparts[1]}
				} else {
					npparts = []Part{}
				}
			}
			if len(npparts) != 0 {
				npaths = append(npaths, npparts)
			}
		}
		paths = npaths
	}

	res := make([]*Path, len(paths))
	for i, parts := range paths {
		res[i] = PartsToPath(parts...)
	}
	return res
}

func partsplit(sppart, ppart Part) []Part {
	// TODO - is this faster than using PartsIntersection() which uses BB?
	tvals, err := util.IntersectionTValsP(sppart[0], sppart[1], ppart[0], ppart[1])
	if err != nil || tvals[0] < 0 || tvals[0] > 1 || tvals[1] < 0 || tvals[1] > 1 {
		return nil
	}
	t := tvals[1]
	if t > 0 && t < 1 {
		// Split ppart into two
		dx, dy := ppart[1][0]-ppart[0][0], ppart[1][1]-ppart[0][1]
		dx *= t
		dy *= t
		ip := []float64{ppart[0][0] + dx, ppart[0][1] + dy}
		return []Part{{ppart[0], ip}, {ip, ppart[1]}}
	}
	if t < 1 {
		// t is at start of part
		return []Part{nil, ppart}
	}
	// t is at end of part
	return []Part{ppart, nil}
}
