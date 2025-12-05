package graphics2d

import (
	"fmt"
	"math"

	"github.com/jphsd/graphics2d/util"
)

// Aff3 is a 3x3 affine transformation matrix in row major order,
// where the bottom row is implicitly [0 0 1].
//
// m[3*r+c] is the element in the r'th row and c'th column.
type Aff3 [6]float64

// NewAff3 creates the identity affine transform.
func NewAff3() *Aff3 {
	var res Aff3
	res[3*0+0] = 1
	res[3*1+1] = 1
	return &res
}

// Translate creates an affine translation transform.
func Translate(x, y float64) *Aff3 {
	xfm := NewAff3()
	xfm.Translate(x, y)
	return xfm
}

// Rotate creates an affine rotation transform.
func Rotate(th float64) *Aff3 {
	xfm := NewAff3()
	xfm.Rotate(th)
	return xfm
}

// RotateAbout creates an affine rotation transform about a point.
func RotateAbout(th, ax, ay float64) *Aff3 {
	xfm := NewAff3()
	xfm.RotateAbout(th, ax, ay)
	return xfm
}

// Scale creates an affine scale transform.
func Scale(sx, sy float64) *Aff3 {
	xfm := NewAff3()
	xfm.Scale(sx, sy)
	return xfm
}

// ScaleAbout creates an affine scale transform about a point.
func ScaleAbout(sx, sy, ax, ay float64) *Aff3 {
	xfm := NewAff3()
	xfm.ScaleAbout(sx, sy, ax, ay)
	return xfm
}

// Shear creates an affine shear transform.
func Shear(shx, shy float64) *Aff3 {
	xfm := NewAff3()
	xfm.Shear(shx, shy)
	return xfm
}

// ShearAbout creates an affine shear transform about a point.
func ShearAbout(shx, shy, ax, ay float64) *Aff3 {
	xfm := NewAff3()
	xfm.ShearAbout(shx, shy, ax, ay)
	return xfm
}

// Reflect creates an affine reflection transform.
func Reflect(x1, y1, x2, y2 float64) *Aff3 {
	xfm := NewAff3()
	xfm.Reflect(x1, y1, x2, y2)
	return xfm
}

// Determinant calculates the affine transform's matrix determinant.
func (a *Aff3) Determinant() float64 {
	return a[3*0+0]*a[3*1+1] - a[3*0+1]*a[3*1+0]
}

// Translate adds a translation to the affine transform.
func (a *Aff3) Translate(x, y float64) *Aff3 {
	a[3*0+2] = x*a[3*0+0] + y*a[3*0+1] + a[3*0+2]
	a[3*1+2] = x*a[3*1+0] + y*a[3*1+1] + a[3*1+2]
	return a
}

// Rotate adds a rotation to the affine transform. The rotation is about {0, 0}.
func (a *Aff3) Rotate(th float64) *Aff3 {
	sin, cos := math.Sin(th), math.Cos(th)
	m0, m1 := a[3*0+0], a[3*0+1]
	a[3*0+0] = cos*m0 + sin*m1
	a[3*0+1] = -sin*m0 + cos*m1
	m0, m1 = a[3*1+0], a[3*1+1]
	a[3*1+0] = cos*m0 + sin*m1
	a[3*1+1] = -sin*m0 + cos*m1
	return a
}

// RotateAbout adds a rotation about a point to the affine transform.
func (a *Aff3) RotateAbout(th, ax, ay float64) *Aff3 {
	// Reverse order
	a.Translate(ax, ay)
	a.Rotate(th)
	a.Translate(-ax, -ay)
	return a
}

