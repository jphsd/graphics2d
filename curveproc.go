package graphics2d

import "github.com/jphsd/graphics2d/util"

// CurveStyle determines how the curve behaves relative to the path points. With Bezier, the
// path will intersect the mid-point of each path step. With Catmul, the path will intersect
// point.
type CurveStyle int

// Constants for curve styles.
const (
	Bezier CurveStyle = iota
	Quad
	CatmullRom
)

// CurveProc replaces the steps on a path with cubics. The locations of the control points
// are controlled by the Style setting and whether or not the path is closed.
type CurveProc struct {
	Scale float64
	Style CurveStyle
}

// Process implements the PathProcessor interface.
func (cp *CurveProc) Process(p *Path) []*Path {
	steps := p.Steps()
	ns := len(steps)
	if ns < 2 {
		return []*Path{p}
	}

	// Truncate steps to end points
	points := make([][]float64, ns)
	for i, step := range steps {
		points[i] = step[len(step)-1]
	}

	if p.closed && util.EqualsP(points[0], points[ns-1]) {
		ns--
	}

	res := []*Path{}

	// Bezier - curve passes through mid points. If path is open, then passes through start and end too.

	if cp.Style == Bezier {
		// Calc mid points
		mp := make([][]float64, ns)
		for i := 0; i < ns-1; i++ {
			mp[i] = util.Centroid(points[i], points[i+1])
		}
		mp[ns-1] = util.Centroid(points[ns-1], points[0])

		// Create path
		if p.closed {
			res = append(res, NewPath(mp[0]))
		} else {
			res = append(res, NewPath(points[0]))
			res[0].AddStep(mp[0])
		}
		for i := 1; i < ns-1; i++ {
			c1 := Lerp(cp.Scale, mp[i-1], points[i])
			c2 := Lerp(cp.Scale, mp[i], points[i])
			res[0].AddStep(c1, c2, mp[i])
		}
		if p.closed {
			c1 := Lerp(cp.Scale, mp[ns-2], points[ns-1])
			c2 := Lerp(cp.Scale, mp[ns-1], points[ns-1])
			res[0].AddStep(c1, c2, mp[ns-1])
			c1 = Lerp(cp.Scale, mp[ns-1], points[0])
			c2 = Lerp(cp.Scale, mp[0], points[0])
			res[0].AddStep(c1, c2, mp[0])
			res[0].Close()
		} else {
			res[0].AddStep(points[ns-1])
		}

		return res
	}

	// Quad - curve passes through mid points. If path is open, then passes through start and end too.

	if cp.Style == Quad {
		// Calc mid points
		mp := make([][]float64, ns)
		for i := 0; i < ns-1; i++ {
			mp[i] = util.Centroid(points[i], points[i+1])
		}
		mp[ns-1] = util.Centroid(points[ns-1], points[0])

		// Create path
		if p.closed {
			res = append(res, NewPath(mp[0]))
		} else {
			res = append(res, NewPath(points[0]))
			res[0].AddStep(mp[0])
		}
		for i := 1; i < ns-1; i++ {
			res[0].AddStep(points[i], mp[i])
		}
		if p.closed {
			res[0].AddStep(points[ns-1], mp[ns-1])
			res[0].AddStep(points[0], mp[0])
			res[0].Close()
		} else {
			res[0].AddStep(points[ns-1])
		}

		return res
	}

	// Catmull-Rom - curve passes through all points

	// Calc opposite tangents
	ops := make([][]float64, ns)
	for i := 1; i < ns-1; i++ {
		ops[i] = []float64{(points[i+1][0] - points[i-1][0]) / 2, (points[i+1][1] - points[i-1][1]) / 2}
	}
	if p.closed {
		ops[0] = []float64{(points[1][0] - points[ns-1][0]) / 2, (points[1][1] - points[ns-1][1]) / 2}
		ops[ns-1] = []float64{(points[0][0] - points[ns-2][0]) / 2, (points[0][1] - points[ns-2][1]) / 2}
	} else {
		ops[0] = []float64{0, 0}
		ops[ns-1] = ops[0]
	}

	// Create path
	res = append(res, NewPath(points[0]))
	if p.closed {
		for i := 0; i < ns-1; i++ {
			c1, c2 := cp.calcControlOpp(points[i], ops[i], points[i+1], ops[i+1])
			res[0].AddStep(c1, c2, points[i+1])
		}
		c1, c2 := cp.calcControlOpp(points[ns-1], ops[ns-1], points[0], ops[0])
		res[0].AddStep(c1, c2, points[0])
		res[0].Close()
	} else {
		// Insert quads for start and end
		_, c2 := cp.calcControlOpp(points[0], ops[0], points[1], ops[1])
		res[0].AddStep(c2, points[1])
		for i := 1; i < ns-2; i++ {
			c1, c2 := cp.calcControlOpp(points[i], ops[i], points[i+1], ops[i+1])
			res[0].AddStep(c1, c2, points[i+1])
		}
		c1, _ := cp.calcControlOpp(points[ns-2], ops[ns-2], points[ns-1], ops[ns-1])
		res[0].AddStep(c1, points[ns-1])
	}

	return res
}

// Lerp performs a linear interpolation between two points.
func Lerp(t float64, p1, p2 []float64) []float64 {
	return []float64{util.Lerp(t, p1[0], p2[0]), util.Lerp(t, p1[1], p2[1])}
}

func (cp *CurveProc) calcControlOpp(p1, op1, p2, op2 []float64) ([]float64, []float64) {
	dx1, dy1 := op1[0]*cp.Scale, op1[1]*cp.Scale
	dx2, dy2 := -op2[0]*cp.Scale, -op2[1]*cp.Scale
	return []float64{p1[0] + dx1, p1[1] + dy1}, []float64{p2[0] + dx2, p2[1] + dy2}
}
