package image

import (
	"fmt"
	"image"
	"image/color"
)

// RemapGray remaps a color map with a look up table. The lookup table must be 256 long.
func RemapGray(img *image.Gray, lut []uint8) *image.Gray {
	if len(lut) != 256 {
		panic(fmt.Errorf("lut must be 256 long"))
	}
	imgR := img.Bounds()

	w, h := img.Rect.Dx(), img.Rect.Dy()
	res := image.NewGray(image.Rect(0, 0, w, h))

	sp, dp := img.Pix, res.Pix
	if img.Stride == res.Stride {
		// All at once
		for i := 0; i < len(res.Pix); i++ {
			dp[i] = lut[sp[i]]
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so := img.PixOffset(0+imgR.Min.X, i+imgR.Min.Y)
			do := res.PixOffset(0, i)
			for j := 0; j < w; j++ {
				dp[do] = lut[sp[so]]
				so++
				do++
			}
		}
	}
	return res
}

// RemapGray2RGBA remaps a grayscale image to a RGBA one using a look up table. The lookup table must be 256 long.
func RemapGray2RGBA(img image.Image, lut []color.Color) *image.RGBA {
	if len(lut) != 256 {
		panic(fmt.Errorf("lut must be 256 long"))
	}

	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			gc := color.GrayModel.Convert(img.At(x+imgR.Min.X, y+imgR.Min.Y)).(color.Gray)
			res.Set(x, y, lut[gc.Y])
		}
	}

	return res
}

// RemapRGB remaps all the color maps with R, G, B look up tables. The lookup tables must all be 256 long.
// The maps are applied after alpha pre-multiplication.
func RemapRGB(img *image.RGBA, lutR, lutG, lutB []uint8) *image.RGBA {
	if len(lutR) != 256 || len(lutG) != 256 || len(lutB) != 256 {
		panic(fmt.Errorf("luts must be 256 long"))
	}
	imgR := img.Bounds()
	// Convert luts to alpha pre-multiplied
	lRA := PremulLut(lutR)
	lGA := PremulLut(lutG)
	lBA := PremulLut(lutB)

	w, h := img.Rect.Dx(), img.Rect.Dy()
	res := image.NewRGBA(image.Rect(0, 0, w, h))

	sp, dp := img.Pix, res.Pix
	if img.Stride == res.Stride {
		// All at once
		for i := 0; i < len(res.Pix); i += 4 {
			a := sp[i+3]
			dp[i] = lRA[sp[i]][a]
			dp[i+1] = lGA[sp[i+1]][a]
			dp[i+2] = lBA[sp[i+2]][a]
			dp[i+3] = a
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so := img.PixOffset(0+imgR.Min.X, i+imgR.Min.Y)
			do := res.PixOffset(0, i)
			for j := 0; j < w; j++ {
				sj, dj := so+j, do+j
				a := sp[sj+3]
				dp[dj] = lRA[sp[sj]][a]
				dp[dj+1] = lGA[sp[sj+1]][a]
				dp[dj+2] = lBA[sp[sj+2]][a]
				dp[dj+3] = a
			}
		}
	}
	return res
}

// RemapRGBSingle remaps all the color maps with a single look up table. The lookup table must be 256 long.
// The maps are applied after alpha pre-multiplication.
func RemapRGBSingle(img *image.RGBA, lut []uint8) *image.RGBA {
	if len(lut) != 256 {
		panic(fmt.Errorf("lut must be 256 long"))
	}
	imgR := img.Bounds()
	// Convert lut to alpha pre-multiplied
	lA := PremulLut(lut)

	w, h := img.Rect.Dx(), img.Rect.Dy()
	res := image.NewRGBA(image.Rect(0, 0, w, h))

	sp, dp := img.Pix, res.Pix
	if img.Stride == res.Stride {
		// All at once
		for i := 0; i < len(res.Pix); i += 4 {
			a := sp[i+3]
			dp[i] = lA[sp[i]][a]
			dp[i+1] = lA[sp[i+1]][a]
			dp[i+2] = lA[sp[i+2]][a]
			dp[i+3] = a
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so := img.PixOffset(0+imgR.Min.X, i+imgR.Min.Y)
			do := res.PixOffset(0, i)
			for j := 0; j < w; j++ {
				sj, dj := so+j, do+j
				a := sp[sj+3]
				dp[dj] = lA[sp[sj]][a]
				dp[dj+1] = lA[sp[sj+1]][a]
				dp[dj+2] = lA[sp[sj+2]][a]
				dp[dj+3] = a
			}
		}
	}
	return res
}
