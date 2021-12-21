package texture

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

// Ordered Dithers
var (
	// B2Dither is the Bayer ordered dither matrix 2x2
	B2Dither = [][]float64{
		{1 / 4.0, 3 / 4.0},
		{4 / 4.0, 2 / 4.0},
	}

	// B4Dither is the Bayer ordered dither matrix 4x4
	B4Dither = [][]float64{
		{1 / 16.0, 9 / 16.0, 3 / 16.0, 11 / 16.0},
		{13 / 16.0, 5 / 16.0, 15 / 16.0, 7 / 16.0},
		{4 / 16.0, 12 / 16.0, 2 / 16.0, 10 / 16.0},
		{16 / 16.0, 8 / 16.0, 14 / 16.0, 6 / 16.0},
	}

	// B8Dither is the Bayer ordered dither matrix 8x8
	B8Dither = [][]float64{
		{0.015625, 0.515625, 0.140625, 0.640625, 0.046875, 0.546875, 0.171875, 0.671875},
		{0.765625, 0.265625, 0.890625, 0.390625, 0.796875, 0.296875, 0.921875, 0.421875},
		{0.203125, 0.703125, 0.078125, 0.578125, 0.234375, 0.734375, 0.109375, 0.609375},
		{0.953125, 0.453125, 0.828125, 0.328125, 0.984375, 0.484375, 0.859375, 0.359375},
		{0.062500, 0.562500, 0.187500, 0.687500, 0.031250, 0.531250, 0.156250, 0.656250},
		{0.812500, 0.312500, 0.937500, 0.437500, 0.781250, 0.281250, 0.906250, 0.406250},
		{0.250000, 0.750000, 0.125000, 0.625000, 0.218750, 0.718750, 0.093750, 0.593750},
		{1.000000, 0.500000, 0.875000, 0.375000, 0.968750, 0.468750, 0.843750, 0.343750},
	}
)

// OrderedDither returns an ordered dithered image.
func OrderedDither(dst draw.Image, c1, c2 color.Color, t float64, mat [][]float64) {
	r := dst.Bounds()

	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	nr, nc := len(mat), len(mat[0])
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if t < mat[y%nr][x%nc] {
				dst.Set(x, y, c1)
			} else {
				dst.Set(x, y, c2)
			}
		}
	}
}

// ToOrderedDither takes the grayscale version of img and renders it as a black and white dithered image
// using the specified dither mask.
func ToOrderedDither(img image.Image, mat [][]float64) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)

	gray, ok := img.(*image.Gray)
	if !ok {
		gray = image.NewGray(r)
		draw.Draw(gray, r, img, r.Min, draw.Src)
	}

	nr, nc := len(mat), len(mat[0])
	scale := 1 / float64(0xff)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		goffs := gray.PixOffset(r.Min.X, y)
		roffs := res.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			val := float64(gray.Pix[goffs]) * scale
			if val < mat[y%nr][x%nc] {
				res.Pix[roffs] = 0x00 // Black
			} else {
				res.Pix[roffs] = 0xff // White
			}
			goffs++
			roffs++
		}
	}

	return res
}

// Error Diffusion Dithers

// ErrorDiffusion holds the error diffusion matrix and a stochastic value to use when noise is required.
type ErrorDiffusion struct {
	Mat  [][]float64
	Stoc float64
}

// SimplestError returns a simpler error diffusion matrix 3x3
func SimplestError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 1 / 2.0},
		{0, 1 / 2.0, 0},
	}, 1 / 4.0}
}

// Simple3Error returns a simple error diffusion matrix 3x3
func Simple3Error() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 1 / 3.0},
		{0, 1 / 3.0, 1 / 3.0},
	}, 1 / 6.0}
}

// FSError returns the Floyd Steinberg error diffusion matrix 3x3
func FSError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 7 / 16.0},
		{3 / 16.0, 5 / 16.0, 1 / 16.0},
	}, 1 / 32.0}
}

