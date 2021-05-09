package image

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

// Ordered Dithers

// B2Dither returns the Bayer ordered dither matrix 2x2
func B2Dither() [][]float64 {
	return [][]float64{
		{1 / 4.0, 3 / 4.0},
		{4 / 4.0, 2 / 4.0},
	}
}

// B4Dither returns the Bayer ordered dither matrix 4x4
func B4Dither() [][]float64 {
	return [][]float64{
		{1 / 16.0, 9 / 16.0, 3 / 16.0, 11 / 16.0},
		{13 / 16.0, 5 / 16.0, 15 / 16.0, 7 / 16.0},
		{4 / 16.0, 12 / 16.0, 2 / 16.0, 10 / 16.0},
		{16 / 16.0, 8 / 16.0, 14 / 16.0, 6 / 16.0},
	}
}

// B8Dither returns the Bayer ordered dither matrix 8x8
func B8Dither() [][]float64 {
	return [][]float64{
		{0.015625, 0.515625, 0.140625, 0.640625, 0.046875, 0.546875, 0.171875, 0.671875},
		{0.765625, 0.265625, 0.890625, 0.390625, 0.796875, 0.296875, 0.921875, 0.421875},
		{0.203125, 0.703125, 0.078125, 0.578125, 0.234375, 0.734375, 0.109375, 0.609375},
		{0.953125, 0.453125, 0.828125, 0.328125, 0.984375, 0.484375, 0.859375, 0.359375},
		{0.062500, 0.562500, 0.187500, 0.687500, 0.031250, 0.531250, 0.156250, 0.656250},
		{0.812500, 0.312500, 0.937500, 0.437500, 0.781250, 0.281250, 0.906250, 0.406250},
		{0.250000, 0.750000, 0.125000, 0.625000, 0.218750, 0.718750, 0.093750, 0.593750},
		{1.000000, 0.500000, 0.875000, 0.375000, 0.968750, 0.468750, 0.843750, 0.343750},
	}
}

// OrderedDither returns a ordered dithered image.
func OrderedDither(dst draw.Image, r image.Rectangle, c1, c2 color.Color, t float64, f func() [][]float64) {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	mat := f()
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

// Error Diffusion Dithers

// SimplestError returns a simpler error diffusion matrix 3x3
func SimplestError() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 1 / 2.0},
		{0, 1 / 2.0, 0},
	}, 1 / 4.0
}

// Simple3Error returns a simple error diffusion matrix 3x3
func Simple3Error() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 1 / 3.0},
		{0, 1 / 3.0, 1 / 3.0},
	}, 1 / 6.0
}

// FSError returns the Floyd Steinberg error diffusion matrix 3x3
func FSError() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 7 / 16.0},
		{3 / 16.0, 5 / 16.0, 1 / 16.0},
	}, 1 / 32.0
}

// JJNError returns the Jarvis Judice Ninke error diffusion matrix 5x5
func JJNError() ([][]float64, float64) {
	return [][]float64{
		// First two rows are empty
		{0, 0, 0, 7 / 48.0, 5 / 48.0},
		{3 / 48.0, 5 / 48.0, 7 / 48.0, 5 / 48.0, 3 / 48.0},
		{1 / 48.0, 3 / 48.0, 5 / 48.0, 3 / 48.0, 1 / 48.0},
	}, 1 / 96.0
}

// StukiError returns the Stuki error diffusion matrix 5x5
func StukiError() ([][]float64, float64) {
	return [][]float64{
		// First two rows are empty
		{0, 0, 0, 8 / 42.0, 4 / 42.0},
		{2 / 42.0, 4 / 42.0, 8 / 42.0, 4 / 42.0, 2 / 42.0},
		{1 / 42.0, 2 / 42.0, 4 / 42.0, 2 / 42.0, 1 / 42.0},
	}, 1 / 84.0
}

// BurkesError returns the Burkes error diffusion matrix 5x3
// (Stuki without the last row)
func BurkesError() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 0, 8 / 32.0, 4 / 32.0},
		{2 / 32.0, 4 / 32.0, 8 / 32.0, 4 / 32.0, 2 / 32.0},
	}, 1 / 64.0
}

// AtkinsonError returns the Atkinson error diffusion matrix 5x5
func AtkinsonError() ([][]float64, float64) {
	// 8 and not 6!
	return [][]float64{
		// First two rows are empty
		{0, 0, 0, 1 / 8.0, 1 / 8.0},
		{0, 1 / 8.0, 1 / 8.0, 1 / 8.0, 0},
		{0, 0, 1 / 8.0, 0, 0},
	}, 1 / 16.0
}

// MAYBE Modified Atkinson over 6?

// SierraError returns the Sierra error diffusion matrix 5x5
func SierraError() ([][]float64, float64) {
	return [][]float64{
		// First two rows are empty
		{0, 0, 0, 5 / 32.0, 3 / 32.0},
		{2 / 32.0, 4 / 32.0, 5 / 32.0, 4 / 32.0, 2 / 32.0},
		{0, 2 / 32.0, 3 / 32.0, 2 / 32.0, 0},
	}, 1 / 64.0
}

// Sierra2Error returns the Sierra two row error diffusion matrix 5x3
func Sierra2Error() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 0, 4 / 16.0, 3 / 16.0},
		{1 / 16.0, 2 / 16.0, 3 / 16.0, 2 / 16.0, 1 / 16.0},
	}, 1 / 32.0
}

// SierraLError returns the Sierra lite error diffusion matrix 3x3
func SierraLError() ([][]float64, float64) {
	return [][]float64{
		// First row is empty
		{0, 0, 2 / 4.0},
		{1 / 4.0, 1 / 4.0, 0},
	}, 1 / 8.0
}

func ErrorDither(dst draw.Image, r image.Rectangle, c1, c2 color.Color, t float64,
	ef func() ([][]float64, float64), noise bool) {
	mat, stoc := ef()
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
		_ = runErrorRow(rows, mat, t, stoc, noise)
	}

	// Generate image
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := runErrorRow(rows, mat, t, stoc, noise)
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

func runErrorRow(rows, mat [][]float64, t, stoc float64, noise bool) []float64 {
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
