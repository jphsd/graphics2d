package image

import "image"

// Histogram computes the histogram of the image and the
// first and last non-zero entry positions
func Histogram(img *image.Gray) ([]int, int, int) {
	res := make([]int, 256)
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()

	// Count frequency of each value
	sp := img.Pix
	if img.Stride == w {
		// All at once
		for i := 0; i < len(sp); i++ {
			res[sp[i]]++
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so := img.PixOffset(0, i)
			for j := 0; j < w; j++ {
				res[sp[so+j]]++
			}
		}
	}

	// Find first and last
	first := -1
	last := 0
	for i := 0; i < len(res); i++ {
		if res[i] > 0 {
			last = i
			if first == -1 {
				first = i
			}
		}
	}

	return res, first, last
}

// CDF computes the normalized cummulative distribution function of the
// histogram.
func CDF(hist []int) []float64 {
	res := make([]float64, len(hist))

	// Create the cummulative sum and capture the first non-zero
	// value (CDFmin)
	minv := 0.0
	minp := -1
	sum := 0
	for i := 0; i < len(hist); i++ {
		v := hist[i]
		sum += v
		res[i] = float64(sum)
		if minp == -1 && v > 0 {
			minv = float64(v)
			minp = i
		}
	}

	// Normalize result
	div := 1 / (float64(sum) - minv)
	for i := minp; i < len(res); i++ {
		res[i] = (res[i] - minv) * div
	}

	return res
}
