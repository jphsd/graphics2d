package image

import (
	"github.com/jphsd/graphics2d/util"
	"image"
	"image/color"
	"math"
)

// CircNonLinImage renders the supplied non-linear function into a circle to generate a Gray16 image
// of the specified dimensions.
func CircNonLinImage(w, h, inset int, f util.NonLinear, inv bool) *image.Gray16 {
	res := &image.Gray16{}
	if inv {
		res = NewGray16(w, h, color.White)
	} else {
		res = NewGray16(w, h, color.Black)
	}
	dx, dy := 2/float64(w-2*inset), 2/float64(h-2*inset)
	oy := -1.0
	for y := inset; y < h-inset; y++ {
		ox := -1.0
		oy2 := oy * oy
		for x := inset; x < w-inset; x++ {
			d := math.Sqrt(ox*ox + oy2)
			v := 0.0
			if d > 1 {
				ox += dx
				continue
			}
			v = f.Transform(1 - d)
			if inv {
				v = 1 - v
			}
			v *= 0xffff
			gcol := color.Gray16{uint16(v)}
			res.Set(x, y, gcol)
			ox += dx
		}
		oy += dy
	}
	return res
}
