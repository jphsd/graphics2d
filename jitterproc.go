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
func (j JitterProc) Process(p *Path) []*Path {
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

	for i := range np {
		part := parts[i]
		end := len(part) - 1
		dx, dy := part[end][0]-part[0][0], part[end][1]-part[0][1]
		l := math.Hypot(dx, dy)
		nx, ny := dy/l, -dx/l
		if p.closed && i == 0 {
			// Jitter first if closed
			dl := l * (rand.Float64()*2 - 1) * j.Perc / 2
			part[0][0] += nx * dl
			part[0][1] += ny * dl
			res = NewPath(part[0])
		}
		if i != np-1 || p.closed {
			dl := l * (rand.Float64()*2 - 1) * j.Perc / 2
			part[end][0] += nx * dl
			part[end][1] += ny * dl
			res.AddStep(part[1:]...)
		}
	}

	if p.closed {
		res.Close()
	}

	return []*Path{res}
}

// CircularJitterProc takes a path and jitters its internal step points by a random amount within the defined radius.
type CircularJitterProc struct {
	Radius float64
}

// Process implements the PathProcessor interface.
func (sp CircularJitterProc) Process(p *Path) []*Path {
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

	for i := range np {
		part := parts[i]
		end := len(part) - 1
		if p.closed && i == 0 {
			part[0] = rjitter(part[0], sp.Radius)
			res = NewPath(part[0])
		}
		if i != np-1 || p.closed {
			part[end] = rjitter(part[end], sp.Radius)
		}
		res.AddStep(part[1:]...)
	}

	if p.closed {
		res.Close()
	}

	return []*Path{res}
}

func rjitter(pt []float64, r float64) []float64 {
	th := rand.Float64() * TwoPi
	dx, dy := math.Cos(th)*r, math.Sin(th)*r
	return []float64{pt[0] + dx, pt[1] + dy}
}
