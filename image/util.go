package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	_ "image/png"
	//_ "github.com/adrium/goheif"
)

// NewRGBA is a wrapper for image.RGBA which returns a new image of the desired size filled with color.
func NewRGBA(w, h int, col color.Color) *RGBA {
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	r, g, b, _ := col.RGBA()
	if r != 0 || g != 0 || b != 0 {
		bg := NewUniform(col)
		draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
	}
	return res
}

// CopyRGBA clones an RGBA image.
func CopyRGBA(in *RGBA) *RGBA {
	res := &RGBA{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewRGBA64 is a wrapper for RGBA64 which returns a new image of the desired size filled with color.
func NewRGBA64(w, h int, col color.Color) *RGBA64 {
	res := image.NewRGBA64(Rect(0, 0, w, h))
	bg := NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
	return res
}

// CopyRGBA64 clones an RGBA64 image.
func CopyRGBA64(in *RGBA64) *RGBA64 {
	res := &RGBA64{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewAlpha is a wrapper for Alpha which returns a new image of the desired size filled with color.
func NewAlpha(w, h int, col color.Color) *Alpha {
	res := image.NewAlpha(Rect(0, 0, w, h))
	bg := NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
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
	bg := NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
	return res
}

// CopyAlpha16 clones an Alpha16 image.
func CopyAlpha16(in *Alpha16) *Alpha16 {
	res := &Alpha16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewGray is a wrapper for Gray which returns a new image of the desired size filled with color.
func NewGray(w, h int, col color.Color) *Gray {
	res := image.NewGray(Rect(0, 0, w, h))
	bg := NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
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
	bg := NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, Point{}, draw.Src)
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

// SaveImage is a utility function to save an image as a .png.
func SaveImage(img Image, name string) error {
	fDst, err := os.Create(fmt.Sprintf("%s.png", name))
	if err != nil {
		return err
	}
	defer fDst.Close()
	return png.Encode(fDst, img)
}

// ReadImage is a utility function to read an image from a file.
// The following formats are supported - bmp, jpeg, png, tiff, webp.
func ReadImage(name string) (Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
