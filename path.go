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
	steps  [][]float64
	closed bool
	bounds image.Rectangle
	// Caching flattened path
	flattened *Path
	tolerance float64
	// Processed from
	parent *Path
}

func NewPath(start []float64) *Path {
	np := Path{make([][]float64, 1), false, image.Rectangle{}, nil, 1, nil}
	np.steps[0] = start
	return &np
}

func (p *Path) AddStep(step []float64) error {
	if p.closed {
		return fmt.Errorf("Path is closed, adding a step is forbidden")
	}
	if len(step)%2 != 0 {
		return fmt.Errorf("Step locations must be an even number of float64")
	}
	p.steps = append(p.steps, step)
	p.bounds = image.Rectangle{}
	p.flattened = nil
	return nil
}

func (p *Path) Steps() [][]float64 {
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
	if util.Equalsf64(d, p.tolerance) && p.flattened != nil {
		return p.flattened
	}
	p.tolerance = d
	d2 := d * d
	res := make([][]float64, 1)
	res[0] = p.steps[0]
	cp := res[0]
	// For all remaining steps in path
	for i := 1; i < len(p.steps); i++ {
		fp := flattenStep(0, d2, toPoints(cp, p.steps[i]))[1:]
		cp = fp[len(fp)-1]
		res = append(res, fp...)
	}

	p.flattened = &Path{res, p.closed, image.Rectangle{}, nil, 1, p}
	return p.flattened
}

func toPoints(cp []float64, pts []float64) [][]float64 {
	res := make([][]float64, len(pts)/2+1)
	res[0] = cp
	for i, j := 1, 0; i < len(res); i++ {
		res[i] = pts[j : j+2]
		j += 2
	}
	return res
}

func flattenStep(n int, d2 float64, pts [][]float64) [][]float64 {
	if cpWithinD2(d2, pts) {
		return [][]float64{pts[0], pts[len(pts)-1]}
	}
	lr := util.SplitCurve(pts, 0.5)
	res := append([][]float64{}, flattenStep(n+1, d2, lr[0])...)
	// rhs needs reversing
	rhs := flattenStep(n+1, d2, lr[1])[1:] // Trim off first point
	res = append(res, rhs...)
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
		for _, s := range p.steps {
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
		// Bump Max by 1 as image.Rectangle is exclusive on the high end
		rect.Max.X += 1
		rect.Max.Y += 1
		p.bounds = rect
	}
	return p.bounds
}

// Deepish copy - points themselves aren't duplicated
func (p *Path) Copy() *Path {
	steps := make([][]float64, len(p.steps))
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
		for i := 0; i < len(step); i += 2 {
			str += fmt.Sprintf("%f %f ", step[i], step[i+1])
		}
	}
	if p.closed {
		str += "C"
	}
	return str
}

// Apply an affine transfrom to the points in a path to
// create a new one.
func (p *Path) Transform(xfm *Aff3) *Path {
	// x' = xfm[3*0+0]*x + xfm[3*0+1]*y + xfm[3*0+2]
	// y' = xfm[3*1+0]*x + xfm[3*1+1]*y + xfm[3*1+2]
	steps := make([][]float64, len(p.steps))
	for i, step := range p.steps {
		nstep := make([]float64, len(step))
		for j := 0; j < len(step); j += 2 {
			x, y := step[j+0], step[j+1]
			nstep[j+0] = xfm[0]*x + xfm[1]*y + xfm[2]
			nstep[j+1] = xfm[3]*x + xfm[4]*y + xfm[5]
		}
		steps[i] = nstep
	}
	return &Path{steps, p.closed, image.Rectangle{}, nil, 1, nil}
}
