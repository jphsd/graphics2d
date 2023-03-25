package image

import (
	"image"
	"image/color"
	"math"
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
		lut[i] = ColorRGBALerp(t, c1, c2)
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

// TODO - is this the right place for this?

// ColorRGBALerp calculates the color value at t [0,1] given a start and end color in RGB space.
func ColorRGBALerp(t float64, start, end color.Color) color.RGBA {
	rs, gs, bs, as := start.RGBA() // uint32 [0,0xffff]
	re, ge, be, ae := end.RGBA()
	rt := uint32(math.Floor((1-t)*float64(rs) + t*float64(re) + 0.5))
	gt := uint32(math.Floor((1-t)*float64(gs) + t*float64(ge) + 0.5))
	bt := uint32(math.Floor((1-t)*float64(bs) + t*float64(be) + 0.5))
	at := uint32(math.Floor((1-t)*float64(as) + t*float64(ae) + 0.5))
	rt >>= 8 // uint32 [0,0xff]
	gt >>= 8
	bt >>= 8
	at >>= 8
	return color.RGBA{uint8(rt), uint8(gt), uint8(bt), uint8(at)}
}
