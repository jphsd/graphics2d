package graphics2d

import (
	"encoding/json"
	"fmt"
	"image"
	"math"

	"github.com/jphsd/graphics2d/util"
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
	// Caching flattened, simplified and reversed paths, and tangents
	flattened  *Path
	tolerance  float64
	simplified *Path
	reversed   *Path
	tangents   [][][]float64
	// Processed from
	parent *Path
}

// NewPath creates a new path starting at start.
func NewPath(start []float64) *Path {
	np := &Path{}
	np.steps = make([][][]float64, 1)
	if InvalidPoint(start) {
		panic("invalid point (NaN) for path start")
	}
	np.steps[0] = [][]float64{start}
	return np
}

// AddStep takes an array of points and treats n-1 of them as control points and the
// last as a point on the curve.
func (p *Path) AddStep(points ...[]float64) error {
	n := len(points)
	if n == 0 {
		return nil
	}
	if p.closed {
		return fmt.Errorf("path is closed, adding a step is forbidden")
	}

	lastStep := p.steps[len(p.steps)-1]
	last := lastStep[len(lastStep)-1]
	npoints := make([][]float64, 0, n)
	for i, pt := range points {
		if InvalidPoint(pt) {
			panic(fmt.Sprintf("invalid step point (NaN) at %d", i))
		}
		if util.EqualsP(last, pt) {
			// Ignore coincident points
			continue
		}
		npoints = append(npoints, pt)
		last = pt
	}
	if len(npoints) == 0 {
		return nil
	}

	p.steps = append(p.steps, npoints)
	p.bounds = image.Rectangle{}
	p.flattened = nil
	p.simplified = nil
	p.tangents = nil
	p.reversed = nil
	return nil
}

// AddSteps adds multiple steps to the path.
func (p *Path) AddSteps(steps ...[][]float64) error {
	for _, step := range steps {
		err := p.AddStep(step...)
		if err != nil {
			return err
		}
	}
	return nil
}

// Concatenate adds the paths to this path. If any path is closed then an error
// is returned. If the paths aren't coincident, then they are joined with a line.
func (p *Path) Concatenate(paths ...*Path) error {
	if p.closed {
		return fmt.Errorf("path is closed, adding a step is forbidden")
	}
	for _, path := range paths {
		if path.closed {
			return fmt.Errorf("can't add a closed path")
		}
	}

	lstep := p.steps[len(p.steps)-1]
	last := lstep[len(lstep)-1]
	for _, path := range paths {
		steps := path.Steps()
		if util.EqualsP(last, steps[0][0]) {
			// End of p is coincident with sep[0][0] of path
			p.AddSteps(steps[1:]...)
		} else {
			// Line to steps[0][0]
			p.AddSteps(steps...)
		}
		lstep = steps[len(steps)-1]
		last = lstep[len(lstep)-1]
	}
	return nil
}

// ConcatenatePaths concatenates all the paths into a new path. If any path is closed then an error
// is returned. If the paths aren't coincident, then they are joined with a line.
func ConcatenatePaths(paths ...*Path) (*Path, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	path := paths[0].Copy()
	err := path.Concatenate(paths[1:]...)
	if err != nil {
		return nil, err
	}
	return path, nil
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
	if p.closed && !util.EqualsP(cp, p.steps[0][0]) {
		fpts = append(fpts, [][]float64{cp, p.steps[0][0]})
	}
	return fpts
}

// Close marks the path as closed.
func (p *Path) Close() {
	p.AddStep(p.steps[0][0])
	p.closed = true
}

// Closed returns true if the path is closed.
func (p *Path) Closed() bool {
	return p.closed
}

// PartsToPath constructs a new path by concatenating the parts.
func PartsToPath(parts ...[][]float64) *Path {
	if len(parts) == 0 {
		return nil
	}

	res := NewPath(parts[0][0])
	if len(parts[0]) == 1 {
		return res
	}

	for i, part := range parts {
		for j, pt := range part {
			if InvalidPoint(pt) {
				panic(fmt.Sprintf("invalid point (NaN) in part %d,%d", i, j))
			}
		}
		res.AddStep(part[1:]...)
	}
	return res
}

