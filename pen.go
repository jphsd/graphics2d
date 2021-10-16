package graphics2d

import (
	g2dc "github.com/jphsd/graphics2d/color"
	"image"
	"image/color"
)

// Pen describes the color/image, stroke and shape to image transform to
// use when rendering shapes. If Stroke is nil then the shape's paths are
// used as is and forced closed (i.e. this is a fill). If Xfm is nil then
// the identity xfm is assumed.
type Pen struct {
	Filler image.Image
	Stroke PathProcessor
	Xfm    *Aff3
}

// Predefined pens.
var (
	Black   = NewPen(g2dc.Black, 1)
	White   = NewPen(g2dc.White, 1)
	Red     = NewPen(g2dc.Red, 1)
	Green   = NewPen(g2dc.Green, 1)
	Blue    = NewPen(g2dc.Blue, 1)
	Yellow  = NewPen(g2dc.Yellow, 1)
	Magenta = NewPen(g2dc.Magenta, 1)
	Cyan    = NewPen(g2dc.Cyan, 1)
	Orange  = NewPen(g2dc.Orange, 1)
)

// NewPen returns a pen that will render a shape with the given pen
// width and color into an image.
func NewPen(color color.Color, width float64) *Pen {
	return &Pen{image.NewUniform(color), NewStrokeProc(width), nil}
}