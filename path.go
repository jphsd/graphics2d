package graphics2d

import (
	"fmt"
	"image"
	"math"

	. "github.com/jphsd/graphics2d/util"
)

/*
 * Path is a simple path - continuous without interruption, and may be closed.
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

// Path contains the housekeeping necessary for path building.
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

// NewPath creates a new path starting at start.
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

// AddSteps adds multiple steps to the path.
func (p *Path) AddSteps(steps [][][]float64) error {
	for _, step := range steps {
		err := p.AddStep(step)
		if err != nil {
			return err
		}
	}
	return nil
}

// Concatenate adds the path to this path. If either path is closed then an error
// is returned. If the paths aren't coincident, then they are joind with a line.
func (p *Path) Concatenate(path *Path) error {
	if p.closed {
		return fmt.Errorf("path is closed, adding a step is forbidden")
	}
	if path.closed {
		return fmt.Errorf("can't add a closed path")
	}
	lstep := p.steps[len(p.steps)-1]
	last := lstep[len(lstep)-1]

	steps := path.Steps()
	if EqualsP(last, steps[0][0]) {
		// End of p is coincident with sep[0][0] of path
		p.AddSteps(steps[1:])
	} else {
		// Line to steps[0][0]
		p.AddSteps(steps)
	}
	return nil
}

// Steps returns a shallow copy of all the steps in the path.
func (p *Path) Steps() [][][]float64 {
	return p.steps[:]
}

// Parts returns the steps of a path, each prepended with its start.
func (p *Path) Parts() [][][]float64 {
	n := len(p.steps)
	if n == 1 {
		// This is a point
		return [][][]float64{{p.steps[0][0]}}
	}
	fpts := make([][][]float64, n-1, n)
	cp := p.steps[0][0]
	for i := 1; i < n; i++ {
		pts := toPoints(cp, p.steps[i])
		fpts[i-1] = pts
		cp = pts[len(pts)-1]
	}
	if p.closed && !EqualsP(cp, p.steps[0][0]) {
		fpts = append(fpts, [][]float64{cp, p.steps[0][0]})
	}
	return fpts
}

// Close marks the path as closed.
func (p *Path) Close() {
	p.closed = true
}

// Closed returns true if the path is closed.
func (p *Path) Closed() bool {
	return p.closed
}

// PartsToPath constructs a new path by concatenating the parts.
func PartsToPath(pts [][][]float64) (*Path, error) {
	res := NewPath(pts[0][0])
	if len(pts[0]) == 1 {
		return res, nil
	}
	for i, part := range pts {
		if EqualsP(part[0], part[len(part)-1]) {
			return nil, fmt.Errorf("part %d start and end are coincident", i)
		}
		res.AddStep(part[1:])
	}
	return res, nil
}

// Flatten works by recursively subdividing the path until the control points are within d of
// the line through the end points.
func (p *Path) Flatten(d float64) *Path {
	if p.flattened != nil && d >= p.tolerance {
		// Path has already been flattened at least to the degree we're looking for
		return p.flattened
	}
	p.tolerance = d
	res := make([][][]float64, 1)
	res[0] = p.steps[0]
	cp := res[0][0]
	// For all remaining steps in path
	for i := 1; i < len(p.steps); i++ {
		fp := flattenPart(d*d, toPoints(cp, p.steps[i]))
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

func flattenPart(d2 float64, pts [][]float64) [][][]float64 {
	if cpWithinD2(d2, pts) {
		return [][][]float64{{pts[0], pts[len(pts)-1]}}
	}
	lr := SplitCurve(pts, 0.5)
	res := append([][][]float64{}, flattenPart(d2, lr[0])...)
	res = append(res, flattenPart(d2, lr[1])...)
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
		pd2 := DistanceToLineSquared(start, end, cp)
		if pd2 > d2 {
			return false
		}
	}

	return true
}

// FlattenPart works by subdividing the curve until its control points are within d2 (d squared)
// of the line through the end points.
func FlattenPart(d float64, pts [][]float64) [][][]float64 {
	return flattenPart(d*d, pts)
}

// Bounds calculates a rectangle that the Path is guaranteed to fit within. It's unlikely to
// be the minimal bounding rectangle for the path since the control points are also included.
// If a tight bounding rectangle is required then use CalcExtremities().
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

// Copy performs a Deepish copy - points themselves aren't duplicated.
func (p *Path) Copy() *Path {
	steps := make([][][]float64, len(p.steps))
	copy(steps, p.steps)
	return &Path{steps, p.closed, p.bounds, nil, 1, p.parent}
}

// Open performs a Deepish copy like Copy() but leaves the path open.
func (p *Path) Open() *Path {
	steps := make([][][]float64, len(p.steps))
	copy(steps, p.steps)
	return &Path{steps, false, p.bounds, nil, 1, p.parent}
}

// Parent returns the path's parent
func (p *Path) Parent() *Path {
	return p.parent
}

// Process applies a processor to a path.
func (p *Path) Process(proc PathProcessor) []*Path {
	paths := proc.Process(p)
	// Fix parent
	for _, path := range paths {
		path.parent = p
	}

	return paths
}

// String converts a path into a string.
func (p *Path) String() string {
	step := p.steps[0]
	str := fmt.Sprintf("P %f,%f ", step[0][0], step[0][1])
	for i := 1; i < len(p.steps); i++ {
		step = p.steps[i]
		str += "S "
		for _, pts := range step {
			str += fmt.Sprintf("%f,%f ", pts[0], pts[1])
		}
	}
	if p.closed {
		str += "C"
	}
	return str
}

// Transform applies an affine transform to the points in a path to create a new one.
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
	for i := 1; i < len(p.steps); i++ {
		points := toPoints(cp, p.steps[i])
		if len(points) > 2 {
			// Split based on extremities
			parts := simplifyExtremities(points)
			fp := [][][]float64{}
			for _, part := range parts {
				// check parts are simple
				fp = append(fp, simplifyStep(part)...)
			}
			// Turn points back into path step
			for _, ns := range fp {
				res = append(res, ns[1:])
				cp = ns[len(ns)-1]
			}
		} else {
			res = append(res, points[1:])
			cp = points[len(points)-1]
		}
	}

	return &Path{res, p.closed, image.Rectangle{}, nil, 1, p}
}

// Chop curve into pieces based on maxima, minima and inflections in x and y.
func simplifyExtremities(points [][]float64) [][][]float64 {
	tvals := CalcExtremities(points)
	nt := len(tvals)
	if nt < 3 {
		return [][][]float64{points}
	}
	rhs := points
	res := make([][][]float64, nt-1)
	for i := 1; i < nt-1; i++ {
		lr := SplitCurve(rhs, tvals[i])
		res[i-1] = lr[0]
		rhs = lr[1]
	}
	res[nt-2] = rhs
	return res
}

// simplifyStep recursively cuts the curve in half until cpSafe is
// satisfied.
func simplifyStep(points [][]float64) [][][]float64 {
	if cpSafe(points) {
		return [][][]float64{points}
	}
	lr := SplitCurve(points, 0.5)
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
	side := CrossProduct(start, end, points[1]) < 0
	for i := 2; i < n-1; i++ {
		if (CrossProduct(start, end, points[i]) < 0) != side {
			return false
		}
	}

	c := Centroid(points...)
	v := DeCasteljau(points, 0.5)
	bb := BoundingBox(points...)
	dx := bb[1][0] - bb[0][0]
	dy := bb[1][1] - bb[0][1]
	dx, dy = dx/40, dy/40
	// Crude check v is within 5% of c based on bb size
	return v[0] < c[0]+dx && v[0] > c[0]-dx && v[1] < c[1]+dy && v[1] > c[1]-dy
}

// ReversePath returns a new path describing the current path in reverse order (i.e start and end switched).
func (p *Path) ReversePath() *Path {
	path, _ := PartsToPath(ReverseParts(p.Parts()))
	return path
}

// ReverseParts reverses the order (and points) of the supplied part slice.
func ReverseParts(pts [][][]float64) [][][]float64 {
	n := len(pts)
	res := make([][][]float64, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = reversePoints(pts[j])
		j--
	}
	return res
}

// [pts][x/y]
func reversePoints(cp [][]float64) [][]float64 {
	n := len(cp)
	res := make([][]float64, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = cp[j]
		j--
	}
	return res
}
