package image

import (
	"image"
	"image/color"
	"image/draw"
)

// NewRGBA is a wrapper for image.RGBA which returns a new image of the desired size filled with color.
func NewRGBA(w, h int, col color.Color) *RGBA {
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	_, _, _, a := col.RGBA()
	if a != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyRGBA clones an RGBA image.
func CopyRGBA(in *RGBA) *RGBA {
	res := &RGBA{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// SetAlpha sets the alpha channel of an RGBA image to the supplied alpha value.
func SetAlpha(in *RGBA, c color.Color) *RGBA {
	_, _, _, a := c.RGBA()
	a >>= 8
	rect := in.Bounds()
	dx, dy := rect.Dx(), rect.Dy()
	x, y := rect.Min.X, rect.Min.Y
	offs := in.PixOffset(x, y)
	for range dy {
		coffs := offs + 3
		for range dx {
			in.Pix[coffs] = uint8(a)
			coffs += 4
		}
		offs += in.Stride
	}

	return in
}

// NewRGBA64 is a wrapper for RGBA64 which returns a new image of the desired size filled with color.
func NewRGBA64(w, h int, col color.Color) *RGBA64 {
	res := image.NewRGBA64(Rect(0, 0, w, h))
	_, _, _, a := col.RGBA()
	if a != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyRGBA64 clones an RGBA64 image.
func CopyRGBA64(in *RGBA64) *RGBA64 {
	res := &RGBA64{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// SetAlpha16 sets the alpha channel of an RGBA64 image to the supplied alpha value.
func SetAlpha16(in *RGBA64, c color.Color) *RGBA64 {
	_, _, _, a := c.RGBA()
	ah, al := a>>8, a&0xff
	rect := in.Bounds()
	dx, dy := rect.Dx(), rect.Dy()
	x, y := rect.Min.X, rect.Min.Y
	offs := in.PixOffset(x, y)
	for range dy {
		coffs := offs + 6
		for range dx {
			in.Pix[coffs] = uint8(ah)
			in.Pix[coffs+1] = uint8(al)
			coffs += 8
		}
		offs += in.Stride
	}

	return in
}
