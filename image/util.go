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
func NewRGBA(w, h int, col color.Color) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyRGBA clones an RGBA image.
func CopyRGBA(in *image.RGBA) *image.RGBA {
	res := &image.RGBA{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewRGBA64 is a wrapper for image.RGBA64 which returns a new image of the desired size filled with color.
func NewRGBA64(w, h int, col color.Color) *image.RGBA64 {
	res := image.NewRGBA64(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyRGBA64 clones an RGBA64 image.
func CopyRGBA64(in *image.RGBA64) *image.RGBA64 {
	res := &image.RGBA64{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewAlpha is a wrapper for image.Alpha which returns a new image of the desired size filled with color.
func NewAlpha(w, h int, col color.Color) *image.Alpha {
	res := image.NewAlpha(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyAlpha clones an Alpha image.
func CopyAlpha(in *image.Alpha) *image.Alpha {
	res := &image.Alpha{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewAlpha16 is a wrapper for image.Alpha16 which returns a new image of the desired size filled with color.
func NewAlpha16(w, h int, col color.Color) *image.Alpha16 {
	res := image.NewAlpha16(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyAlpha16 clones an Alpha16 image.
func CopyAlpha16(in *image.Alpha16) *image.Alpha16 {
	res := &image.Alpha16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewGray is a wrapper for image.Gray which returns a new image of the desired size filled with color.
func NewGray(w, h int, col color.Color) *image.Gray {
	res := image.NewGray(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyGray clones a Gray image.
func CopyGray(in *image.Gray) *image.Gray {
	res := &image.Gray{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// NewGray16 is a wrapper for image.Gray16 which returns a new image of the desired size filled with color.
func NewGray16(w, h int, col color.Color) *image.Gray16 {
	res := image.NewGray16(image.Rect(0, 0, w, h))
	bg := image.NewUniform(col)
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyGray16 clones a Gray16 image.
func CopyGray16(in *image.Gray16) *image.Gray16 {
	res := &image.Gray16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
	return res
}

// ToGray creates a grayscale copy of img
func ToGray(img image.Image) *image.Gray {
	r := img.Bounds()
	res := image.NewGray(r)
	draw.Draw(res, r, img, r.Min, draw.Src)
	return res
}

// ToGray16 creates a grayscale copy of img
func ToGray16(img image.Image) *image.Gray16 {
	r := img.Bounds()
	res := image.NewGray16(r)
	draw.Draw(res, r, img, r.Min, draw.Src)
	return res
}

// AlphaToGray does a shallow copy (vs going through the ColorModel).
func AlphaToGray(a *image.Alpha) *image.Gray {
	return &image.Gray{a.Pix, a.Stride, a.Rect}
}

// GrayToAlpha does a shallow copy (vs going through the ColorModel).
func GrayToAlpha(a *image.Gray) *image.Alpha {
	return &image.Alpha{a.Pix, a.Stride, a.Rect}
}

// Alpha16ToGray16 does a shallow copy (vs going through the ColorModel).
func Alpha16ToGray16(a *image.Alpha16) *image.Gray16 {
	return &image.Gray16{a.Pix, a.Stride, a.Rect}
}

// Gray16ToAlpha16 does a shallow copy (vs going through the ColorModel).
func Gray16ToAlpha16(a *image.Gray16) *image.Alpha16 {
	return &image.Alpha16{a.Pix, a.Stride, a.Rect}
}

// AlphaToGrayDeep does a deep copy (vs going through the ColorModel).
func AlphaToGrayDeep(a *image.Alpha) *image.Gray {
	res := &image.Gray{make([]uint8, len(a.Pix)), a.Stride, a.Rect}
	copy(res.Pix, a.Pix)
	return res
}

// GrayToAlphaDeep does a deep copy (vs going through the ColorModel).
func GrayToAlphaDeep(g *image.Gray) *image.Alpha {
	res := &image.Alpha{make([]uint8, len(g.Pix)), g.Stride, g.Rect}
	copy(res.Pix, g.Pix)
	return res
}

// Alpha16ToGray16Deep does a deep copy (vs going through the ColorModel).
func Alpha16ToGray16Deep(a *image.Alpha16) *image.Gray16 {
	res := &image.Gray16{make([]uint8, len(a.Pix)), a.Stride, a.Rect}
	copy(res.Pix, a.Pix)
	return res
}

// Gray16ToAlpha16Deep does a deep copy (vs going through the ColorModel).
func Gray16ToAlpha16Deep(g *image.Gray16) *image.Alpha16 {
	res := &image.Alpha16{make([]uint8, len(g.Pix)), g.Stride, g.Rect}
	copy(res.Pix, g.Pix)
	return res
}

// SaveImage is a utility function to save an image as a .png.
func SaveImage(img image.Image, name string) error {
	fDst, err := os.Create(fmt.Sprintf("%s.png", name))
	if err != nil {
		return err
	}
	defer fDst.Close()
	return png.Encode(fDst, img)
}

// ReadImage is a utility function to read an image from a file.
func ReadImage(name string) (image.Image, error) {
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
