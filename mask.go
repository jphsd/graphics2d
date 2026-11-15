package graphics2d

import (
	"image"
	"image/color"
	"image/draw"
	//"github.com/jphsd/graphics2d/color"
)

// Draw is a wrapper around draw.Draw since attempting to alias it didn't work.
func Draw(dst Image, rect image.Rectangle, src image.Image, sp image.Point, op Op) {
	draw.Draw(dst, rect, src, sp, draw.Op(op))
}

// DrawMask is a wrapper around draw.DrawMask since attempting to alias it didn't work.
func DrawMask(dst Image, rect image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op Op) {
	draw.DrawMask(dst, rect, src, sp, mask, mp, draw.Op(op))
}

// DrawColor is a wrapper around Draw that uses a Uniform color as the src.
func DrawColor(dst Image, rect image.Rectangle, col color.Color, op Op) {
	Draw(dst, rect, &image.Uniform{col}, image.Point{}, op)
}

// DrawColorMask is a wrapper around DrawMask that uses a Uniform color as the src.
func DrawColorMask(dst Image, rect image.Rectangle, col color.Color, mask image.Image, mp image.Point, op Op) {
	DrawMask(dst, rect, &image.Uniform{col}, image.Point{}, mask, mp, op)
}

// NewRGBAFromMask returns a new RGBA image with the same bounds as the mask,
// filled by drawing the source through the mask.
func NewRGBAFromMask(src image.Image, sp image.Point, mask image.Image) *image.RGBA {
	rect := mask.Bounds()
	dst := image.NewRGBA(rect)
	DrawMask(dst, rect, src, sp, mask, rect.Min, Over)
	return dst
}
