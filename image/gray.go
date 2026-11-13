package image

import (
	"image"
	"image/color"
	"image/draw"
)

// NewGray is a wrapper for Gray which returns a new image of the desired size filled with color.
func NewGray(w, h int, col color.Color) *Gray {
	res := image.NewGray(Rect(0, 0, w, h))
	r, _, _, _ := col.RGBA()
	if r != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyGray clones a Gray image.
func CopyGray(in *Gray) *Gray {
	res := &Gray{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewGray16 is a wrapper for Gray16 which returns a new image of the desired size filled with color.
func NewGray16(w, h int, col color.Color) *Gray16 {
	res := image.NewGray16(Rect(0, 0, w, h))
	r, _, _, _ := col.RGBA()
	if r != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyGray16 clones a Gray16 image.
func CopyGray16(in *Gray16) *Gray16 {
	res := &Gray16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// ToGray creates a grayscale copy of img
func ToGray(img Image) *Gray {
	r := img.Bounds()
	res := image.NewGray(r)
	draw.Draw(res, r, img, r.Min, draw.Src)
	return res
}

// ToGray16 creates a grayscale copy of img
func ToGray16(img Image) *Gray16 {
	r := img.Bounds()
	res := image.NewGray16(r)
	draw.Draw(res, r, img, r.Min, draw.Src)
	return res
}