// InvalidPoint checks that both values are valid (i.e. not NaN)
func InvalidPoint(p []float64) bool {
	return p[0] != p[0] || p[1] != p[1]
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

	path := &Path{}
	path.steps = res
	path.closed = p.closed
	path.parent = p
	p.flattened = path
	return path
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
	lr := util.SplitCurve(pts, 0.5)
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
		pd2, _, _ := util.DistanceToLineSquared(start, end, cp)
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
					rect.Max.X = cx + 1
					rect.Max.Y = cy + 1
				} else {
					if rect.Min.X > fx {
						rect.Min.X = fx
					}
					if rect.Min.Y > fy {
						rect.Min.Y = fy
					}
					if rect.Max.X <= cx {
						rect.Max.X = cx + 1
					}
					if rect.Max.Y <= cy {
						rect.Max.Y = cy + 1
					}
				}
			}
		}
		p.bounds = rect
	}
	return p.bounds
}

// Copy performs a deepish copy - points themselves aren't duplicated.
func (p *Path) Copy() *Path {
	steps := make([][][]float64, len(p.steps))
	copy(steps, p.steps)

	path := &Path{}
	path.steps = steps
	path.closed = p.closed
	path.bounds = p.bounds
	path.parent = p.parent
	return path
}

// Open performs a deepish copy like Copy() but leaves the path open.
func (p *Path) Open() *Path {
	path := p.Copy()
	path.closed = false
	return path
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

	path := &Path{}
	path.steps = steps
	path.closed = p.closed
	path.parent = p
	return path
}

// Simplify breaks up a path into steps where for any step, its control points are all on the
// same side and its midpoint is well behaved. If a step doesn't meet the criteria, it is
// recursively subdivided in half until it does.
func (p *Path) Simplify() *Path {
	if p.simplified != nil {
		return p.simplified
	}
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

	path := &Path{}
	path.steps = res
	path.closed = p.closed
	path.parent = p
	p.simplified = path

	return path
}

