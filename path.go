package graphics2d

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"strings"

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
	bbox   [][]float64
	// Caching flattened, simplified and reversed paths, and tangents
	flattened  *Path
	tolerance  float64
	simplified *Path
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
// Adding a step to a closed path will cause an error as will adding an invalid point.
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
	nstep := make([][]float64, 0, n)
	for i, pt := range points {
		if InvalidPoint(pt) {
			return fmt.Errorf("invalid step point (NaN) at %d", i)
		}
		if util.EqualsP(last, pt) {
			// Ignore coincident points
			continue
		}
		nstep = append(nstep, pt)
		last = pt
	}
	if len(nstep) == 0 {
		// Nothing added
		return nil
	}

	p.steps = append(p.steps, nstep)
	// Reset cached items
	p.bbox = nil
	p.flattened = nil
	p.simplified = nil
	p.tangents = nil
	return nil
}

// LineTo is a chain wrapper around AddStep.
func (p *Path) LineTo(point []float64) *Path {
	p.AddStep(point)
	return p
}

// CurveTo is a chain wrapper around AddStep.
func (p *Path) CurveTo(points ...[]float64) *Path {
	p.AddStep(points...)
	return p
}

// ArcTo is a chain wrapper around MakeRoundedParts.
// If r is too large for the supplied tangents, then it is truncated.
func (p *Path) ArcTo(p1, p2 []float64, r float64) *Path {
	last := p.steps[len(p.steps)-1]
	p0 := last[len(last)-1]
	parts := MakeRoundedParts(p0, p1, p2, r)
	p.AddStep(parts[0][0]) // in case arc doesn't start at p0
	for _, part := range parts {
		p.AddStep(part[1:]...)
	}
	p.AddStep(p2) // in case arc doesn't end at p2
	return p
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
			for _, step := range steps[1:] {
				p.AddStep(step...)
			}
		} else {
			// Line to steps[0][0]
			for _, step := range steps {
				p.AddStep(step...)
			}
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
	if n == 0 {
		return nil
	}
	cp := p.steps[0][0]
	if n == 1 {
		// This is a point
		// Deep copy
		return [][][]float64{{{cp[0], cp[1]}}}
	}
	parts := make([][][]float64, n-1, n)
	for i := 1; i < n; i++ {
		part := toPart(cp, p.steps[i])
		parts[i-1] = part
		cp = part[len(part)-1]
	}
	if p.closed && !util.EqualsP(cp, p.steps[0][0]) {
		parts = append(parts, [][]float64{cp, p.steps[0][0]})
	}
	return parts
}

// Close marks the path as closed.
func (p *Path) Close() *Path {
	if p.closed {
		return p
	}
	p.AddStep(p.steps[0][0])
	p.closed = true
	return p
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
	d2 := d * d
	res := make([][][]float64, 1)
	sp := p.Simplify()
	res[0] = sp.steps[0]
	cp := res[0][0]
	// For all remaining steps in path
	for i := 1; i < len(sp.steps); i++ {
		fp := flattenPart(d2, toPart(cp, sp.steps[i]))
		for _, ns := range fp {
			// ns length is always 2
			res = append(res, [][]float64{ns[1]})
			cp = ns[1]
		}
	}

	path := &Path{}
	path.steps = res
	path.closed = p.closed
	path.parent = p
	p.flattened = path
	return path
}

// Deep copy
func toPart(cp []float64, pts [][]float64) [][]float64 {
	res := make([][]float64, len(pts)+1)
	res[0] = []float64{cp[0], cp[1]}
	for i, pt := range pts {
		res[i+1] = []float64{pt[0], pt[1]}
	}
	return res
}

// flattenPart successively splits the part until the control points are within d2 of the line
// joing the part start with the part end. Returns a list of line parts.
func flattenPart(d2 float64, part [][]float64) [][][]float64 {
	if cpWithinD2(d2, part) {
		return [][][]float64{{part[0], part[len(part)-1]}}
	}
	lr := util.SplitCurve(part, 0.5)
	res := append([][][]float64{}, flattenPart(d2, lr[0])...)
	res = append(res, flattenPart(d2, lr[1])...)
	return res
}

