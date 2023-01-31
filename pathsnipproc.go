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
func NewPathSnipProc(path *Path) *PathSnipProc {
	return &PathSnipProc{RenderFlatten, path}
}

// Process implements the PathProcessor interface.
func (psp *PathSnipProc) Process(p *Path) []*Path {
	// Flatten the paths and get the parts
	spparts := psp.Path.Flatten(psp.Flatten).Parts()
	if len(spparts) == 0 {
		return []*Path{p.Copy()}
	}
	pparts := p.Flatten(psp.Flatten).Parts()
	if len(pparts) == 0 {
		return []*Path{p.Copy()}
	}

	paths := [][][][]float64{pparts}
	for _, sppart := range spparts {
		npaths := [][][][]float64{}
		for _, pparts := range paths {
			npparts := [][][]float64{}
			for _, ppart := range pparts {
				splitparts := partsplit(sppart, ppart)
				if splitparts == nil {
					npparts = append(npparts, ppart)
					continue
				}
				npparts = append(npparts, splitparts[0])
				npaths = append(npaths, npparts)
				npparts = [][][]float64{splitparts[1]}
			}
			npaths = append(npaths, npparts)
		}
		paths = npaths
	}

	res := make([]*Path, len(paths))
	for i, parts := range paths {
		res[i] = PartsToPath(parts...)
	}
	return res
}

func partsplit(sppart, ppart [][]float64) [][][]float64 {
	tvals, err := util.IntersectionTValsP(sppart[0], sppart[1], ppart[0], ppart[1])
	if err != nil || tvals[0] < 0 || tvals[0] > 1 || tvals[1] < 0 || tvals[1] > 1 {
		return nil
	}
	// Split ppart into two
	dx, dy := ppart[1][0]-ppart[0][0], ppart[1][1]-ppart[0][1]
	dx *= tvals[1]
	dy *= tvals[1]
	ip := []float64{ppart[0][0] + dx, ppart[0][1] + dy}
	return [][][]float64{{ppart[0], ip}, {ip, ppart[1]}}
}
