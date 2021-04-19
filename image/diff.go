package image

import (
	"fmt"
	"image"
)

// Variance provides a normalized metric on how different two images are. It's not a strict
// statistical variance but the sum of the differences over the image size and number of differences.
func Variance(img1, img2 *image.Gray) (float64, int) {
	img1R, img2R := img1.Bounds(), img2.Bounds()
	w, h := img1R.Dx(), img1R.Dy()
	// Bounds check
	if w != img2R.Dx() || h != img2R.Dy() {
		panic(fmt.Errorf("image sizes don't match"))
	}
	n := 0
	sum := float64(0)
	for y := 0; y < h; y++ {
		i1 := img1.PixOffset(0+img1R.Min.X, y+img1R.Min.Y)
		i2 := img2.PixOffset(0+img2R.Min.X, y+img2R.Min.Y)
		for x := 0; x < w; x, i1, i2 = x+1, i1+1, i2+1 {
			v1, v2 := img1.Pix[i1], img2.Pix[i2]
			if v1 != v2 {
				if v1 > v2 {
					sum += float64(v1-v2) / 255
				} else {
					sum += float64(v2-v1) / 255
				}
				n++
			}
		}
	}
	return sum / float64(w*h), n
}