// JJNError returns the Jarvis Judice Ninke error diffusion matrix 5x5
func JJNError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First two rows are empty
		{0, 0, 0, 7 / 48.0, 5 / 48.0},
		{3 / 48.0, 5 / 48.0, 7 / 48.0, 5 / 48.0, 3 / 48.0},
		{1 / 48.0, 3 / 48.0, 5 / 48.0, 3 / 48.0, 1 / 48.0},
	}, 1 / 96.0}
}

// StukiError returns the Stuki error diffusion matrix 5x5
func StukiError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First two rows are empty
		{0, 0, 0, 8 / 42.0, 4 / 42.0},
		{2 / 42.0, 4 / 42.0, 8 / 42.0, 4 / 42.0, 2 / 42.0},
		{1 / 42.0, 2 / 42.0, 4 / 42.0, 2 / 42.0, 1 / 42.0},
	}, 1 / 84.0}
}

// BurkesError returns the Burkes error diffusion matrix 5x3
// (Stuki without the last row)
func BurkesError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 0, 8 / 32.0, 4 / 32.0},
		{2 / 32.0, 4 / 32.0, 8 / 32.0, 4 / 32.0, 2 / 32.0},
	}, 1 / 64.0}
}

// AtkinsonError returns the Atkinson error diffusion matrix 5x5
func AtkinsonError() *ErrorDiffusion {
	// 8 and not 6!
	return &ErrorDiffusion{[][]float64{
		// First two rows are empty
		{0, 0, 0, 1 / 8.0, 1 / 8.0},
		{0, 1 / 8.0, 1 / 8.0, 1 / 8.0, 0},
		{0, 0, 1 / 8.0, 0, 0},
	}, 1 / 16.0}
}

// MAYBE Modified Atkinson over 6?

// SierraError returns the Sierra error diffusion matrix 5x5
func SierraError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First two rows are empty
		{0, 0, 0, 5 / 32.0, 3 / 32.0},
		{2 / 32.0, 4 / 32.0, 5 / 32.0, 4 / 32.0, 2 / 32.0},
		{0, 2 / 32.0, 3 / 32.0, 2 / 32.0, 0},
	}, 1 / 64.0}
}

// Sierra2Error returns the Sierra two row error diffusion matrix 5x3
func Sierra2Error() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 0, 4 / 16.0, 3 / 16.0},
		{1 / 16.0, 2 / 16.0, 3 / 16.0, 2 / 16.0, 1 / 16.0},
	}, 1 / 32.0}
}

// SierraLError returns the Sierra lite error diffusion matrix 3x3
func SierraLError() *ErrorDiffusion {
	return &ErrorDiffusion{[][]float64{
		// First row is empty
		{0, 0, 2 / 4.0},
		{1 / 4.0, 1 / 4.0, 0},
	}, 1 / 8.0}
}

// ErrorDither returns an error diffused dithered image using the error matrix. If noise is true, the diffusion
// weights are randomly +1/-1/0 times the stochastic value.
func (ed *ErrorDiffusion) ErrorDither(dst draw.Image, c1, c2 color.Color, t float64, noise bool) {
	r := dst.Bounds()

	mat := ed.Mat
	nr, nc := len(mat), len(mat[0])
	n := r.Dx() + 2*nc // allow for diffusion at edges
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	rows := make([][]float64, nr)
	for i := 0; i < nr; i++ {
		rows[i] = loadRow(t, n)
	}

	// Warm it up
	for i := 0; i < nr; i++ {
		_ = ed.runErrorRow(rows, t, noise)
	}

	// Generate image
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := ed.runErrorRow(rows, t, noise)
		i := nc
		for x := r.Min.X; x < r.Max.X; x++ {
			if row[i] < 0.5 {
				dst.Set(x, y, c1)
			} else {
				dst.Set(x, y, c2)
			}
			i++
		}
	}
}

