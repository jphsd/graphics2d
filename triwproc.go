package graphics2d

// TriangleWaveProc applies a triangle wave along a path with a defined wave length and amplitude.
// The wave starts and ends on zero-crossing points and the last half wave is truncated to the
// path length remaining. The internal zero-crossing points can be optionally preserved.
type TriangleWaveProc struct {
	HalfLambda float64 // Half wave length
	Scale      float64 // Ratio of amplitude to lambda
	KeepZero   bool    // Keeps internal zero-point crossings if set
	Flip       bool    // Flips the wave phase by 180 (pi) if set
}

// NewTriangleWaveProc creates a new TriangleWaveProc with the supplied wave length and amplitutde.
func NewTriangleWaveProc(lambda, amplitude float64) TriangleWaveProc {
	return TriangleWaveProc{lambda / 2, amplitude / lambda, false, false}
}

// Process implements the PathProcessor interface.
func (tp TriangleWaveProc) Process(p *Path) []*Path {
	// Chunk up path into pieces
	pp1 := NewMunchProc(tp.HalfLambda).Process(p)
	p1, _ := ConcatenatePaths(pp1...)

	n := len(p1.steps)
	last := p1.steps[0][0]
	path := NewPath(last)
	left := !tp.Flip
	for i := 1; i < n; i++ {
		// Each step is a half wave
		cur := p1.steps[i][0]
		dx, dy := cur[0]-last[0], cur[1]-last[1]
		ndx, ndy := dy*tp.Scale, -dx*tp.Scale
		if left {
			ndx, ndy = -ndx, -ndy
		}
		mid := []float64{(last[0] + cur[0]) / 2, (last[1] + cur[1]) / 2}
		path.AddStep([]float64{mid[0] + ndx, mid[1] + ndy})
		path.AddStep([]float64{cur[0], cur[1]})
		last = cur
		left = !left
	}

	if tp.KeepZero {
		return []*Path{path}
	}

	// Filter out internal zero-crossing points
	n = len(path.steps)
	fpath := NewPath(path.steps[0][0])
	for i := 1; i < n-1; i++ {
		if i%2 == 0 {
			continue
		}
		fpath.AddStep(path.steps[i][0])
	}
	fpath.AddStep(path.steps[n-1][0])

	if p.Closed() {
		fpath.Close()
	}

	return []*Path{fpath}
}
