package graphics2d

import (
	"github.com/jphsd/graphics2d/util"
	"math/rand"
)

// MPDProc contains the variables that control the degree to which a step is chopped up into smaller
// line segments. Unlike JitterProc, the step end points don't vary. Can be used with MunchProc to get
// a hand drawn look.
type MPDProc struct {
	Perc  float64 // Percentage of step length used as initial displacement
	Itrs  int     // Number of iterations to perform
	Scale float64 // Multiplier used on displacement per iteration
}

// Process implements the PathProcessor interface.
func (m *MPDProc) Process(p *Path) []*Path {
	parts := p.Parts()
	np := len(parts)
	if np == 0 {
		return []*Path{p}
	}

	nparts := [][][]float64{}
	for _, part := range parts {
		a := part[0]
		b := part[len(part)-1]
		points := m.MPD(a, b)
		lp := len(points)
		cp := points[0]
		for i := 1; i < lp; i++ {
			np := points[i]
			nparts = append(nparts, [][]float64{cp, np})
			cp = np
		}
	}

	return []*Path{PartsToPath(nparts...)}
}

// MPD takes two points and adds points between them using the mid-point displacement algorithm
// drive by the parameters stored in the MPDProc structure.
func (m *MPDProc) MPD(a, b []float64) [][]float64 {
	if m.Itrs == 0 {
		return [][]float64{a, b}
	}
	v := util.Vec(a, b)
	d := util.VecMag(v)
	n := []float64{-v[1] / d, v[0] / d}
	return m.mpdhelper(a, b, n, m.Itrs, d*m.Perc)
}

func (m *MPDProc) mpdhelper(a, b, n []float64, itr int, disp float64) [][]float64 {
	mpx, mpy := (a[0]+b[0])/2, (a[1]+b[1])/2
	d := (rand.Float64()*2 - 1) * disp
	c := []float64{mpx + d*n[0], mpy + d*n[1]}
	if itr == 1 {
		return [][]float64{a, c, b}
	}
	ndisp := disp * m.Scale
	lhs := m.mpdhelper(a, c, n, itr-1, ndisp)
	rhs := m.mpdhelper(c, b, n, itr-1, ndisp)
	return append(lhs, rhs[1:]...)
}