// QuadrantRotate adds a rotation (n * 90 degrees) to the affine transform. The rotation is about {0, 0}.
// It avoids rounding issues with the trig functions.
func (a *Aff3) QuadrantRotate(n int) *Aff3 {
	n %= 4
	switch n {
	case 0: // 360
		break
	case 1: // 90
		a[3*0+0], a[3*0+1], a[3*1+0], a[3*1+1] = a[3*0+1], -a[3*0+0], a[3*1+1], -a[3*1+0]
	case 2: // 180
		a[3*0+0], a[3*0+1], a[3*1+0], a[3*1+1] = -a[3*0+0], -a[3*0+1], -a[3*1+0], -a[3*1+1]
	case 3: // 270
		a[3*0+0], a[3*0+1], a[3*1+0], a[3*1+1] = -a[3*0+1], a[3*0+0], -a[3*1+1], a[3*1+0]
	}
	return a
}

// QuadrantRotateAbout adds a rotation (n * 90 degrees) about a point to the affine transform.
// It avoids rounding issues with the trig functions.
func (a *Aff3) QuadrantRotateAbout(n int, ax, ay float64) *Aff3 {
	// Reverse order
	a.Translate(ax, ay)
	a.QuadrantRotate(n)
	a.Translate(-ax, -ay)
	return a
}

// Scale adds a scaling to the affine transform centered on {0, 0}.
func (a *Aff3) Scale(sx, sy float64) *Aff3 {
	a[3*0+0] *= sx
	a[3*1+1] *= sy
	a[3*0+1] *= sy
	a[3*1+0] *= sx
	return a
}

// ScaleAbout adds a scale about a point to the affine transform.
func (a *Aff3) ScaleAbout(sx, sy, ax, ay float64) *Aff3 {
	// Reverse order
	a.Translate(ax, ay)
	a.Scale(sx, sy)
	a.Translate(-ax, -ay)
	return a
}

// Shear adds a shear to the affine transform centered on {0, 0}.
func (a *Aff3) Shear(shx, shy float64) *Aff3 {
	m0, m1 := a[3*0+0], a[3*0+1]
	a[3*0+0] = m0 + m1*shy
	a[3*0+1] = m0*shx + m1

	m0, m1 = a[3*1+0], a[3*1+1]
	a[3*1+0] = m0 + m1*shy
	a[3*1+1] = m0*shx + m1
	return a
}

// ShearAbout adds a shear about a point to the affine transform.
func (a *Aff3) ShearAbout(shx, shy, ax, ay float64) *Aff3 {
	// Reverse order
	a.Translate(ax, ay)
	a.Shear(shx, shy)
	a.Translate(-ax, -ay)
	return a
}

// Concatenate concatenates an affine transform to the affine transform.
func (a *Aff3) Concatenate(aff Aff3) *Aff3 {
	m00, m01, m10, m11 := a[3*0+0], a[3*0+1], a[3*1+0], a[3*1+1]
	t00, t01, t02, t10, t11, t12 := aff[3*0+0], aff[3*0+1], aff[3*0+2], aff[3*1+0], aff[3*1+1], aff[3*1+2]

	a[3*0+0] = t00*m00 + t10*m01
	a[3*0+1] = t01*m00 + t11*m01
	a[3*0+2] += t02*m00 + t12*m01

	a[3*1+0] = t00*m10 + t10*m11
	a[3*1+1] = t01*m10 + t11*m11
	a[3*1+2] += t02*m10 + t12*m11
	return a
}

// PreConcatenate preconcatenates an affine transform to the affine transform.
func (a *Aff3) PreConcatenate(aff Aff3) *Aff3 {
	m00, m01, m02, m10, m11, m12 := a[3*0+0], a[3*0+1], a[3*0+2], a[3*1+0], a[3*1+1], a[3*1+2]
	t00, t01, t02, t10, t11, t12 := aff[3*0+0], aff[3*0+1], aff[3*0+2], aff[3*1+0], aff[3*1+1], aff[3*1+2]

	t02 += m02*t00 + m12*t01
	t12 += m02*t10 + m12*t11
	a[3*0+2] = t02
	a[3*1+2] = t12

	a[3*0+0] = m00*t00 + m10*t01
	a[3*1+0] = m00*t10 + m10*t11

	a[3*0+1] = m01*t00 + m11*t01
	a[3*1+1] = m01*t10 + m11*t11
	return a
}

