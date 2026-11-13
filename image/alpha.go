package image

import (
	"image"
	"image/color"
	"image/draw"
)

// NewAlpha is a wrapper for Alpha which returns a new image of the desired size filled with color.
func NewAlpha(w, h int, col color.Color) *Alpha {
	res := image.NewAlpha(Rect(0, 0, w, h))
	_, _, _, a := col.RGBA()
	if a != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyAlpha clones an Alpha image.
func CopyAlpha(in *Alpha) *Alpha {
	res := &Alpha{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewAlpha16 is a wrapper for Alpha16 which returns a new image of the desired size filled with color.
func NewAlpha16(w, h int, col color.Color) *Alpha16 {
	res := image.NewAlpha16(Rect(0, 0, w, h))
	_, _, _, a := col.RGBA()
	if a != 0 {
		draw.Draw(res, res.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
	}
	return res
}

// CopyAlpha16 clones an Alpha16 image.
func CopyAlpha16(in *Alpha16) *Alpha16 {
	res := &Alpha16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// AlphaToGray does a shallow copy (vs going through the ColorModel).
func AlphaToGray(a *Alpha) *Gray {
	return &Gray{a.Pix, a.Stride, a.Rect}
}

// GrayToAlpha does a shallow copy (vs going through the ColorModel).
func GrayToAlpha(a *Gray) *Alpha {
	return &Alpha{a.Pix, a.Stride, a.Rect}
}

// Alpha16ToGray16 does a shallow copy (vs going through the ColorModel).
func Alpha16ToGray16(a *Alpha16) *Gray16 {
	return &Gray16{a.Pix, a.Stride, a.Rect}
}

// Gray16ToAlpha16 does a shallow copy (vs going through the ColorModel).
func Gray16ToAlpha16(a *Gray16) *Alpha16 {
	return &Alpha16{a.Pix, a.Stride, a.Rect}
}

// AlphaToGrayDeep does a deep copy (vs going through the ColorModel).
func AlphaToGrayDeep(a *Alpha) *Gray {
	res := &Gray{make([]uint8, len(a.Pix)), a.Stride, a.Rect}
	copy(res.Pix, a.Pix)
	return res
}

// GrayToAlphaDeep does a deep copy (vs going through the ColorModel).
func GrayToAlphaDeep(g *Gray) *Alpha {
	res := &Alpha{make([]uint8, len(g.Pix)), g.Stride, g.Rect}
	copy(res.Pix, g.Pix)
	return res
}

// Alpha16ToGray16Deep does a deep copy (vs going through the ColorModel).
func Alpha16ToGray16Deep(a *Alpha16) *Gray16 {
	res := &Gray16{make([]uint8, len(a.Pix)), a.Stride, a.Rect}
	copy(res.Pix, a.Pix)
	return res
}

// Gray16ToAlpha16Deep does a deep copy (vs going through the ColorModel).
func Gray16ToAlpha16Deep(g *Gray16) *Alpha16 {
	res := &Alpha16{make([]uint8, len(g.Pix)), g.Stride, g.Rect}
	copy(res.Pix, g.Pix)
	return res
}
