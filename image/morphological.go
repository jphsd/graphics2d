package image

import (
	"image"
	"image/draw"
	"sort"
)

// MinOp finds the lowest value in values that has the support set.
func MinOp(values [][]uint8, support [][]bool) uint8 {
	min := uint8(0xff)
	for i := 0; i < len(support); i++ {
		sr := support[i]
		vr := values[i]
		for j := 0; j < len(sr); j++ {
			if sr[j] {
				v := vr[j]
				if v < min {
					min = v
				}
			}
		}
	}
	return min
}

// MaxOp finds the highest value in values that has the support set.
func MaxOp(values [][]uint8, support [][]bool) uint8 {
	max := uint8(0)
	for i := 0; i < len(support); i++ {
		sr := support[i]
		vr := values[i]
		for j := 0; j < len(sr); j++ {
			if sr[j] {
				v := vr[j]
				if v > max {
					max = v
				}
			}
		}
	}
	return max
}

// AvgOp finds the average value (rounded down) of all the values that have the support set.
func AvgOp(values [][]uint8, support [][]bool) uint8 {
	sum, n := 0, 0
	for i := 0; i < len(support); i++ {
		sr := support[i]
		vr := values[i]
		for j := 0; j < len(sr); j++ {
			if sr[j] {
				n++
				sum += int(vr[j])
			}
		}
	}
	sum /= n
	return uint8(sum)
}

// MedOp finds the midway value of all the values, sorted, that have the support set.
func MedOp(values [][]uint8, support [][]bool) uint8 {
	vals, n := []int{}, 0
	for i := 0; i < len(support); i++ {
		sr := support[i]
		vr := values[i]
		for j := 0; j < len(sr); j++ {
			if sr[j] {
				n++
				vals = append(vals, int(vr[j]))
			}
		}
	}
	sort.Ints(vals)
	return uint8(vals[n/2])
}

// Sorry, no ModeOp for those of you looking for it.

// Morphological runs op over the image img using support and supplying def when a location
// falls outside of the image boundary. The support dimensions must be odd (not checked for).
func Morphological(img image.Image, op func([][]uint8, [][]bool) uint8, support [][]bool, def uint8) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)

	// Convert to grayscale if necessary
	gray, ok := img.(*image.Gray)
	if !ok {
		gray := image.NewGray(r)
		draw.Draw(gray, r, img, r.Min, draw.Src)
	}

	sw, sh := len(support[0]), len(support)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		roffs := gray.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			values := morphHelper(x, y, sw, sh, gray, def)
			res.Pix[roffs] = op(values, support)
			roffs++
		}
	}
	return res
}

func morphHelper(x, y, sw, sh int, img *image.Gray, def uint8) [][]uint8 {
	r := img.Bounds()
	// if entire value set is within the image, create new slice references
	// else copy def and valid values
	res := make([][]uint8, sh)
	w1, h1 := sw/2, sh/2
	x1, y1, x2, y2 := x-w1, y-h1, x+w1, y+h1
	if x1 >= r.Min.X && y1 >= r.Min.Y && x2 < r.Max.X && y2 < r.Max.Y {
		// We can construct new slices using the image
		for i := 0; i < sh; i++ {
			so := img.PixOffset(x1, y1+i)
			res[i] = img.Pix[so : so+sw]
		}
		return res
	}

	da := make([]uint8, sw)
	for i := 0; i < sw; i++ {
		da[i] = def
	}
	for i := 0; i < sh; i++ {
		if y1+i < r.Min.Y || y1+i >= r.Max.Y {
			// y out of bounds
			res[i] = da
			continue
		}
		if x1 >= r.Min.X && x2 < r.Max.X {
			// x1, x2 and y in bounds
			so := img.PixOffset(x1, y1+i)
			res[i] = img.Pix[so : so+sw]
			continue
		}
		res[i] = make([]uint8, sw)
		for j := x1; j < x1+sw; j++ {
			if j < r.Min.X || j >= r.Max.X {
				res[i][j-x1] = def
			} else {
				so := img.PixOffset(j, y1+i)
				res[i][j-x1] = img.Pix[so]
			}
		}
	}
	return res
}

// Dilate replaces each pixel with the max of the pixels in its support.
func Dilate(img image.Image, support [][]bool) *image.Gray {
	return Morphological(img, MaxOp, support, 0)
}

// Erode replaces each pixel with the min of the pixels in its support.
func Erode(img image.Image, support [][]bool) *image.Gray {
	return Morphological(img, MinOp, support, 0xff)
}

// Open applies a dilation to a prior erosion.
func Open(img image.Image, support [][]bool) *image.Gray {
	return Dilate(Erode(img, support), support)
}

// Close applies an erosion to a prior dilation.
func Close(img image.Image, support [][]bool) *image.Gray {
	return Erode(Dilate(img, support), support)
}

// TopHat is the subtraction of an open from the original.
func TopHat(img image.Image, support [][]bool) *image.Gray {
	r := img.Bounds()
	return Sub(img, Open(img, support), r.Min)
}

// BotHat is the subtraction of the original from a close.
func BotHat(img image.Image, support [][]bool) *image.Gray {
	r := img.Bounds()
	return Sub(Close(img, support), img, r.Min)
}

// HitOrMiss keeps support1 and not support2 in the image. It requires that the
// intersection of the two supports be empty.
func HitOrMiss(img image.Image, support1, support2 [][]bool) *image.Gray {
	r := img.Bounds()
	return And(Erode(img, support1), Erode(Not(img), support2), r.Min)
}

// Thin thins the image by repeatedly subtracting HitOrMiss with the selected support pairs.
func Thin(img image.Image) *image.Gray {
	// Apply all eight pairs
	res := ThinStep(img, C11, D11)
	res = ThinStep(res, C21, D21)
	res = ThinStep(res, C12, D12)
	res = ThinStep(res, C22, D22)
	res = ThinStep(res, C13, D13)
	res = ThinStep(res, C23, D23)
	res = ThinStep(res, C14, D14)
	return ThinStep(res, C24, D24)
}

// ThinStep thins the image by subtracting HitOrMiss with a support pair.
func ThinStep(img image.Image, support1, support2 [][]bool) *image.Gray {
	r := img.Bounds()
	return Sub(img, HitOrMiss(img, support1, support2), r.Min)
}

// Skeleton repeatedly thins the image until it's no longer changing.
// This operation can take a while to converge.
func Skeleton(img image.Image) *image.Gray {
	r := img.Bounds()
	prev := image.Image(&image.Gray{})
	res := img
	for !Equal(res, prev, r.Min) {
		prev = res
		res = Thin(res)
	}
	return res.(*image.Gray)
}

/*
func Prune()
*/
