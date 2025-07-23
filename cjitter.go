package graphics2d

import (
	"math"
	"math/rand"
)

// CircularJitterProc takes a path and jitters its internal step points by a random amount within the defined radius.
type CircularJitterProc struct {
	Radius float64
}

// Process implements the PathProcessor interface.
func (sp *CircularJitterProc) Process(p *Path) []*Path {
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
		part[end] = jitter(part[end], sp.Radius)
		res.AddStep(part[1:]...)
	}
	res.AddStep(parts[np-1][1:]...)

	if p.closed {
		res.Close()
	}

	return []*Path{res}
}

func jitter(pt []float64, r float64) []float64 {
	th := rand.Float64() * TwoPi
	dx, dy := math.Cos(th)*r, math.Sin(th)*r
	return []float64{pt[0] + dx, pt[1] + dy}
}
