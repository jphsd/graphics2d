package graphics2d

import (
	"fmt"
	"image"
	"math"

	"github.com/jphsd/graphics2d/util"
)

/*
 * Path is a simple path - continous without interruption, and may be closed.
 * It's created with an initial starting point via NewPath() and then steps are
 * added to it with AddStep(). Steps can be single points or sets of control points.
 * For a given step size, the following curve is generated -
 *  2 - Line
 *  4 - Quad
 *  6 - Cubic
 *  ... - Higher order curves
 * DeCasteljau's algorithm is used to calculate the curve
 * Closing a path means steps can no longer be added to it and the last step is
 * connected to the first as a line.
 */

type Path struct {
	// step, point, ordinal
	steps  [][][]float64
	closed bool
	bounds image.Rectangle
	// Caching flattened path
	flattened *Path
	tolerance float64
	// Processed from
	parent *Path
}

func NewPath(start []float64) *Path {
	np := Path{make([][][]float64, 1), false, image.Rectangle{}, nil, 1, nil}
	np.steps[0] = [][]float64{start}
	return &np
}

// AddStep takes an array of points and treats n-1 of them as control points and the
// last as a point on the curve.
func (p *Path) AddStep(points [][]float64) error {
	if p.closed {
		return fmt.Errorf("path is closed, adding a step is forbidden")
	}
	p.steps = append(p.steps, points)
	p.bounds = image.Rectangle{}
	p.flattened = nil
	return nil
}

func (p *Path) Steps() [][][]float64 {
	return p.steps[:]
}

func (p *Path) Close() {
	p.closed = true
}

func (p *Path) Closed() bool {
	return p.closed
}

// Recursively subdivide path until the control points are within d of
// the line through the end points. Emit new path.
func (p *Path) Flatten(d float64) *Path {
	if p.flattened != nil && util.Equalsf64(d, p.tolerance) {
		return p.flattened
	}
	p.tolerance = d
	d2 := d * d
	res := make([][][]float64, 1)
	res[0] = p.steps[0]
	cp := res[0][0]
	// For all remaining steps in path
	for i := 1; i < len(p.steps); i++ {
		fp := flattenStep(0, d2, toPoints(cp, p.steps[i]))
		for _, ns := range fp {
			res = append(res, ns[1:])
			cp = ns[len(ns)-1]
		}
	}

	p.flattened = &Path{res, p.closed, image.Rectangle{}, nil, 1, p}
	return p.flattened
}

func toPoints(cp []float64, pts [][]float64) [][]float64 {
	res := make([][]float64, len(pts)+1)
	res[0] = cp
	copy(res[1:], pts)
	return res
}

func flattenStep(n int, d2 float64, pts [][]float64) [][][]float64 {
	if cpWithinD2(d2, pts) {
		return [][][]float64{{pts[0], pts[len(pts)-1]}}
	}
	lr := util.SplitCurve(pts, 0.5)
	res := append([][][]float64{}, flattenStep(n+1, d2, lr[0])...)
	res = append(res, flattenStep(n+1, d2, lr[1])...)
	return res
}

func cpWithinD2(d2 float64, pts [][]float64) bool {
	// First and last are end points
	l := len(pts)
	if l == 2 {
		// Trivial case
		return true
	}
	start, cpts, end := pts[0], pts[1:l-1], pts[l-1]
	for _, cp := range cpts {
		pd2 := util.DistanceToLineSquared(start, end, cp)
		if pd2 > d2 {
			return false
		}
	}

	return true
}

// Path guaranteed to fit within bounds
func (p *Path) Bounds() image.Rectangle {
	if p.bounds.Empty() {
		rect := image.Rectangle{}
		for _, pts := range p.steps {
			for _, s := range pts {
				fx, fy := int(math.Floor(s[0])), int(math.Floor(s[1]))
				cx, cy := int(math.Ceil(s[0])), int(math.Ceil(s[1]))
				if rect.Empty() {
					rect.Min.X = fx
					rect.Min.Y = fy
					rect.Max.X = cx
					rect.Max.Y = cy
				} else {
					if rect.Min.X > fx {
						rect.Min.X = fx
					}
					if rect.Min.Y > fy {
						rect.Min.Y = fy
					}
					if rect.Max.X < cx {
						rect.Max.X = cx
					}
					if rect.Max.Y < cy {
						rect.Max.Y = cy
					}
				}
			}
		}
		// Bump Max by 1 as image.Rectangle is exclusive on the high end
		rect.Max.X += 1
		rect.Max.Y += 1
		p.bounds = rect
	}
	return p.bounds
}

