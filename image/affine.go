package image

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// Aff5 is a 5x5 affine transformation matrix in row major order, where the
// bottom row is implicitly [0 0 0 0 1].
//
// m[5*r+c] is the element in the r'th row and c'th column.
type Aff5 [20]float64

// NewAff5 creates the identity transform.
func NewAff5() *Aff5 {
	var res Aff5
	res[5*0+0] = 1
	res[5*1+1] = 1
	res[5*2+2] = 1
	res[5*3+3] = 1
	return &res
}

// Identity returns true if the transform is the identity.
func (a *Aff5) Identity() bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 5; j++ {
			if i == j && !util.Equals(a[5*i+j], 1) {
				return false
			} else if !util.Equals(a[5*i+j], 0) {
				return false
			}
		}
	}
	return true
}

// Apply applies the transform to the set of supplied points.
func (a *Aff5) Apply(pts ...[]float64) [][]float64 {
	npts := make([][]float64, len(pts))
	for i, pt := range pts {
		npt := make([]float64, 4)
		r, g, b, aa := pt[0], pt[1], pt[2], pt[4]
		// x' = a[3*0+0]*x + a[3*0+1]*y + a[3*0+2]
		// y' = a[3*1+0]*x + a[3*1+1]*y + a[3*1+2]
		npt[0] = a[0]*r + a[1]*g + a[2]*b + a[3]*aa + a[4]
		npt[1] = a[5]*r + a[6]*g + a[7]*b + a[8]*aa + a[9]
		npt[2] = a[10]*r + a[11]*g + a[12]*b + a[13]*aa + a[14]
		npt[3] = a[15]*r + a[16]*g + a[17]*b + a[18]*aa + a[19]
		npts[i] = npt
	}
	return npts
}

// Transform transforms an r, g, b, a tuple into r', g', b', a' according to
// the values in a.
func (a *Aff5) Transform(tup []uint8) []uint8 {
	pt := make([]float64, 4)
	for i := 0; i < len(pt); i++ {
		pt[i] = float64(tup[i]) / 0xff
	}
	pt = a.Apply(pt)[0]
	res := make([]uint8, 4)
	for i := 0; i < len(res); i++ {
		// Clamp first
		if pt[i] < 0 {
			pt[i] = 0
		} else if pt[i] > 1 {
			pt[i] = 1
		}
		res[i] = uint8(pt[i] * 0xff)
	}
	return res
}

// Matrix values are taken from https://www.w3.org/TR/SVG11/filters.html#feColorMatrixTypeAttribute
// and https://www.w3.org/TR/filter-effects-1/

// Saturate returns a transform that will modify the saturation of an image by s (0, 1)
func Saturate(s float64) *Aff5 {
	if s < 0 {
		s = 0
	} else if s > 1 {
		s = 1
	}
	oms := 1 - s

	res := NewAff5()
	res[5*0+0] = 0.2126 + 0.7874*s
	res[5*0+1] = 0.7152 * oms
	res[5*0+2] = 0.0722 * oms
	res[5*1+0] = 0.2126 * oms
	res[5*1+1] = 0.7152 + 0.2848*s
	res[5*1+2] = 0.0722 * oms
	res[5*2+0] = 0.2126 * oms
	res[5*2+1] = 0.7152 * oms
	res[5*2+2] = 0.0722 + 0.9278*s
	return res
}

// HueRotate returns a transform that will rotate the hue selection by th radians.
func HueRotate(th float64) *Aff5 {
	for th < -math.Pi*2 {
		th += math.Pi * 2
	}
	for th > math.Pi*2 {
		th -= math.Pi * 2
	}

	cth, sth := math.Cos(th), math.Sin(th)
	res := NewAff5()
	res[5*0+0] = 0.2126 + 0.7874*cth - 0.2126*sth
	res[5*0+1] = 0.7152 - 0.7152*cth - 0.7152*sth
	res[5*0+2] = 0.0722 - 0.0722*cth + 0.9278*sth

	res[5*1+0] = 0.2126 - 0.2126*cth + 0.143*sth // JH
	res[5*1+1] = 0.7152 + 0.2848*cth + 0.140*sth // JH
	res[5*1+2] = 0.0722 - 0.0722*cth - 0.283*sth // JH

	res[5*2+0] = 0.2126 - 0.2126*cth - 0.7874*sth
	res[5*2+1] = 0.7152 - 0.7152*cth + 0.7152*sth
	res[5*2+2] = 0.0722 + 0.9278*cth + 0.0722*sth

	return res
}

// LuminanceToAlpha returns a transform that will write just the luminance to the alpha channel.
func LuminanceToAlpha() *Aff5 {
	var res Aff5

	res[5*3+0] = 0.2126
	res[5*3+1] = 0.7152
	res[5*3+2] = 0.0722

	return &res
}

// Sepia returns a transform that will create a sepia tinted image.
func Sepia(t float64) *Aff5 {
	omt := 1 - t
	res := NewAff5()
	res[5*0+0] = 0.393 + 0.607*omt
	res[5*0+1] = 0.769 * t
	res[5*0+2] = 0.189 * t

	res[5*1+0] = 0.349 * t
	res[5*1+1] = 0.686 + 0.314*omt
	res[5*1+2] = 0.168 * t

	res[5*2+0] = 0.272 * t
	res[5*2+1] = 0.534 * t
	res[5*2+2] = 0.131 + 0.869*omt

	return res
}

// Grayscale returns a transform that will create a gray tinted image.
func Grayscale(t float64) *Aff5 {
	omt := 1 - t
	res := NewAff5()
	res[5*0+0] = 0.2126 + 0.7874*omt
	res[5*0+1] = 0.7152 * t
	res[5*0+2] = 0.0722 * t

	res[5*1+0] = 0.2126 * t
	res[5*1+1] = 0.7152 + 0.2848*omt
	res[5*1+2] = 0.0722 * t

	res[5*2+0] = 0.2126 * t
	res[5*2+1] = 0.7152 * t
	res[5*2+2] = 0.0722 + 0.9278*omt

	return res
}
