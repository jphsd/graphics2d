package graphics2d

// SquareWaveProc applies a square wave along a path with a defined wave length and amplitude.
// The wave starts and ends on zero-crossing points and the last half wave is truncated to the
// path length remaining. The internal zero-crossing points can be optionally preserved.
type SquareWaveProc struct {
	HalfLambda float64 // Half wave length
	Scale      float64 // Ratio of amplitude to lambda
	KeepZero   bool    // Keeps internal zero-point crossings if set
	Flip       bool    // Flips the wave phase by 180 if set
}

// NewSquareWaveProc creates a new SquareWaveProc with the supplied path processors.
func NewSquareWaveProc(lambda, amplitude float64) *SquareWaveProc {
	return &SquareWaveProc{lambda / 2, amplitude / lambda, false, false}
}

// Process implements the PathProcessor interface.
func (sp *SquareWaveProc) Process(p *Path) []*Path {
	// Chunk up path into pieces
	pp1 := NewMunchProc(sp.HalfLambda).Process(p)
	p1, _ := ConcatenatePaths(pp1...)

	n := len(p1.steps)
	last := p1.steps[0][0]
	path := NewPath(last)
	left := !sp.Flip
	for i := 1; i < n; i++ {
		// Each step is a half wave
		cur := p1.steps[i][0]
		dx, dy := cur[0]-last[0], cur[1]-last[1]
		ndx, ndy := dy*sp.Scale, -dx*sp.Scale
		if left {
			ndx, ndy = -ndx, -ndy
		}
		path.AddStep([]float64{last[0] + ndx, last[1] + ndy})
		path.AddStep([]float64{cur[0] + ndx, cur[1] + ndy})
		path.AddStep([]float64{cur[0], cur[1]})
		last = cur
		left = !left
	}

	if sp.KeepZero {
		return []*Path{path}
	}

	// Filter out internal zero-crossing points
	n = len(path.steps)
	fpath := NewPath(path.steps[0][0])
	for i := 1; i < n-1; i++ {
		if i%3 == 0 {
			continue
		}
		fpath.AddStep(path.steps[i][0])
	}
	fpath.AddStep(path.steps[n-1][0])

	return []*Path{fpath}
}
