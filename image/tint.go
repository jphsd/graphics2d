package image

import (
	g2dcol "github.com/jphsd/graphics2d/color"
	"image"
	"image/color"
)

// Tinter uses the gray or red channel of an image to create an image tinted with a lerp'd gradient
// between two colors.
type Tinter struct {
	Img  image.Image
	Rect image.Rectangle
	Lut  []color.RGBA
	G16  bool
}

// NewTinter creates a new Tinter and creates the internal lut.
func NewTinter(img image.Image, c1, c2 color.Color) *Tinter {
	lut := make([]color.RGBA, 256)

	// Build lut - simple lerp between c1 and c2
	dt := 1.0 / 256.0
	t := 0.0
	for i := 0; i < 256; i++ {
		lut[i] = g2dcol.ColorRGBALerp(t, c1, c2)
		t += dt
	}
	cm := img.ColorModel()
	g16 := false
	if cm == color.Gray16Model {
		g16 = true
	}
	return &Tinter{img, img.Bounds(), lut, g16}
}

// ColorModel implements the ColorModel function in the Image interface.
func (t *Tinter) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the Bounds function in the Image interface.
func (t *Tinter) Bounds() image.Rectangle {
	return t.Rect
}

// At implements the At function in the Image interface.
func (t *Tinter) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(t.Rect)) {
		return color.RGBA{}
	}
	var v uint32
	c := t.Img.At(x, y)
	if t.G16 {
		gc := c.(color.Gray16)
		v = uint32(gc.Y)
	} else {
		v, _, _, _ = c.RGBA()
	}
	v = (v & 0xff00) >> 8
	return t.Lut[v]
}
