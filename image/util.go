package image

import (
	"image"
	"image/color"
	"image/draw"
)

// NewRGBA is a wrapper for image.RGBA which returns a new image of the desired size filled with color.
func NewRGBA(w, h int, col color.Color) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// NewRGBAVal is a wrapper for image.RGBA which returns a new image of the desired size filled with color.
func NewRGBAVal(w, h int, r, g, b, a uint8) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.RGBA{r, g, b, a})
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// NewAlphaVal is a wrapper for image.Alpha which returns a new image of the desired size filled with color.
func NewAlphaVal(w, h int, a uint8) *image.Alpha {
	res := image.NewAlpha(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.Alpha{a})
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// NewGray is a wrapper for image.Gray which returns a new image of the desired size filled with color.
func NewGrayVal(w, h int, g uint8) *image.Gray {
	res := image.NewGray(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.Gray{g})
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}
