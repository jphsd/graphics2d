package graphics2d

// HandyProc applies a modified form of line rendering as outlined in Wood12.
// Note the lines are not smoothed and closed paths are not preserved.
type HandyProc struct {
	N int     // Repetitions
	R float64 // Jitter radius
}

// Process implements the PathProcessor interface.
func (hp *HandyProc) Process(p *Path) []*Path {
	steps := p.Steps()
	ns := len(steps)
	op1 := steps[0][0]

	paths := make([]*Path, hp.N)
	for i, _ := range paths {
		paths[i] = NewPath(jitter(op1, hp.R))
	}

	for i := 1; i < ns; i++ {
		nc := len(steps[i]) - 1
		op2 := steps[i][nc]
		if nc == 0 {
			// This is a linear step, add extra points
			opa, opb := Lerp(0.5, op1, op2), Lerp(0.75, op1, op2)
			for _, path := range paths {
				path.AddStep(jitter(opa, hp.R))
				path.AddStep(jitter(opb, hp.R))
				path.AddStep(jitter(op2, hp.R))
			}
		} else {
			// Just jitter control and end points
			for _, path := range paths {
				sps := make([][]float64, nc+1)
				for j := range nc + 1 {
					sps[j] = jitter(steps[i][j], hp.R)
				}
				path.AddStep(sps...)
			}
		}
		op1 = op2
	}

	return paths
}
