package graphics2d

import (
	"math"
	"math/rand"
)

// JitterProc contains the percentage degree to which segment endpoints will be moved to their left or right (relative
// to their tangents) based on their length. Can be used with MunchProc to get a hand drawn look.
type JitterProc struct {
	Perc float64 // Percentage
}

// Process implements the PathProcessor interface.
func (j *JitterProc) Process(p *Path) []*Path {
	parts := p.Parts()
	np := len(parts)
	if np == 0 {
		return []*Path{p}
	}

	// Path start and end are left unchanged
	res := NewPath(parts[0][0])
	if len(parts[0]) == 1 {
		return []*Path{res}
	}

	for i := 0; i < np-1; i++ {
		part := parts[i]
		end := len(part) - 1
		dx, dy := part[end][0]-part[0][0], part[end][1]-part[0][1]
		l := math.Hypot(dx, dy)
		nx, ny := dy/l, -dx/l
		l *= (rand.Float64()*2 - 1) * j.Perc / 2
		part[end][0] += nx * l
		part[end][1] += ny * l
		res.AddStep(part[1:]...)
	}
	res.AddStep(parts[np-1][1:]...)

	if p.closed {
		res.Close()
	}

	return []*Path{res}
}