func cpWithinD2(d2 float64, part [][]float64) bool {
	// First and last are end points
	l := len(part)
	if l == 2 {
		// Trivial case
		return true
	}
	start, cpts, end := part[0], part[1:l-1], part[l-1]
	for _, cp := range cpts {
		pd2, _, _ := util.DistanceToLineSquared(start, end, cp)
		if pd2 > d2 {
			return false
		}
	}

	return true
}

// FlattenPart works by subdividing the curve until its control points are within d^2 (d squared)
// of the line through the end points.
func FlattenPart(d float64, part [][]float64) [][][]float64 {
	return flattenPart(d*d, part)
}

// PartLength returns the approximate length of a part by flattening it to the supplied degree of flatness.
func PartLength(d float64, part [][]float64) float64 {
	parts := flattenPart(d*d, part)
	sum := 0.0
	for _, part := range parts {
		sum += util.DistanceE(part[0], part[1])
	}
	return sum
}

// BoundingBox calculates a bounding box that the Path is guaranteed to fit within. It's unlikely to
// be the minimal bounding box for the path since the control points are also included.
// If a tight bounding box is required then use CalcExtremities().
func (p *Path) BoundingBox() [][]float64 {
	if p.bbox == nil {
		bb := [][]float64{p.steps[0][0], p.steps[0][0]}
		for _, step := range p.steps {
			bbp := util.BoundingBox(step...)
			bb = util.BoundingBox(bb[0], bb[1], bbp[0], bbp[1])
		}
		p.bbox = bb
	}
	return p.bbox
}

// Bounds calculates a rectangle that the Path is guaranteed to fit within. It's unlikely to
// be the minimal bounding rectangle for the path since the control points are also included.
// If a tight bounding rectangle is required then use CalcExtremities().
func (p *Path) Bounds() image.Rectangle {
	return util.BBToRect(p.BoundingBox())
}

// Copy performs a deep copy
func (p *Path) Copy() *Path {
	steps := make([][][]float64, len(p.steps))
	copy(steps, p.steps)

	path := &Path{}
	path.steps = make([][][]float64, len(steps))
	for i, step := range steps {
		path.steps[i] = copyStep(step)
	}
	path.closed = p.closed
	path.parent = p.parent
	return path
}

func copyStep(step [][]float64) [][]float64 {
	res := make([][]float64, len(step))
	for i, pt := range step {
		res[i] = copyPoint(pt)
	}
	return res
}

func copyPoint(point []float64) []float64 {
	// Only preserve x and y
	return []float64{point[0], point[1]}
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
		if path != nil {
			path.parent = p
		}
	}

	return paths
}

// String converts a path into a string.
// P %f,%f [S %d [%f,%f ]][C]
func (p *Path) String() string {
	step := p.steps[0]
	str := fmt.Sprintf("P %f,%f ", step[0][0], step[0][1])
	for i := 1; i < len(p.steps); i++ {
		step = p.steps[i]
		str += "S "
		str += fmt.Sprintf("%d ", len(step))
		for _, pts := range step {
			str += fmt.Sprintf("%f,%f ", pts[0], pts[1])
		}
	}
	if p.closed {
		str += "C"
	}
	return str
}

// StringToPath converts a string created using path.String() back into a path.
// Returns nil if the string isn't parsable into a path.
func StringToPath(str string) *Path {
	parts := strings.Split(str, " ")
	np := len(parts)

	if np < 2 || parts[0] != "P" {
		return nil
	}
	var x, y float64
	n, err := fmt.Sscanf(parts[1], "%f,%f", &x, &y)
	if n != 2 || err != nil {
		return nil
	}
	path := NewPath([]float64{x, y})
	if np < 3 {
		return path
	}

	// Handle steps
	i := 2
	for i+1 < np && parts[i] == "S" {
		i++
		var s int
		n, err = fmt.Sscanf(parts[i], "%d", &s)
		if n != 1 || err != nil {
			return nil
		}
		step := make([][]float64, s)
		i++
		for j := 0; i < np && j < s; j++ {
			n, err := fmt.Sscanf(parts[i], "%f,%f", &x, &y)
			if n != 2 || err != nil {
				return nil
			}
			step[j] = []float64{x, y}
			i++
		}
		path.AddStep(step...)
	}

	if i < np && parts[i] == "C" {
		path.Close()
	}

	return path
}

