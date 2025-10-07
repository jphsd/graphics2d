package graphics2d

import (
	"github.com/jphsd/graphics2d/util"
	"math"
)

// VWTraceProc path processor creates a variable width trace of a path
// using a function to determine the length of a line from the point
// where two path parts meet, at an angle that bisects the angle between the two parts.
// The line end points are used to create path that will always be open.
type VWTraceProc struct {
	Width   float64                        // Distance from path
	Flatten float64                        // See Path.Flatten
	Func    func(float64, float64) float64 // Func(t, Width) where t [0,1]
}

// Process implements the path processor interface.
func (s VWTraceProc) Process(p *Path) []*Path {
	// Flatten
	parts := p.Parts()
	nparts := make([][][]float64, 0, len(parts))
	for _, part := range parts {
		nparts = append(nparts, FlattenPart(s.Flatten, part)...)
	}
	parts = nparts

	// Calculate part lengths
	np := len(parts)
	sum := 0.0
	plens := make([]float64, len(parts))
	for i, part := range parts {
		plens[i] = util.DistanceE(part[0], part[1])
		sum += plens[i]
	}

	points := make([][]float64, 0, np+1)
	clen := 0.0
	for i, part := range parts {
		clen += plens[i]
		if i == np-1 {
			// Last part...
			if p.Closed() {
				npart := parts[0]
				a, a1, _ := util.AngleBetweenLines(part[0], part[1], part[1], npart[1])
				var ba float64
				if a < 0 {
					ba = (-Pi - a) / 2
				} else {
					ba = (Pi - a) / 2
				}
				th := a1 - Pi - ba
				dx, dy := math.Cos(th), math.Sin(th)
				w0 := s.Func(0, s.Width)
				w1 := s.Func(1, s.Width)
				if ba < 0 {
					w0 = -w0
					w1 = -w1
				}
				points = append(points, []float64{part[1][0] + w1*dx, part[1][1] + w1*dy})
				points = append(points, []float64{part[1][0] + w0*dx, part[1][1] + w0*dy})
				// Extra line will be provided by Close()
			} else {
				dx, dy := part[1][0]-part[0][0], part[1][1]-part[0][1]
				d := math.Hypot(dx, dy)
				dx /= d
				dy /= d
				nx, ny := -dy, dx
				w := s.Func(clen/sum, s.Width)
				points = append(points, []float64{part[1][0] + w*nx, part[1][1] + w*ny})
			}
			break
		}
		npart := parts[i+1]
		a, a1, _ := util.AngleBetweenLines(part[0], part[1], part[1], npart[1]) // [-Pi,Pi]
		var ba float64                                                          // bisection angle
		if a < 0 {
			ba = (-Pi - a) / 2
		} else {
			ba = (Pi - a) / 2
		}
		th := a1 - Pi - ba
		dx, dy := math.Cos(th), math.Sin(th)
		w := s.Func(clen/sum, s.Width)
		if ba < 0 {
			w = -w
		}
		points = append(points, []float64{part[1][0] + w*dx, part[1][1] + w*dy})
	}
	if !p.Closed() {
		// Figure first point and construct path
		part := parts[0]
		dx, dy := part[1][0]-part[0][0], part[1][1]-part[0][1]
		d := math.Hypot(dx, dy)
		dx /= d
		dy /= d
		nx, ny := -dy, dx
		w := s.Func(0, s.Width)
		path := NewPath([]float64{part[0][0] + w*nx, part[0][1] + w*ny})
		for _, pt := range points {
			path.AddStep(pt)
		}
		return []*Path{path}
	}
	npts := len(points)
	path := NewPath(points[npts-1])
	for i, pt := range points {
		if i == npts-1 {
			break
		}
		path.AddStep(pt)
	}
	return []*Path{path}
}
