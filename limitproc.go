package graphics2d

import (
	"github.com/jphsd/graphics2d/util"
	"math"
)

// LimitProc takes a path and chunks it up into at most Limit length parts.
// Smaller parts are left alone. Similar to FSnipProc but without the pattern.
type LimitProc struct {
	Limit float64 // Maximum allowed length of any path part
}

// Process implements the PathProcessor interface.
func (lp LimitProc) Process(p *Path) []*Path {
	paths := []*Path{}
	for _, part := range p.Parts() {
		n := len(part)
		p1, p2 := part[0], part[n-1]
		dx, dy := p2[0]-p1[0], p2[1]-p1[1]
		l := math.Hypot(dx, dy)
		if l < lp.Limit {
			paths = append(paths, PartsToPath(part))
		} else {
			// Chop it up
			nf := l / lp.Limit
			n := int(nf)
			if !util.Equals(l-float64(n)*lp.Limit, 0) {
				// Round up
				n++
			}
			for i := range n {
				t := 1 / float64(n-i)
				p3 := Lerp(t, p1, p2)
				paths = append(paths, PartsToPath([][]float64{p1, p3}))
				p1 = p3
			}
		}
	}
	return paths
}