// Simplify breaks up a path into steps where for any step, its control points are all on the
// same side and its midpoint is well behaved. If a step doesn't meet the criteria, it is
// recursively subdivided in half until it does.
func (p *Path) Simplify() *Path {
	if p.simplified != nil {
		return p.simplified
	}
	res := [][][]float64{}
	parts := p.Parts()
	for _, part := range parts {
		if len(part) > 2 {
			// Split based on extremities
			nparts := SimplifyExtremities(part)
			for _, npart := range nparts {
				// check simplified parts are simple
				res = append(res, SimplifyPart(npart)...)
			}
		} else {
			res = append(res, part)
		}
	}

	path := PartsToPath(res...)
	path.closed = p.closed
	path.parent = p
	path.simplified = path // self reference
	p.simplified = path

	return path
}

// SimplifyExtremities chops curve into pieces based on maxima, minima and inflections in x and y.
func SimplifyExtremities(part [][]float64) [][][]float64 {
	tvals := util.CalcExtremities(part)
	nt := len(tvals)
	if nt < 3 {
		return [][][]float64{part}
	}
	// Convert tvals to relative tvals
	rtvals := make([]float64, nt)
	lt := 0.0
	ll := 1.0
	for i := range nt {
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

	rhs := part
	res := make([][][]float64, nt-1)
	for i := 1; i < nt-1; i++ {
		lr := util.SplitCurve(rhs, rtvals[i])
		res[i-1] = lr[0]
		rhs = lr[1]
	}
	res[nt-2] = rhs
	return res
}

// SimplifyPart recursively cuts the curve in half until CPSafe is
// satisfied.
func SimplifyPart(part [][]float64) [][][]float64 {
	if CPSafe(part) {
		return [][][]float64{part}
	}
	lr := util.SplitCurve(part, 0.5)
	res := append([][][]float64{}, SimplifyPart(lr[0])...)
	res = append(res, SimplifyPart(lr[1])...)
	return res
}

// SafeFraction if greater than 0 causes Simplify to perform a check of the mid-point against
// the part centroid. If the two are within SafeFraction of the distance from p[0] to the centroid
// then no further subdivision of the curve is performed.
var SafeFraction float64 = -1

// CPSafe returns true if all the control points are on the same side of
// the line formed by start and the last part points and the point at t = 0.5 is close
// to the centroid of the part.
func CPSafe(part [][]float64) bool {
	n := len(part)
	if n < 3 {
		// Either a point or line
		return true
	}

	start := part[0]
	end := part[n-1]
	side := util.CrossProduct(start, end, part[1]) < 0
	for i := 2; i < n-1; i++ {
		if (util.CrossProduct(start, end, part[i]) < 0) != side {
			return false
		}
	}

	if n == 3 {
		return true
	}

	if SafeFraction > 0 {
		// Check mid-point against centroid
		// scale against distance between p0 and centroid
		centroid := util.Centroid(part...)
		hp := util.DeCasteljau(part, 0.5)
		p0dx := centroid[0] - part[0][0]
		p0dy := centroid[1] - part[0][1]
		p0ds := p0dx*p0dx + p0dy*p0dy
		hpdx := centroid[0] - hp[0]
		hpdy := centroid[1] - hp[1]
		hpds := hpdx*hpdx + hpdy*hpdy
		return math.Sqrt(hpds) < math.Sqrt(p0ds)*SafeFraction
	}

	return true
}

// Reverse returns a new path describing the current path in reverse order (i.e start and end switched).
func (p *Path) Reverse() *Path {
	path := PartsToPath(ReverseParts(p.Parts())...)

	// If other aspects have already been calculated - reverse them too
	if p.flattened != nil {
		path.flattened = p.flattened.Reverse()
		path.tolerance = p.tolerance
	}
	if p.simplified != nil {
		if p != p.simplified {
			path.simplified = p.simplified.Reverse()
		} else {
			path.simplified = path
		}
	}
	if p.tangents != nil {
		path.tangents = reverseTangents(p.tangents)
	}
	if p.closed {
		path.closed = true
	}
	path.parent = p.parent

	return path
}

// ReverseParts reverses the order (and points) of the supplied part slice.
func ReverseParts(parts [][][]float64) [][][]float64 {
	n := len(parts)
	res := make([][][]float64, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = ReversePoints(parts[j])
		j--
	}
	return res
}

// [pts][x/y]
func ReversePoints(cp [][]float64) [][]float64 {
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

// unit converts a normal to a unit normal
func unit(dx, dy float64) (float64, float64) {
	d := math.Hypot(dx, dy)
	if util.Equals(0, d) {
		return 0, 0
	}
	return dx / d, dy / d
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
	parts := p.Parts()
	sum := 0.0
	for _, part := range parts {
		sum += PartLength(flat, part)
	}
	return sum
}

// ProjectPoint returns the point, it's t on the path closest to pt and the distance^2.
// Note t can be very non-linear.
func (p *Path) ProjectPoint(pt []float64) ([]float64, float64, float64) {
	sp := p.Simplify()
	parts := sp.Parts()
	n := len(parts)
	dtp := 1.0 / float64(n)

	// Construct start and end distances
	d := make([]float64, n+1)
	for i := range n {
		d[i] = dist2(pt, parts[i][0])
	}
	if p.closed {
		d[n] = d[0]
	} else {
		d[n] = dist2(pt, parts[n-1][len(parts[n-1])-1])
	}

	// Iterate through all parts of the simplified path, since there may be crossing
	// points.
	c, ppt, bt, bd := -1, []float64{}, 0.0, math.MaxFloat64 // Best fit so far
	for i, part := range parts {
		pp, d, t := bs(pt, 0, d[i], 1, d[i+1], part)
		if d < bd {
			bd = d
			c = i
			bt = t
			ppt = pp
		}
	}

	return ppt, (float64(c) + bt) * dtp, bd
}

func dist2(a, b []float64) float64 {
	dx, dy := b[0]-a[0], b[1]-a[1]
	return dx*dx + dy*dy
}

// Returns closest point on part to pt, dist2 and t [0-1]
func bs(pt []float64, ts, ds, te, de float64, part [][]float64) ([]float64, float64, float64) {
	dt := te - ts
	if dt < 0.00001 {
		prj := util.DeCasteljau(part, ts)
		return prj, ds, ts
	}

	t := []float64{ts, ts + dt/4, ts + dt/2, ts + 3*dt/4, te}
	dl := dist2(pt, util.DeCasteljau(part, t[1]))
	dm := dist2(pt, util.DeCasteljau(part, t[2]))
	dr := dist2(pt, util.DeCasteljau(part, t[3]))
	d := []float64{ds, dl, dm, dr, de}
	nd := len(d)

	ci, cd := 0, d[0]
	for i := 1; i < nd; i++ {
		if d[i] < cd {
			ci = i
			cd = d[i]
		}
	}
	if ci == 0 {
		// Only search t[0] to t[1]
		return bs(pt, t[0], d[0], t[1], d[1], part)
	}
	if ci == nd-1 {
		// Only search t[3] to t[4]
		return bs(pt, t[3], d[3], t[4], d[4], part)
	}

	// Search both t[ci-1] to t[ci] and t[ci] to t[ci+1]
	pl, d1, tl := bs(pt, t[ci-1], d[ci-1], t[ci], d[ci], part)
	pr, d2, tr := bs(pt, t[ci], d[ci], t[ci+1], d[ci+1], part)
	if d1 < d2 {
		return pl, d1, tl
	}
	return pr, d2, tr
}

// PointInPath returns if a point is contained within a closed path according to the
// setting of util.WindingRule. If the path is not closed then false is returned, regardless.
func (p *Path) PointInPath(pt []float64) bool {
	if !p.closed {
		return false
	}
	ppts, _ := p.PolyLine()
	return util.PointInPoly(pt, ppts...)
}

// PolyLine converts a path into a polygon line. If the second result is true, the result is a polygon.
func (p *Path) PolyLine() ([][]float64, bool) {
	fp := p.Flatten(RenderFlatten)
	parts := fp.Parts()
	np := len(parts)
	poly := make([][]float64, np)
	for i, part := range parts {
		poly[i] = part[0]
		if i == np-1 && !p.closed {
			poly = append(poly, part[len(part)-1])
		}
	}
	return poly, p.closed
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
	p.bbox = nil
	p.flattened = nil
	p.tolerance = 0
	p.simplified = nil
	p.tangents = nil
	p.parent = nil

	return nil
}