// InverseOf returns the inverse of the affine transform.
func (a *Aff3) InverseOf() (*Aff3, error) {
	res := a.Copy()
	if err := res.Invert(); err != nil {
		return nil, err
	}

	return res, nil
}

// Invert inverts the affine transform.
func (a *Aff3) Invert() error {
	det := a.Determinant()
	if util.Equals(math.Abs(det), 0) {
		return fmt.Errorf("Determinant is zero => non-invertible")
	}

	m00, m01, m02, m10, m11, m12 := a[3*0+0], a[3*0+1], a[3*0+2], a[3*1+0], a[3*1+1], a[3*1+2]
	a[3*0+0], a[3*0+1], a[3*0+2] = m11/det, -m01/det, (m01*m12-m02*m11)/det
	a[3*1+0], a[3*1+1], a[3*1+2] = -m10/det, m00/det, (m02*m10-m00*m12)/det
	return nil
}

// String converts the affine transform into a string.
func (a *Aff3) String() string {
	return fmt.Sprintf("{{%f, %f, %f}, {%f, %f, %f}, {0, 0, 1}}",
		a[3*0+0], a[3*0+1], a[3*0+2],
		a[3*1+0], a[3*1+1], a[3*1+2])
}

// Identity returns true if the affine transform is the identity.
func (a *Aff3) Identity() bool {
	if !util.Equals(a[3*0+0], 1) {
		return false
	}
	if !util.Equals(a[3*0+1], 0) {
		return false
	}
	if !util.Equals(a[3*0+2], 0) {
		return false
	}
	if !util.Equals(a[3*1+0], 0) {
		return false
	}
	if !util.Equals(a[3*1+1], 1) {
		return false
	}
	if !util.Equals(a[3*1+2], 0) {
		return false
	}
	return true
}

// Copy returns a copy of the affine transform.
func (a *Aff3) Copy() *Aff3 {
	res := *a
	return &res
}

// Reflect performs a reflection along the axis defined by the two non-coincident points.
func (a *Aff3) Reflect(x1, y1, x2, y2 float64) *Aff3 {
	dx, dy := x2-x1, y2-y1
	if util.Equals(dy, 0) {
		// Horizontal - no rotation required
		a.Translate(0, y1)
		a.Scale(1, -1)
		a.Translate(0, -y1)
		return a
	}
	if util.Equals(dx, 0) {
		// Vertical - no rotation required
		a.Translate(x1, 0)
		a.Scale(-1, 1)
		a.Translate(-x1, 0)
		return a
	}
	th := math.Atan2(dy, dx)
	if th < 0 {
		th += TwoPi
	}
	a.Translate(x1, y1)
	a.Rotate(th)
	a.Scale(1, -1)
	a.Rotate(-th)
	a.Translate(-x1, -y1)
	return a
}

// LineTransform produces an affine transform that maps the line {p1, p2} to {p1', p2'}.
// Assumes neither of the lines are degenerate.
func LineTransform(x1, y1, x2, y2, x1p, y1p, x2p, y2p float64) *Aff3 {
	// Calculate the offset, the rotation and the scale
	ox, oy := x1p-x1, y1p-y1
	dx, dy, dxp, dyp := x2-x1, y2-y1, x2p-x1p, y2p-y1p
	th := math.Atan2(dyp, dxp) - math.Atan2(dy, dx)
	s := math.Hypot(dxp, dyp) / math.Hypot(dx, dy)
	// Reverse order
	xfm := RotateAbout(th, x1p, y1p)
	xfm.ScaleAbout(s, s, x1p, y1p)
	xfm.Translate(ox, oy)
	return xfm
}

