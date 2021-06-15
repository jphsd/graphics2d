package image

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Displace displaces the image within the rectangle by the amounts referenced in the other
// two images (one each for X and Y) scaled. If the displaced location exceeds the bounds of
// the rectangle, it is wrapped. Black maps to scale * -0.5 and white to scale * 0.5.
func Displace(img draw.Image, r image.Rectangle, xchan image.Image, xoffs image.Point, xscale float64,
	ychan image.Image, yoffs image.Point, yscale float64) {
	dx, dy := r.Dx(), r.Dy()
	tmp := image.NewRGBA(image.Rect(0, 0, dx, dy))
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			fx := col2f64(xchan.At(x+xoffs.X, y+xoffs.Y)) * xscale
			nx := x + int(math.Floor(fx+0.5))
			for nx < r.Min.X {
				nx += dx
			}
			for nx > r.Max.X {
				nx -= dx
			}
			fy := col2f64(ychan.At(x+yoffs.X, y+yoffs.Y)) * yscale
			ny := y + int(math.Floor(fy+0.5))
			for ny < r.Min.Y {
				ny += dy
			}
			for ny > r.Max.Y {
				ny -= dy
			}
			tmp.Set(x-r.Min.X, y-r.Min.Y, img.At(nx, ny))
		}
	}
	// Replace original with displaced
	draw.Draw(img, r, tmp, image.Point{}, draw.Src)
}

// Convert color into float64 [-0.5,0.5]
func col2f64(c color.Color) float64 {
	g16c, _ := color.Gray16Model.Convert(c).(color.Gray16)
	res := float64(g16c.Y) / 0xffff
	return res - 0.5
}