func (ed *ErrorDiffusion) runErrorRow(rows [][]float64, t float64, noise bool) []float64 {
	mat, stoc := ed.Mat, ed.Stoc
	nr, nc := len(mat), len(mat[0])
	e := nc / 2 // number of cells to L or R of midpoint
	n := len(rows[0])

	for c := e; c < n-e; c++ {
		pv := rows[0][c]
		pnv := closest(pv)
		rows[0][c] = pnv
		err := pv - pnv
		// Spread the error
		for i := e + 1; i < nc; i++ {
			if noise {
				rows[0][c+i-e] += err * (mat[0][i] + float64(rand.Intn(3)-1)*stoc)
			} else {
				rows[0][c+i-e] += err * mat[0][i]
			}
		}
		for r := 1; r < nr; r++ {
			for i := 0; i < nc; i++ {
				if noise {
					rows[r][c+i-e] += err * (mat[r][i] + float64(rand.Intn(3)-1)*stoc)
				} else {
					rows[r][c+i-e] += err * mat[r][i]
				}
			}
		}
	}

	// Capture result and update rows
	res := rows[0]
	for i := 1; i < nr; i++ {
		rows[i-1] = rows[i]
	}
	rows[nr-1] = loadRow(t, n)
	return res
}

func loadRow(v float64, n int) []float64 {
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = v
	}
	return res
}

func closest(v float64) float64 {
	if v < 0.5 {
		return 0
	}
	return 1
}

// ToErrorDither takes the grayscale version of img and renders it as a black and white dithered image
// using the specified error dither.
func (ed *ErrorDiffusion) ToErrorDither(img image.Image) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)

	gray, ok := img.(*image.Gray)
	if !ok {
		gray = image.NewGray(r)
		draw.Draw(gray, r, img, r.Min, draw.Src)
	}

	mat := ed.Mat
	nr, nc := len(mat), len(mat[0])

	rows := make([][]float64, nr)
	for i := 0; i < nr; i++ {
		rows[i] = loadImgRow(gray, r.Min.Y+i, nc)
	}

	// Convert image
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := ed.runImgErrorRow(rows, gray, y, nc)
		i := nc
		roffs := res.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			if row[i] < 0.5 {
				res.Pix[roffs] = 0x0
			} else {
				res.Pix[roffs] = 0xff
			}
			roffs++
			i++
		}
	}

	return res
}

func (ed *ErrorDiffusion) runImgErrorRow(rows [][]float64, img *image.Gray, y, nc int) []float64 {
	mat := ed.Mat
	nr, nc := len(mat), len(mat[0])
	e := nc / 2 // number of cells to L or R of midpoint
	n := len(rows[0])

	for c := e; c < n-e; c++ {
		pv := rows[0][c]
		pnv := closest(pv)
		rows[0][c] = pnv
		err := pv - pnv
		// Spread the error
		for i := e + 1; i < nc; i++ {
			rows[0][c+i-e] += err * mat[0][i]
		}
		for r := 1; r < nr; r++ {
			for i := 0; i < nc; i++ {
				rows[r][c+i-e] += err * mat[r][i]
			}
		}
	}

	// Capture result and update rows
	// TODO - use copy()
	res := rows[0]
	for i := 1; i < nr; i++ {
		rows[i-1] = rows[i]
	}
	rows[nr-1] = loadImgRow(img, y+nr, nc)
	return res
}

func loadImgRow(img *image.Gray, y, nc int) []float64 {
	r := img.Bounds()
	dx := r.Dx()
	w := dx + 2*nc
	res := make([]float64, w)
	if y >= r.Max.Y {
		for i := 0; i < w; i++ {
			res[i] = 0
		}
		return res
	}

	goffs := img.PixOffset(r.Min.X, y)
	scale := 1 / float64(0xff)
	v := float64(img.Pix[goffs]) * scale
	for i := 0; i < nc; i++ {
		res[i] = v
	}
	for i := nc; i < dx+nc; i++ {
		v = float64(img.Pix[goffs]) * scale
		res[i] = v
		goffs++
	}
	for i := dx + nc; i < w; i++ {
		res[i] = v
	}
	return res
}
