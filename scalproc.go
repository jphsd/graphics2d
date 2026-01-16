package graphics2d

import "math"

// ScallopProc applies a scallop along a path with a defined length.
type ScallopProc struct {
	Lambda float64 // Scallop length
	Flip   bool    // Flips the arc from the right side to the left of the path
}

// Process implements the PathProcessor interface.
func (sp ScallopProc) Process(p *Path) []*Path {
	// Chunk up path into pieces
	pp1 := NewMunchProc(sp.Lambda).Process(p)
	p1, _ := ConcatenatePaths(pp1...)
	var arc *Path
	if sp.Flip {
		arc = PartsToPath(MakeArcParts(0, 0, sp.Lambda/2, Pi, -Pi)...)
	} else {
		arc = PartsToPath(MakeArcParts(0, 0, sp.Lambda/2, Pi, Pi)...)
	}

	n := len(p1.steps)
	last := p1.steps[0][0]
	path := NewPath(last)
	for i := 1; i < n; i++ {
		// Each step is a scallop
		cur := p1.steps[i][0]
		dx, dy := cur[0]-last[0], cur[1]-last[1]
		th := math.Atan2(dy, dx)
		cx, cy := (cur[0]+last[0])/2, (cur[1]+last[1])/2
		xfm := Translate(cx, cy)
		xfm.Rotate(th)
		path.Concatenate(arc.Process(xfm)...)
		last = cur
	}

	if p.Closed() {
		path.Close()
	}

	return []*Path{path}
}
