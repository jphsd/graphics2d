package graphics2d

import (
	"encoding/json"
	"encoding/xml"
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

// Start returns the first point in the path.
func (p *Path) Start() []float64 {
	pt := p.steps[0][0]
	return []float64{pt[0], pt[1]}
}

// Current returns the last point in the path.
func (p *Path) Current() []float64 {
	var pt []float64
	if p.closed {
		pt = p.steps[0][0]
	} else {
		last := p.steps[len(p.steps)-1]
		pt = last[len(last)-1]
	}
	return []float64{pt[0], pt[1]}
}

// Concatenate adds the paths to this path. If any path is closed then an error
// is returned. If the paths aren't coincident, then they are joined with a line.
func (p *Path) Concatenate(paths ...*Path) error {
	if p.closed {
		return fmt.Errorf("path is closed, adding a step is forbidden")
	}
	for _, path := range paths {
		if path != nil && path.closed {
			return fmt.Errorf("can't add a closed path")
		}
	}

	for _, path := range paths {
		if path == nil {
			continue
		}
		last := p.Current()
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
	}
	return nil
}

// ConcatenatePaths concatenates all the paths into a new path. If any path is closed then an error
// is returned. If the paths aren't coincident, then they are joined with a line.
func ConcatenatePaths(paths ...*Path) (*Path, error) {
	if len(paths) == 0 || paths[0] == nil {
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

// AddParts adds parts to the path and returns it.
func (p *Path) AddParts(parts ...Part) *Path {
	for i, part := range parts {
		for j, pt := range part {
			if InvalidPoint(pt) {
				panic(fmt.Sprintf("invalid point (NaN) in part %d,%d", i, j))
			}
		}
		p.AddStep(part[1:]...)
	}
	return p
}

// Parts returns the steps of a path, each prepended with its start.
func (p *Path) Parts() Parts {
	n := len(p.steps)
	if n == 0 {
		return nil
	}
	cp := p.steps[0][0]
	if n == 1 {
		// This is a point
		// Deep copy
		return Parts{{{cp[0], cp[1]}}}
	}
	parts := make(Parts, n-1, n)
	for i := 1; i < n; i++ {
		part := toPart(cp, p.steps[i])
		parts[i-1] = part
		cp = part[len(part)-1]
	}
	if p.closed && !util.EqualsP(cp, p.steps[0][0]) {
		parts = append(parts, Part{cp, p.steps[0][0]})
	}
	return parts
}

// Deep copy
func toPart(cp []float64, pts [][]float64) Part {
	res := make(Part, len(pts)+1)
	res[0] = []float64{cp[0], cp[1]}
	for i, pt := range pts {
		res[i+1] = []float64{pt[0], pt[1]}
	}
	return res
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
		fp := toPart(cp, sp.steps[i]).flatten(d2)
		for _, ns := range fp {
			// ns length is always 2
			res = append(res, Part{ns[1]})
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
func (p *Path) String() string {
	b, _ := p.MarshalText()
	return string(b)
}

// Simplify breaks up a path into steps where for any step, its control points are all on the
// same side and its midpoint is well behaved. If a step doesn't meet the criteria, it is
// recursively subdivided in half until it does.
func (p *Path) Simplify() *Path {
	if p.simplified != nil {
		return p.simplified
	}
	res := Parts{}
	parts := p.Parts()
	for _, part := range parts {
		if len(part) > 2 {
			// Split based on extremities
			nparts := part.SimplifyExtremities()
			for _, npart := range nparts {
				// check simplified parts are simple
				res = append(res, npart.Simplify()...)
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

// Reverse returns a new path describing the current path in reverse order (i.e start and end switched).
func (p *Path) Reverse() *Path {
	path := PartsToPath(p.Parts().Reverse()...)

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

// Length returns the approximate length of a path by flattening it to the desired degree
// and summing the line steps.
func (p *Path) Length(flat float64) float64 {
	parts := p.Parts()
	sum := 0.0
	for _, part := range parts {
		sum += part.Length(flat)
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
		pp, d, t := part.bs(pt, 0, d[i], 1, d[i+1])
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

/*
 * Marshaling functions for JSON, XML/SVG and text.
 */

type jpath struct {
	Steps  [][][]float64
	Closed bool
}

// MarshalJSON implements the encoding/json.Marshaler interface
func (p *Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(jpath{p.steps, p.closed})
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface
func (p *Path) UnmarshalJSON(b []byte) error {
	var pj jpath
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

type xpath struct {
	Desc string `xml:"d,attr"`
}

// MarshalXML implements the encoding/xml.Marshaler interface
func (p *Path) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(xpath{p.StringSVG()}, xml.StartElement{Name: xml.Name{Space: "", Local: "path"}})
}

func (p *Path) StringSVG() string {
	// SVG can't handle high order steps
	fp := p.Flatten(DefaultRenderFlatten)
	pt := fp.steps[0][0]
	desc := fmt.Sprintf("M %.2f %.2f", pt[0], pt[1])
	ns := len(fp.steps)
	if ns == 1 {
		if p.closed {
			desc += " z"
		}
		return desc
	}
	desc += " L"
	for i := 1; i < len(fp.steps); i++ {
		pt := fp.steps[i][0]
		desc += fmt.Sprintf(" %.2f %.2f", pt[0], pt[1])
	}
	if p.closed {
		desc += " z"
	}
	return desc
}

// UnmarshalXML is not supported.
// Use the github.com/jphsd/xml/svg framework instead.
func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return fmt.Errorf("UnmarshalXML is not supported")
}

// MarshalText implements the encoding.TextMarshaler interface.
// P %f,%f[ %d[ %f,%f]][ C]
func (p *Path) MarshalText() ([]byte, error) {
	step := p.steps[0]
	var str strings.Builder
	str.WriteString(fmt.Sprintf("P %f,%f", step[0][0], step[0][1]))
	for i := 1; i < len(p.steps); i++ {
		step = p.steps[i]
		str.WriteString(fmt.Sprintf(" %d", len(step)))
		for _, pts := range step {
			str.WriteString(fmt.Sprintf(" %f,%f", pts[0], pts[1]))
		}
	}
	if p.closed {
		str.WriteString(" C")
	}
	return []byte(str.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Expected format: P %f,%f[ %d[ %f,%f]][ C]
func (p *Path) UnmarshalText(b []byte) error {
	str := string(b)

	parts := strings.Split(str, " ")
	np := len(parts)

	if np < 2 || parts[0] != "P" {
		return fmt.Errorf("Not a valid Path string")
	}

	var x, y float64
	n, err := fmt.Sscanf(parts[1], "%f,%f", &x, &y)
	if n != 2 || err != nil {
		return fmt.Errorf("Not a valid Path string")
	}
	steps := [][][]float64{{{x, y}}}

	// Handle steps
	for i := 2; i < np; {
		if parts[i] == "C" {
			p.steps = steps
			p.closed = true
			return nil
		}
		var s int
		n, err = fmt.Sscanf(parts[i], "%d", &s)
		if n != 1 || err != nil {
			return fmt.Errorf("Not a valid Path string")
		}
		step := make([][]float64, s)
		i++
		for j := 0; i < np && j < s; j++ {
			n, err := fmt.Sscanf(parts[i], "%f,%f", &x, &y)
			if n != 2 || err != nil {
				return fmt.Errorf("Not a valid Path string")
			}
			step[j] = []float64{x, y}
			i++
		}
		steps = append(steps, step)
	}
	p.steps = steps

	// Reset everything else
	p.bbox = nil
	p.flattened = nil
	p.tolerance = 0
	p.simplified = nil
	p.tangents = nil
	p.parent = nil

	return nil
}
