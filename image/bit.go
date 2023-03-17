package image

import (
	"github.com/jphsd/datastruct"
	"image"
	"image/color"
)

// Bit image uses a bit array to store data. The default colors assigned to bit values
// are color.White for set and color.Black for clear. These can be modified.
type Bit struct {
	Bits   datastruct.Bits
	SetC   color.RGBA // color returned when bit is set
	ClearC color.RGBA // color returned when bit is clear
	Stride int
	Thresh uint32 // threshold above which a bit is set in the RGB conversion
	Rect   image.Rectangle
}

// NewTile creates a new image with the supplied image tile.
func NewBit(r image.Rectangle) *Bit {
	w, h := r.Dx(), r.Dy()
	bits := datastruct.NewBits(w * h)
	return &Bit{bits, color.RGBA{0xff, 0xff, 0xff, 0xff}, color.RGBA{0, 0, 0, 0xff}, w, 0x8000, r}
}

// ColorModel implements the ColorModel function in the Image interface.
func (b *Bit) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the Bounds function in the Image interface.
func (b *Bit) Bounds() image.Rectangle {
	return b.Rect
}

// At implements the At function in the Image interface.
func (b *Bit) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(b.Rect)) {
		return color.RGBA{}
	}
	// Convert x, y to bit index
	x -= b.Rect.Min.X
	y -= b.Rect.Min.Y
	i := x + y*b.Stride
	if b.Bits.Get(i) {
		return b.SetC
	}
	return b.ClearC
}

// Set implements the Set function in the draw.Image interface.
func (b *Bit) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	// Convert x, y to bit index
	x -= b.Rect.Min.X
	y -= b.Rect.Min.Y
	i := x + y*b.Stride
	rc, gc, bc, _ := c.RGBA()
	sum := rc + gc + bc
	sum /= 3
	if sum < b.Thresh {
		b.Bits.Clear(i)
	} else {
		b.Bits.Set(i)
	}
}
