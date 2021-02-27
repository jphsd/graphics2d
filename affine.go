package graphics2d

import (
	"fmt"
	"github.com/jphsd/graphics2d/util"
	"math"
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

// QuadrantRotate adds a rotation (n * 90 degrees) about a point to the transform.
// It avoids rounding issues with the trig functions.
func (a *Aff3) QuadrantRotateAbout(n int, ax, ay float64) {
	// Reverse order
	a.Translate(ax, ay)
	a.QuadrantRotate(n)
	a.Translate(-ax, -ay)
}

// Scale adds a scaling to the transform.
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

// Shear adds a shear to the transform.
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
	err := res.Invert()
	return res, err
}

// Invert inverts the transform.
func (a *Aff3) Invert() error {
	det := a.Determinant()
	if util.Equalsf64(math.Abs(det), 0) {
		return fmt.Errorf("Determinant is zero => non-invertible")
	}

	m00, m01, m02, m10, m11, m12 := a[3*0+0], a[3*0+1], a[3*0+2], a[3*1+0], a[3*1+1], a[3*1+2]
	a[3*0+0], a[3*0+1], a[3*0+2] = m11/det, -m10/det, -m01/det
	a[3*1+0], a[3*1+1], a[3*1+2] = m00/det, (m01*m12-m11*m02)/det, (m10*m02-m00*m12)/det
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
	if !util.Equalsf64(a[3*0+0], 1) {
		return false
	}
	if !util.Equalsf64(a[3*0+1], 0) {
		return false
	}
	if !util.Equalsf64(a[3*0+2], 0) {
		return false
	}
	if !util.Equalsf64(a[3*1+0], 0) {
		return false
	}
	if !util.Equalsf64(a[3*1+1], 1) {
		return false
	}
	if !util.Equalsf64(a[3*1+2], 0) {
		return false
	}
	return true
}

// Copy returns a copy of the transform.
func (a *Aff3) Copy() *Aff3 {
	res := *a
	return &res
}
