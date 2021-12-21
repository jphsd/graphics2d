package image

import (
	"image"
	"image/draw"
)

// And returns the result of ANDing img1 with img2, offset by offs (intersection). The images are converted to
// image.Gray if not already so. For a pair of pixels, p1 & p2 returns min(p1, p2).
func And(img1, img2 image.Image, offs image.Point) *image.Gray {
	r := img1.Bounds()
	res := image.NewGray(r)

	rect := r.Intersect(img2.Bounds().Add(r.Min.Sub(offs)))
	if rect.Empty() {
		// Implicit: AND with 0 => 0
		return res
	}

	process(img1, img2, res, offs, min)

	return res
}

// Or returns the result of ORing img1 with img2, offset by offs (union). The images are converted to
// image.Gray if not already so. For a pair of pixels, p1 & p2 returns max(p1, p2).
func Or(img1, img2 image.Image, offs image.Point) *image.Gray {
	r := img1.Bounds()
	res := image.NewGray(r)

	rect := r.Intersect(img2.Bounds().Add(r.Min.Sub(offs)))
	if rect.Empty() {
		// Implicit: OR with 0 => img1
		draw.Draw(res, r, img1, r.Min, draw.Src)
		return res
	}

	process(img1, img2, res, offs, max)

	return res
}

// Xor returns of XORing img1 wih img2, offset by offs (union - intersection). The images are converted to
// image.Gray if not already so. For a pair of pixels, p1 ^ p2 returns min(max(p1, p2), 1-min(p1, p2)).
func Xor(img1, img2 image.Image, offs image.Point) *image.Gray {
	r := img1.Bounds()
	res := image.NewGray(r)

	process(img1, img2, res, offs, max)

	return res
}

// Sub returns the result of subtracting img2 from img1, offset by offs. The images are converted to
// image.Gray if not already so. For a pair of pixels, p1 & p2 returns max(0, p1-p2).
func Sub(img1, img2 image.Image, offs image.Point) *image.Gray {
	r := img1.Bounds()
	res := image.NewGray(r)

	rect := r.Intersect(img2.Bounds().Add(r.Min.Sub(offs)))
	if rect.Empty() {
		// Implicit: Sub with 0 => img1
		draw.Draw(res, r, img1, r.Min, draw.Src)
		return res
	}

	process(img1, img2, res, offs, sub)

	return res
}

// utility function to perform grayscale conversion and apply f.
func process(img1, img2 image.Image, res *image.Gray, offs image.Point, f func(uint8, uint8) uint8) {
	r := img1.Bounds()

	// Convert to grayscale if necessary
	gray1, ok := img1.(*image.Gray)
	if !ok {
		gray1 = image.NewGray(r)
		draw.Draw(gray1, r, img1, r.Min, draw.Src)
	}

	// Convert to grayscale if necessary and restrict to r after offset
	gray2 := image.NewGray(r)
	draw.Draw(gray2, r, img2, offs, draw.Src)

	// gray2 and res have the same bounds and strides, gray1 may or may not
	for y := r.Min.Y; y < r.Max.Y; y++ {
		g1offs := gray1.PixOffset(r.Min.X, y)
		roffs := res.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			a := gray1.Pix[g1offs]
			b := gray2.Pix[roffs]
			v := f(a, b)
			res.Pix[roffs] = v
			g1offs++
			roffs++
		}
	}
}

// Not returns 1 - img.
func Not(img image.Image) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)

	// Convert to grayscale if necessary
	gray, ok := img.(*image.Gray)
	if !ok {
		gray = image.NewGray(r)
		draw.Draw(gray, r, img, r.Min, draw.Src)
	}

	for y := r.Min.Y; y < r.Max.Y; y++ {
		goffs := gray.PixOffset(r.Min.X, y)
		roffs := res.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			res.Pix[roffs] = 0xff - gray.Pix[goffs]
			goffs++
			roffs++
		}
	}

	return res
}

// Equal returns true if img1 within img2 at offset offs, matches.
func Equal(img1, img2 image.Image, offs image.Point) bool {
	r := img1.Bounds()

	// Convert to grayscale if necessary
	gray1, ok := img1.(*image.Gray)
	if !ok {
		gray1 = image.NewGray(r)
		draw.Draw(gray1, r, img1, r.Min, draw.Src)
	}

	// Convert to grayscale if necessary and restrict to r after offset
	gray2 := image.NewGray(r)
	draw.Draw(gray2, r, img2, offs, draw.Src)

	// Compare
	for y := r.Min.Y; y < r.Max.Y; y++ {
		g1offs := gray1.PixOffset(r.Min.X, y)
		g2offs := gray2.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			if gray1.Pix[g1offs] != gray2.Pix[g2offs] {
				return false
			}
			g1offs++
			g2offs++
		}
	}
	return true
}

// Copy creates a grayscale copy of img
func Copy(img image.Image) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)
	draw.Draw(res, r, img, r.Min, draw.Src)
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
func AlphaAnd(img1, img2 *image.Alpha, offs image.Point) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(And(g1, g2, offs))
}

// AlphaOr returns img1 | img2 (union).
// For a pair of pixels, p1 | p2 returns max(p1, p2).
func AlphaOr(img1, img2 *image.Alpha, offs image.Point) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Or(g1, g2, offs))
}

// AlphaXor returns img1 ^ img2 (union - intersection).
// For a pair of pixels, p1 ^ p2 returns min(max(p1, p2), 1-min(p1, p2)).
func AlphaXor(img1, img2 *image.Alpha, offs image.Point) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Xor(g1, g2, offs))
}

// AlphaSub returns img1 - img2 (img1 minus the intersection of img1 and img2).
// For a pair of pixels, p1 - p2 returns min(p1, 1-p2).
func AlphaSub(img1, img2 *image.Alpha, offs image.Point) *image.Alpha {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return GrayToAlpha(Sub(g1, g2, offs))
}

// AlphaNot returns 1 - img.
func AlphaNot(img *image.Alpha) *image.Alpha {
	g1 := AlphaToGray(img)
	return GrayToAlpha(Not(g1))
}

// AlphaEqual returns true if two images are the same.
func AlphaEqual(img1, img2 *image.Alpha, offs image.Point) bool {
	g1, g2 := AlphaToGray(img1), AlphaToGray(img2)
	return Equal(g1, g2, offs)
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
	}
	return a
}

func max(a, b uint8) uint8 {
	if b > a {
		return b
	}
	return a
}

func sub(a, b uint8) uint8 {
	if b > a {
		return 0
	}
	return a - b
}

func xor(a, b uint8) uint8 {
	return min(max(a, b), sub(0xff, min(a, b)))
}