// Chop curve into pieces based on maxima, minima and inflections in x and y.
func simplifyExtremities(points [][]float64) [][][]float64 {
	tvals := util.CalcExtremities(points)
	nt := len(tvals)
	if nt < 3 {
		return [][][]float64{points}
	}
	// Convert tvals to relative tvals
	rtvals := make([]float64, nt)
	lt := 0.0
	ll := 1.0
	for i := 0; i < nt; i++ {
		if i == 0 {
			// Reset
			lt = tvals[i]
			ll = 1 - lt
			rtvals[i] = lt
			continue
		}
		t := tvals[i]
		d := t - lt
		rtvals[i] = d / ll
		lt = t
		ll -= d
	}

	rhs := points
	res := make([][][]float64, nt-1)
	for i := 1; i < nt-1; i++ {
		lr := util.SplitCurve(rhs, rtvals[i])
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
	lr := util.SplitCurve(points, 0.5)
	res := append([][][]float64{}, simplifyStep(lr[0])...)
	res = append(res, simplifyStep(lr[1])...)
	return res
}

// cpSafe returns true if all the control points are on the same side of
// the line formed by start and the last point in step.
func cpSafe(points [][]float64) bool {
	if len(points) < 3 {
		// Either a point or line
		return true
	}

	n := len(points)
	start := points[0]
	end := points[n-1]
	side := util.CrossProduct(start, end, points[1]) < 0
	for i := 2; i < n-1; i++ {
		if (util.CrossProduct(start, end, points[i]) < 0) != side {
			return false
		}
	}
	return true
}

// Reverse returns a new path describing the current path in reverse order (i.e start and end switched).
func (p *Path) Reverse() *Path {
	if p.reversed != nil {
		return p.reversed
	}

	path := PartsToPath(ReverseParts(p.Parts())...)

	// If other aspects have already been calculated - reverse them too
	if p.flattened != nil {
		path.flattened = p.flattened.Reverse()
		path.tolerance = p.tolerance
	}
	if p.simplified != nil {
		path.simplified = p.simplified.Reverse()
	}
	if p.tangents != nil {
		path.tangents = reverseTangents(p.tangents)
	}
	if p.closed {
		path.closed = true
	}
	path.parent = p.parent
	p.reversed = path

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

// [part][start/end][normalized x/y]
func reverseTangents(tangents [][][]float64) [][][]float64 {
	n := len(tangents)
	res := make([][][]float64, n)

	for i, j := 0, n-1; i < n; i++ {
		cur := tangents[j]
		res[i] = [][]float64{{-cur[1][0], -cur[1][1]}, {-cur[0][0], -cur[0][1]}}
		j--
	}
	return res
}

// Line reduces a path to a line between its endpoints. For a closed path or one where the
// start and endpoints are coincident, a single point is returned.
func (p *Path) Line() *Path {
	n := len(p.steps)

	if n == 1 || p.closed {
		return p.Copy()
	}

	first := p.steps[0][0]
	lastStep := p.steps[n-1]
	last := lastStep[len(lastStep)-1]

	path := NewPath(first)
	if util.EqualsP(first, last) {
		return path
	}
	path.AddStep(last)
	return path
}

// Tangents returns the normalized start and end tangents of every part in the path.
// [part][start/end][normalized x/y]
func (p *Path) Tangents() [][][]float64 {
	if p.tangents != nil {
		return p.tangents
	}

	parts := p.Parts()
	n := len(parts)
	res := make([][][]float64, n)
	for i, part := range parts {
		tmp0, tmp1 := util.DeCasteljau(part, 0), util.DeCasteljau(part, 1)
		dx0, dy0 := unit(tmp0[2], tmp0[3])
		dx1, dy1 := unit(tmp1[2], tmp1[3])
		res[i] = [][]float64{{dx0, dy0}, {dx1, dy1}}
	}
	p.tangents = res
	return res
}

// PartsIntersection returns the location of where the two parts intersect or nil. Assumes the parts
// are the result of simplification. Uses a brute force approach for curves with d as the flattening
// value.
func PartsIntersection(part1, part2 [][]float64, d float64) []float64 {
	// Test bounding boxes first
	bb1, bb2 := util.BoundingBox(part1...), util.BoundingBox(part2...)
	if !util.BBOverlap(bb1, bb2) {
		return nil
	}

	// Flatten and calculate bounding boxes for part2
	fparts1, fparts2 := FlattenPart(d, part1), FlattenPart(d, part2)
	bbs2 := make([][][]float64, len(fparts2))
	for i, part := range fparts2 {
		bbs2[i] = util.BoundingBox(part...)
	}

	// Test each line in part1 against lines in part2 until we find an intersection
	for _, part := range fparts1 {
		bb1 = util.BoundingBox(part...)
		s1, e1 := part[0], part[len(part)-1]
		for j, bbp2 := range bbs2 {
			if !util.BBOverlap(bb1, bbp2) {
				continue
			}
			// Bounding boxes of lines overlap - see if they intersect
			s2, e2 := fparts2[j][0], fparts2[j][len(fparts2[j])-1]
			tvals, err := util.IntersectionTValsP(s1, e1, s2, e2)
			if err != nil || tvals[0] < 0 || tvals[0] > 1 || tvals[1] < 0 || tvals[1] > 1 {
				continue
			}
			return []float64{util.Lerp(tvals[0], s1[0], e1[0]), util.Lerp(tvals[0], s1[1], e1[1])}
		}
	}

	return nil
}

// Length returns the approximate length of a path by flattening it to the desired degree
// and summing the line steps.
func (p *Path) Length(flat float64) float64 {
	parts := p.Flatten(flat).Parts()
	sum := 0.0
	for _, part := range parts {
		sum += util.DistanceE(part[0], part[1])
	}
	return sum
}

// ProjectPoint returns the point on the path closest to pt.
func (p *Path) ProjectPoint(pt []float64) []float64 {
	sp := p.Simplify()
	parts := sp.Parts()
	n := len(parts)

	// Iterate through the parts of the simplified path to
	// find the closest parts.
	d := make([]float64, n+1)
	cp := 0
	cd := dist2(pt, parts[0][0])
	d[0] = cd
	for i := 1; i < n; i++ {
		d2 := dist2(pt, parts[i][0])
		d[i] = d2
		if d2 < cd {
			cp, cd = i, d2
		}
	}
	// Check end point
	d2 := dist2(pt, parts[n-1][len(parts[n-1])])
	d[n] = d2
	if d2 < cd {
		cp = n
	}

	if cp == 0 {
		pr, dr := bs(pt, 0, d[0], 1, d[1], parts[0])
		if !p.closed {
			// point lies in first part
			return pr
		}
		// else test parts[n-1]
		pl, dl := bs(pt, 0, d[n-1], 1, d[n], parts[n-1])
		if dl < dr {
			return pl
		}
		return pr
	}

	if cp == n {
		pl, dl := bs(pt, 0, d[n-1], 1, d[n], parts[n-1])
		if !p.closed {
			// point lies in last part
			return pl
		}
		// else test parts[0]
		pr, dr := bs(pt, 0, d[0], 1, d[1], parts[0])
		if dl < dr {
			return pl
		}
		return pr
	}

	// point lies in either cp-1 to ci or, ci to cp+1
	pl, dl := bs(pt, 0, d[cp-1], 1, d[cp], parts[cp-1])
	pr, dr := bs(pt, 0, d[cp], 1, d[cp+1], parts[cp])

	if dl < dr {
		return pl
	}

	return pr
}

func dist2(a, b []float64) float64 {
	dx, dy := b[0]-a[0], b[1]-a[1]
	return dx*dx + dy*dy
}

// Returns closest point on part to pt and dist2
func bs(pt []float64, ts, ds, te, de float64, part [][]float64) ([]float64, float64) {
	td := te - ts
	if td < 0.00001 {
		return util.DeCasteljau(part, ts), ds
	}
	t := []float64{ts, ts + td/4, ts + td/2, ts + 3*td/4, te}
	dl := dist2(pt, util.DeCasteljau(part, t[1]))
	dm := dist2(pt, util.DeCasteljau(part, t[2]))
	dr := dist2(pt, util.DeCasteljau(part, t[3]))
	d := []float64{ds, dl, dm, dr, de}
	ci, cd := 0, d[0]
	for i := 1; i < 5; i++ {
		if d[i] < cd {
			ci = i
		}
	}
	if ci == 0 {
		// Search t[0] to t[1]
		return bs(pt, t[0], d[0], t[1], d[1], part)
	}
	if ci == 4 {
		// Search t[3] to t[4]
		return bs(pt, t[3], d[3], t[4], d[4], part)
	}
	// Search t[ci-1] to t[ci] and t[ci] to t[ci+1]
	pl, d1 := bs(pt, t[ci-1], d[ci-1], t[ci], d[ci], part)
	pr, d2 := bs(pt, t[ci], d[ci], t[ci+1], d[ci+1], part)
	if d1 < d2 {
		return pl, d1
	}
	return pr, d2
}

// PointInPath returns if a point is contained within a closed path according to the
// setting of util.WindingRule. If the path is not closed then false is returned, regardless.
func (p *Path) PointInPath(pt []float64) bool {
	return util.PointInPoly(pt, p.Poly()...)
}

// Poly converts a path into a flat sided polygon. Returns an empty slice if the path isn't closed.
func (p *Path) Poly() [][]float64 {
	if !p.closed {
		return [][]float64{}
	}
	fp := p.Flatten(RenderFlatten)
	parts := fp.Parts()
	poly := make([][]float64, len(parts))
	for i, part := range parts {
		poly[i] = part[0]
	}
	return poly
}

type pathJSON struct {
	Steps  [][][]float64
	Closed bool
}

// MarshalJSON implements the encoding.json.Marshaler interface
func (p *Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(pathJSON{p.steps, p.closed})
}

// UnmarshalJSON implements the encoding.json.Unmarshaler interface
func (p *Path) UnmarshalJSON(b []byte) error {
	var pj pathJSON
	err := json.Unmarshal(b, &pj)
	if err != nil {
		return err
	}
	p.steps = pj.Steps
	p.closed = pj.Closed

	// Reset everything else
	p.bounds = image.Rectangle{}
	p.flattened = nil
	p.tolerance = 0
	p.simplified = nil
	p.reversed = nil
	p.tangents = nil
	p.parent = nil

	return nil
}