// Deepish copy - points themselves aren't duplicated
func (p *Path) Copy() *Path {
	steps := make([][][]float64, len(p.steps))
	copy(steps, p.steps)
	return &Path{steps, p.closed, p.bounds, nil, 1, p.parent}
}

// Return path's parent
func (p *Path) Parent() *Path {
	return p.parent
}

// Apply a processor to a path
func (p *Path) Process(proc PathProcessor) []*Path {
	paths := proc.Process(p)
	// Fix parent
	for _, path := range paths {
		path.parent = p
	}

	return paths
}

// Applies a collection of PathProcessors to a path
func (p *Path) CompoundProcess(procs []PathProcessor) []*Path {
	paths := []*Path{p}
	if len(procs) == 0 {
		return paths
	}

	for _, proc := range procs {
		npaths := []*Path{}
		for _, path := range paths {
			npaths = append(npaths, proc.Process(path)...)
		}
		paths = npaths
	}

	return paths
}

func (p *Path) String() string {
	str := ""
	for _, step := range p.steps {
		str += "S "
		for _, pts := range step {
			str += fmt.Sprintf("%f %f ", pts[0], pts[1])
		}
	}
	if p.closed {
		str += "C"
	}
	return str
}

// Apply an affine transfrom to the points in a path to create a new one.
func (p *Path) Transform(xfm *Aff3) *Path {
	steps := make([][][]float64, len(p.steps))
	for i, step := range p.steps {
		steps[i] = xfm.Apply(step...)
	}
	return &Path{steps, p.closed, image.Rectangle{}, nil, 1, nil}
}

// Simplify breaks up a path into steps where for any step, its control points are all on the
// same side and its midpoint is well behaved. If a step doesn't meet the criteria, it is
// recursively subdivide in half until it does.
func (p *Path) Simplify() *Path {
	res := make([][][]float64, 1)
	res[0] = p.steps[0]
	cp := res[0][0]
	// For all remaining steps in path
	for i := 1; i < len(p.steps); i++ {
		fp := simplifyStep(toPoints(cp, p.steps[i]))
		for _, ns := range fp {
			res = append(res, ns[1:])
			cp = ns[len(ns)-1]
		}
	}

	return &Path{res, p.closed, image.Rectangle{}, nil, 1, p}
}

func simplifyStep(points [][]float64) [][][]float64 {
	if cpSafe(points) {
		return [][][]float64{points}
	}
	lr := util.SplitCurve(points, 0.5)
	res := append([][][]float64{}, simplifyStep(lr[0])...)
	res = append(res, simplifyStep(lr[1])...)
	return res
}

// cpSafe returns true if all the control points are on the same side of
// the line formed by start and the last point in step, and the t=0.5 point
// is close to the geometric center of the polygon defined by the points.
func cpSafe(points [][]float64) bool {
	if len(points) < 3 {
		// Either a point or line
		return true
	}

	n := len(points)
	start := points[0]
	end := points[n-1]
	side := dotprod(start, end, points[1]) < 0
	for i := 2; i < n-1; i++ {
		if (dotprod(start, end, points[i]) < 0) != side {
			return false
		}
	}

	c := util.Centroid(points...)
	v := util.DeCasteljau(points, 0.5)
	bb := util.BoundingBox(points...)
	dx := bb[1][0] - bb[0][0]
	dy := bb[1][1] - bb[0][1]
	dx, dy = dx/10, dy/10
	// Crude check v is within 10% of c based on bb size
	return v[0] < c[0]+dx && v[0] > c[0]-dx && v[1] < c[1]+dy && v[1] > c[1]-dy
}

func dotprod(p1, p2, p3 []float64) float64 {
	return (p3[0]-p1[0])*(p2[1]-p1[1]) - (p3[1]-p1[1])*(p2[0]-p1[0])
}
