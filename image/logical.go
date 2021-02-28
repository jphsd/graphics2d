package image

import (
	"fmt"
	"image"
)

// And returns img1 & img2 (intersection).
// For a pair of pixels, p1 & p2 returns min(p1, p2).
func And(img1, img2 *image.Gray) *image.Gray {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		panic(fmt.Errorf("image sizes don't match"))
	}
	res := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0, y)
		i2 := img2.PixOffset(0, y)
		d := res.Stride * y
		for x := 0; x < w; x, i1, i2, d = x+1, i1+1, i2+1, d+1 {
			res.Pix[d] = min(img1.Pix[i1], img2.Pix[i2])
		}
	}
	return res
}

// Or returns img1 | img2 (union).
// For a pair of pixels, p1 | p2 returns max(p1, p2).
func Or(img1, img2 *image.Gray) *image.Gray {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		panic(fmt.Errorf("image sizes don't match"))
	}
	res := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0, y)
		i2 := img2.PixOffset(0, y)
		d := res.Stride * y
		for x := 0; x < w; x, i1, i2, d = x+1, i1+1, i2+1, d+1 {
			res.Pix[d] = max(img1.Pix[i1], img2.Pix[i2])
		}
	}
	return res
}

// Xor returns img1 ^ img2 (union - intersection).
// For a pair of pixels, p1 ^ p2 returns min(max(p1, p2), 1-min(p1, p2)).
func Xor(img1, img2 *image.Gray) *image.Gray {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		panic(fmt.Errorf("image sizes don't match"))
	}
	res := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0, y)
		i2 := img2.PixOffset(0, y)
		d := res.Stride * y
		for x := 0; x < w; x, i1, i2, d = x+1, i1+1, i2+1, d+1 {
			p1 := img1.Pix[i1]
			p2 := img2.Pix[i2]
			res.Pix[d] = min(max(p1, p2), 0xff-min(p1, p2))
		}
	}
	return res
}

// Sub returns img1 - img2 (img1 - intersection of img1 and img2).
// For a pair of pixels, p1 - p2 returns min(p1, 1-p2).
func Sub(img1, img2 *image.Gray) *image.Gray {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		panic(fmt.Errorf("image sizes don't match"))
	}
	res := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0, y)
		i2 := img2.PixOffset(0, y)
		d := res.Stride * y
		for x := 0; x < w; x, i1, i2, d = x+1, i1+1, i2+1, d+1 {
			// img1 intersected with not(img2)
			res.Pix[d] = img1.Pix[i1] - min(img1.Pix[i1], img2.Pix[i2])
		}
	}
	return res
}

// Not returns 1 - img.
func Not(img *image.Gray) *image.Gray {
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		i1 := img.PixOffset(0, y)
		d := res.Stride * y
		for x := 0; x < w; x, i1, d = x+1, i1+1, d+1 {
			res.Pix[d] = 0xff - img.Pix[i1]
		}
	}
	return res
}

// Equals returns true if two images are the same.
func Equals(img1, img2 *image.Gray) bool {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		return false
	}
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0, y)
		i2 := img2.PixOffset(0, y)
		for x := 0; x < w; x, i1, i2 = x+1, i1+1, i2+1 {
			if img1.Pix[i1] != img2.Pix[i2] {
				return false
			}
		}
	}
	return true
}

// Copy creates a deep copy of img
func Copy(img *image.Gray) *image.Gray {
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewGray(image.Rect(0, 0, w, h))
	if img.Stride == res.Stride {
		// All at once
		copy(img.Pix, res.Pix)
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			start := img.PixOffset(0, i)
			os := res.Pix[0:w]
			is := img.Pix[start : start+w]
			copy(is, os)
		}
	}
	return res
}

// AlphaToGray does a shallow copy (vs going through the ColorModel).
func AlphaToGray(a *image.Alpha) *image.Gray {
	return &image.Gray{a.Pix, a.Stride, a.Rect}
}

// GrayToAlpha does a shallow copy (vs going through the ColorModel).
func GrayToAlpha(a *image.Gray) *image.Alpha {
	return &image.Alpha{a.Pix, a.Stride, a.Rect}
}

// AlphaAnd returns img1 & img2 (intersection).
// For a pair of pixels, p1 & p2 returns min(p1, p2).
func AlphaAnd(img1, img2 *image.Alpha) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(And(g1, g2))
}

// AlphaOr returns img1 | img2 (union).
// For a pair of pixels, p1 | p2 returns max(p1, p2).
func AlphaOr(img1, img2 *image.Alpha) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Or(g1, g2))
}

// AlphaXor returns img1 ^ img2 (union - intersection).
// For a pair of pixels, p1 ^ p2 returns min(max(p1, p2), 1-min(p1, p2)).
func AlphaXor(img1, img2 *image.Alpha) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Xor(g1, g2))
}

// AlphaSub returns img1 - img2 (img1 minus the intersection of img1 and img2).
// For a pair of pixels, p1 - p2 returns min(p1, 1-p2).
func AlphaSub(img1, img2 *image.Alpha) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Sub(g1, g2))
}

// AlphaNot returns 1 - img.
func AlphaNot(img *image.Alpha) *image.Alpha {
	g1 := AlphaToGray(img)
	return GrayToAlpha(Not(g1))
}

// AlphaEquals returns true if two images are the same.
func AlphaEquals(img1, img2 *image.Alpha) bool {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return Equals(g1, g2)
}

// AlphaCopy returns a deep copy of img.
func AlphaCopy(img *image.Alpha) *image.Alpha {
	g1 := AlphaToGray(img)
	return GrayToAlpha(Copy(g1))
}

// vs inline?
func min(a, b uint8) uint8 {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a, b uint8) uint8 {
	if b > a {
		return b
	} else {
		return a
	}
}
