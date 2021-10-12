package graphics2d

import (
	"fmt"
	"math"

	"github.com/jphsd/graphics2d/util"
)

// Aff3 is a 3x3 affine transformation matrix in row major order, where the
// bottom row is implicitly [0 0 1].
//
// m[3*r+c] is the element in the r'th row and c'th column.
type Aff3 [6]float64

// NewAff3 creates the identity transform.
func NewAff3() *Aff3 {
	var res Aff3
	res[3*0+0] = 1
	res[3*1+1] = 1
	return &res
}

// Determinant calculates the transform's matrix determinant.
func (a *Aff3) Determinant() float64 {
	return a[3*0+0]*a[3*1+1] - a[3*0+1]*a[3*1+0]
}

// Translate adds a translation to the transform.
func (a *Aff3) Translate(x, y float64) {
	a[3*0+2] = x*a[3*0+0] + y*a[3*0+1] + a[3*0+2]
	a[3*1+2] = x*a[3*1+0] + y*a[3*1+1] + a[3*1+2]
}

// Rotate adds a rotation to the transform. The rotation is about {0, 0}.
func (a *Aff3) Rotate(th float64) {
	sin, cos := math.Sin(th), math.Cos(th)
	m0, m1 := a[3*0+0], a[3*0+1]
	a[3*0+0] = cos*m0 + sin*m1
	a[3*0+1] = -sin*m0 + cos*m1
	m0, m1 = a[3*1+0], a[3*1+1]
	a[3*1+0] = cos*m0 + sin*m1
	a[3*1+1] = -sin*m0 + cos*m1
}

// RotateAbout adds a rotation about a point to the transform.
func (a *Aff3) RotateAbout(th, ax, ay float64) {
	// Reverse order
	a.Translate(ax, ay)
	a.Rotate(th)
	a.Translate(-ax, -ay)
}

// QuadrantRotate adds a rotation (n * 90 degrees) to the transform. The rotation is about {0, 0}.
// It avoids rounding issues with the trig functions.
func (a *Aff3) QuadrantRotate(n int) {
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
}

// QuadrantRotateAbout adds a rotation (n * 90 degrees) about a point to the transform.
// It avoids rounding issues with the trig functions.
func (a *Aff3) QuadrantRotateAbout(n int, ax, ay float64) {
	// Reverse order
	a.Translate(ax, ay)
	a.QuadrantRotate(n)
	a.Translate(-ax, -ay)
}

// Scale adds a scaling to the transform centered on {0, 0}.
func (a *Aff3) Scale(sx, sy float64) {
	a[3*0+0] *= sx
	a[3*1+1] *= sy
	a[3*0+1] *= sy
	a[3*1+0] *= sx
}

// ScaleAbout adds a scale about a point to the transform.
func (a *Aff3) ScaleAbout(sx, sy, ax, ay float64) {
	// Reverse order
	a.Translate(ax, ay)
	a.Scale(sx, sy)
	a.Translate(-ax, -ay)
}

// Shear adds a shear to the transform centered on {0, 0}.
func (a *Aff3) Shear(shx, shy float64) {
	m0, m1 := a[3*0+0], a[3*0+1]
	a[3*0+0] = m0 + m1*shy
	a[3*0+1] = m0*shx + m1

	m0, m1 = a[3*1+0], a[3*1+1]
	a[3*1+0] = m0 + m1*shy
	a[3*1+1] = m0*shx + m1
}

// ShearAbout adds a shear about a point to the transform.
func (a *Aff3) ShearAbout(shx, shy, ax, ay float64) {
	// Reverse order
	a.Translate(ax, ay)
	a.Shear(shx, shy)
	a.Translate(-ax, -ay)
}

// Concatenate concatenates a transform to the transform.
func (a *Aff3) Concatenate(aff Aff3) {
	m00, m01, m10, m11 := a[3*0+0], a[3*0+1], a[3*1+0], a[3*1+1]
	t00, t01, t02, t10, t11, t12 := aff[3*0+0], aff[3*0+1], aff[3*0+2], aff[3*1+0], aff[3*1+1], aff[3*1+2]

	a[3*0+0] = t00*m00 + t10*m01
	a[3*0+1] = t01*m00 + t11*m01
	a[3*0+2] += t02*m00 + t12*m01

	a[3*1+0] = t00*m10 + t10*m11
	a[3*1+1] = t01*m10 + t11*m11
	a[3*1+2] += t02*m10 + t12*m11
}