// BoxTransform produces an affine transform that maps the line {p1, p2} to {p1', p2'} and
// scales the perpendicular by hp / h. Assumes neither of the lines nor h are degenerate.
func BoxTransform(x1, y1, x2, y2, h, x1p, y1p, x2p, y2p, hp float64) *Aff3 {
	// Calculate the offset, the rotation and the scale
	ox, oy := x1p-x1, y1p-y1
	dx, dy, dxp, dyp := x2-x1, y2-y1, x2p-x1p, y2p-y1p
	th := math.Atan2(dyp, dxp) - math.Atan2(dy, dx)
	s := math.Hypot(dxp, dyp) / math.Hypot(dx, dy)
	// Reverse order
	xfm := RotateAbout(th, x1p, y1p)
	xfm.ScaleAbout(s, hp/h, x1p, y1p)
	xfm.Translate(ox, oy)
	return xfm
}

// BBTransform produces an affine transform that maps bounding box bb1 to bb2.
func BBTransform(bb1, bb2 [][]float64) *Aff3 {
	dx, dy := bb2[0][0]-bb1[0][0], bb2[0][1]-bb1[0][1]
	// Reverse order
	xfm := Translate(dx, dy)
	dx1, dy1, dx2, dy2 := bb1[1][0]-bb1[0][0], bb1[1][1]-bb1[0][1], bb2[1][0]-bb2[0][0], bb2[1][1]-bb2[0][1]
	sx, sy := dx2/dx1, dy2/dy1
	xfm.ScaleAbout(sx, sy, bb1[0][0], bb1[0][1])
	return xfm
}

// CreateAffineTransform returns an affine transform that performs the requested translation,
// scaling and rotation based on {0, 0}.
func CreateAffineTransform(x, y, scale, rotation float64) *Aff3 {
	xfm := Translate(x, y)
	xfm.Scale(scale, scale)
	xfm.Rotate(rotation)
	return xfm
}

// ScaleAndInset produces an affine transform that will scale and translate a set of points bounded by bb so they
// fit inside the inset box described by width, height, iwidth, iheight located at {0, 0}.
// If fix is true then the aspect ratio of bb is maintained.
func ScaleAndInset(width, height, iwidth, iheight float64, fix bool, bb [][]float64) *Aff3 {
	ox, oy := bb[0][0], bb[0][1]
	dx, dy := bb[1][0]-ox, bb[1][1]-oy

	w := width - 2*iwidth
	h := height - 2*iheight

	// Reverse order
	xfm := Translate(width/2, height/2)
	if fix {
		s := dx
		if dy > s {
			s = dy
		}
		xfm.Scale(w/s, h/s)
	} else {
		xfm.Scale(w/dx, h/dy)
	}
	xfm.Translate(-(ox + dx/2), -(oy + dy/2))

	return xfm
}

// FlipY is a convenience function to create an affine transform that has +ve Y point up rather than down.
func FlipY(height float64) *Aff3 {
	yoffs := height / 2
	// Reverse order
	xfm := Translate(0, yoffs)
	xfm.Scale(1, -1)
	xfm.Translate(0, -yoffs)
	return xfm
}

// Apply implements the Transform interface.
func (a *Aff3) Apply(pts ...[]float64) [][]float64 {
	npts := make([][]float64, len(pts))
	for i, pt := range pts {
		d := len(pt)
		npt := make([]float64, d)
		x, y := pt[0], pt[1]
		// x' = a[3*0+0]*x + a[3*0+1]*y + a[3*0+2]
		// y' = a[3*1+0]*x + a[3*1+1]*y + a[3*1+2]
		npt[0] = a[0]*x + a[1]*y + a[2]
		npt[1] = a[3]*x + a[4]*y + a[5]
		// Preserve other values
		for i := 2; i < d; i++ {
			npt[i] = pt[i]
		}
		npts[i] = npt
	}
	return npts
}

// Process implements the PathProcessor interface.
func (a *Aff3) Process(p *Path) []*Path {
	psteps := p.Steps()
	steps := make([][][]float64, len(psteps))
	for i, step := range psteps {
		steps[i] = a.Apply(step...)
	}

	path := NewPath(steps[0][0])
	for _, step := range steps {
		path.AddStep(step...)
	}

	if p.Closed() {
		path.Close()
	}

	return []*Path{path}
}
