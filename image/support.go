package image

import "image"

var (
	// Various supports

	// Z0 the point itself
	Z0 = [][]bool{
		{true},
	}

	// Z4 3x3 Von Neumann 4-way
	Z4 = [][]bool{
		{false, true, false},
		{true, true, true},
		{false, true, false},
	}

	// X3 3x3 X
	X3 = [][]bool{
		{true, false, true},
		{false, true, false},
		{true, false, true},
	}

	// Z8 3x3 Moore 8-way
	Z8 = [][]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	// Cross5x5 5x5 cross
	Cross5x5 = [][]bool{
		{false, false, true, false, false},
		{false, false, true, false, false},
		{true, true, true, true, true},
		{false, false, true, false, false},
		{false, false, true, false, false},
	}

	// X5 5x5 X
	X5 = [][]bool{
		{true, false, false, false, true},
		{false, true, false, true, false},
		{false, false, true, false, false},
		{false, true, false, true, false},
		{true, false, false, false, true},
	}

	// Star5x5 5x5 star
	Star5x5 = [][]bool{
		{true, false, true, false, true},
		{false, true, true, true, false},
		{true, true, true, true, true},
		{false, true, true, true, false},
		{true, false, true, false, true},
	}

	// Diamond5x5 5x5 diamond
	Diamond5x5 = [][]bool{
		{false, false, true, false, false},
		{false, true, true, true, false},
		{true, true, true, true, true},
		{false, true, true, true, false},
		{false, false, true, false, false},
	}

	// Ball5x5 5x5 ball
	Ball5x5 = [][]bool{
		{false, true, true, true, false},
		{true, true, true, true, true},
		{true, true, true, true, true},
		{true, true, true, true, true},
		{false, true, true, true, false},
	}

	// All5x5 5x5 block
	All5x5 = [][]bool{
		{true, true, true, true, true},
		{true, true, true, true, true},
		{true, true, true, true, true},
		{true, true, true, true, true},
		{true, true, true, true, true},
	}

	// HitOrMiss support pairs {C1, D1} and {C2, D2} for thinning, and their rotations.

	C11 = [][]bool{
		{false, false, false},
		{false, true, false},
		{true, true, true},
	}
	D11 = [][]bool{
		{true, true, true},
		{false, false, false},
		{false, false, false},
	}

	C12 = [][]bool{
		{true, false, false},
		{true, true, false},
		{true, false, false},
	}
	D12 = [][]bool{
		{false, false, true},
		{false, false, true},
		{false, false, true},
	}

	C13 = [][]bool{
		{true, true, true},
		{false, true, false},
		{false, false, false},
	}
	D13 = [][]bool{
		{false, false, false},
		{false, false, false},
		{true, true, true},
	}

	C14 = [][]bool{
		{false, false, true},
		{false, true, true},
		{false, false, true},
	}
	D14 = [][]bool{
		{true, false, false},
		{true, false, false},
		{true, false, false},
	}

	C21 = [][]bool{
		{false, false, false},
		{true, true, false},
		{true, true, false},
	}
	D21 = [][]bool{
		{false, true, true},
		{false, false, true},
		{false, false, false},
	}

	C22 = [][]bool{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}
	D22 = [][]bool{
		{false, false, false},
		{false, false, true},
		{false, true, true},
	}

	C23 = [][]bool{
		{false, true, true},
		{false, true, true},
		{false, false, false},
	}
	D23 = [][]bool{
		{false, false, false},
		{true, false, false},
		{true, true, false},
	}

	C24 = [][]bool{
		{false, false, false},
		{false, true, true},
		{false, true, true},
	}
	D24 = [][]bool{
		{true, true, false},
		{true, false, false},
		{false, false, false},
	}
)

// SupportToGray converts a support to a gray scale image.
func SupportToGray(suppt [][]bool) *image.Gray {
	w, h := len(suppt[0]), len(suppt)
	res := image.NewGray(image.Rect(0, 0, w, h))
	k := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if suppt[i][j] {
				res.Pix[k] = 0xff
			} else {
				res.Pix[k] = 0
			}
			k++
		}
	}
	return res
}

// GrayToSupport converts a gray scale image to a support.
func GrayToSupport(img *image.Gray) [][]bool {
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := make([][]bool, h)
	for i := 0; i < h; i++ {
		res[i] = make([]bool, w)
		so := img.PixOffset(0, i)
		for j := 0; j < w; j++ {
			res[i][j] = img.Pix[so+j] >= 128
		}
	}
	return res
}
