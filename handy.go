package graphics2d

import (
	"math"
	"math/rand"
)

// HandyProc applies a modified form of line rendering as outlined in Wood12.
// Note the lines are not smoothed and closed paths ae not preserved.
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
		op2 := steps[i][len(steps[i])-1]
		opa, opb := Lerp(0.5, op1, op2), Lerp(0.75, op1, op2)
		for _, path := range paths {
			path.AddStep(jitter(opa, hp.R))
			path.AddStep(jitter(opb, hp.R))
			path.AddStep(jitter(op2, hp.R))
		}
		op1 = op2
	}

	return paths
}

func jitter(pt []float64, r float64) []float64 {
	th := rand.Float64() * TwoPi
	dx, dy := math.Cos(th)*r, math.Sin(th)*r
	return []float64{pt[0] + dx, pt[1] + dy}
}
