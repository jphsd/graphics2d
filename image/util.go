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

// CopyRGBA clones an RGBA image.
func CopyRGBA(in *image.RGBA) *image.RGBA {
	res := &image.RGBA{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
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

// NewAlphaVal is a wrapper for image.Alpha which returns a new image of the desired size filled with color.
func NewAlphaVal(w, h int, a uint8) *image.Alpha {
	res := image.NewAlpha(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.Alpha{a})
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyAlpha clones an Alpha image.
func CopyAlpha(in *image.Alpha) *image.Alpha {
	res := &image.Alpha{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
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

// NewGrayVal is a wrapper for image.Gray which returns a new image of the desired size filled with color.
func NewGrayVal(w, h int, g uint8) *image.Gray {
	res := image.NewGray(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.Gray{g})
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

// NewGray16Val is a wrapper for image.Gray which returns a new image of the desired size filled with color.
func NewGray16Val(w, h int, g uint16) *image.Gray16 {
	res := image.NewGray16(image.Rect(0, 0, w, h))
	bg := image.NewUniform(color.Gray16{g})
	draw.Draw(res, res.Bounds(), bg, image.Point{}, draw.Src)
	return res
}

// CopyGray16 clones a Gray16 image.
func CopyGray16(in *image.Gray16) *image.Gray16 {
	res := &image.Gray16{make([]uint8, len(in.Pix)), in.Stride, in.Rect}
	copy(res.Pix, in.Pix)
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