// PreConcatenate preconcatenates a transform to the transform.
func (a *Aff3) PreConcatenate(aff Aff3) {
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
}

// InverseOf returns the inverse of the transform.
func (a *Aff3) InverseOf() (*Aff3, error) {
	res := a.Copy()
	if err := res.Invert(); err != nil {
		return nil, err
	}

	return res, nil
}

// Invert inverts the transform.
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

// String converts the transform into a string.
func (a *Aff3) String() string {
	return fmt.Sprintf("{{%f, %f, %f}, {%f, %f, %f}, {0, 0, 1}}",
		a[3*0+0], a[3*0+1], a[3*0+2],
		a[3*1+0], a[3*1+1], a[3*1+2])
}

// Identity returns true if the transform is the identity.
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

// Copy returns a copy of the transform.
func (a *Aff3) Copy() *Aff3 {
	res := *a
	return &res
}

// Reflect performs a reflection along the axis defined by the two non-coincident points.
func (a *Aff3) Reflect(x1, y1, x2, y2 float64) {
	dx, dy := x2-x1, y2-y1
	if util.Equals(dy, 0) {
		// Horizontal - no rotation required
		a.Translate(0, y1)
		a.Scale(1, -1)
		a.Translate(0, -y1)
		return
	}
	if util.Equals(dx, 0) {
		// Vertical - no rotation required
		a.Translate(x1, 0)
		a.Scale(-1, 1)
		a.Translate(-x1, 0)
		return
	}
	th := math.Atan2(dy, dx)
	if th < 0 {
		th += math.Pi * 2
	}
	a.Translate(x1, y1)
	a.Rotate(th)
	a.Scale(1, -1)
	a.Rotate(-th)
	a.Translate(-x1, -y1)
}

// LineTransform produces a transform that maps the line {p1, p2} to {p1', p2'}.
// Assumes neither of the lines are degenerate.
func LineTransform(x1, y1, x2, y2, x1p, y1p, x2p, y2p float64) *Aff3 {
	// Calculate the offset, the rotation and the scale
	ox, oy := x1p-x1, y1p-y1
	dx, dy, dxp, dyp := x2-x1, y2-y1, x2p-x1p, y2p-y1p
	th := math.Atan2(dyp, dxp) - math.Atan2(dy, dx)
	s := math.Sqrt(dxp*dxp+dyp*dyp) / math.Sqrt(dx*dx+dy*dy)
	xfm := NewAff3()
	// Reverse order
	xfm.RotateAbout(th, x1p, y1p)
	xfm.ScaleAbout(s, s, x1p, y1p)
	xfm.Translate(ox, oy)
	return xfm
}

// LineTransform2 produces a transform that maps the line {p1, p2} to {p1', p2'} and
// scales the perpendicular by hp / h. Assumes neither of the lines nor h are degenerate.
func LineTransform2(x1, y1, x2, y2, h, x1p, y1p, x2p, y2p, hp float64) *Aff3 {
	// Calculate the offset, the rotation and the scale
	ox, oy := x1p-x1, y1p-y1
	dx, dy, dxp, dyp := x2-x1, y2-y1, x2p-x1p, y2p-y1p
	th := math.Atan2(dyp, dxp) - math.Atan2(dy, dx)
	s := math.Sqrt(dxp*dxp+dyp*dyp) / math.Sqrt(dx*dx+dy*dy)
	xfm := NewAff3()
	// Reverse order
	xfm.RotateAbout(th, x1p, y1p)
	xfm.ScaleAbout(s, hp/h, x1p, y1p)
	xfm.Translate(ox, oy)
	return xfm
}

// CreateTransform returns a transform that performs the requested translation,
// scaling and rotation based on {0, 0}.
func CreateTransform(x, y, scale, rotation float64) *Aff3 {
	xfm := NewAff3()
	xfm.Translate(x, y)
	xfm.Scale(scale, scale)
	xfm.Rotate(rotation)
	return xfm
}

// Apply applies the transform to the set of supplied points.
func (a *Aff3) Apply(pts ...[]float64) [][]float64 {
	npts := make([][]float64, len(pts))
	for i, pt := range pts {
		npt := make([]float64, 2)
		x, y := pt[0], pt[1]
		// x' = a[3*0+0]*x + a[3*0+1]*y + a[3*0+2]
		// y' = a[3*1+0]*x + a[3*1+1]*y + a[3*1+2]
		npt[0] = a[0]*x + a[1]*y + a[2]
		npt[1] = a[3]*x + a[4]*y + a[5]
		npts[i] = npt
	}
	return npts
}
