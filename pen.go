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
	BlackPen     = NewPen(g2dc.Black, 1)
	DarkGrayPen  = NewPen(g2dc.DarkGray, 1)
	GrayPen      = NewPen(g2dc.Gray, 1)
	LightGrayPen = NewPen(g2dc.LightGray, 1)
	WhitePen     = NewPen(g2dc.White, 1)
	RedPen       = NewPen(g2dc.Red, 1)
	GreenPen     = NewPen(g2dc.Green, 1)
	BluePen      = NewPen(g2dc.Blue, 1)
	YellowPen    = NewPen(g2dc.Yellow, 1)
	MagentaPen   = NewPen(g2dc.Magenta, 1)
	CyanPen      = NewPen(g2dc.Cyan, 1)
	OrangePen    = NewPen(g2dc.Orange, 1)
	BrownPen     = NewPen(g2dc.Brown, 1)
)

// NewPen returns a pen that will render a shape with the given pen
// width and color into an image.
func NewPen(color color.Color, width float64) *Pen {
	return &Pen{image.NewUniform(color), NewStrokeProc(width), nil}
}

// NewProcessorPen returns a pen that will render a shape with the given pen
// width and color into an image after applying the supplied path processor.
func NewProcessorPen(color color.Color, width float64, proc PathProcessor) *Pen {
	return &Pen{image.NewUniform(color), NewCompoundProc(proc, NewStrokeProc(width)), nil}
}

// NewNamedPen returns a pen that will render a shape with the given width and named color
// into an image. If the name is not matched then a black pen will be returned.
func NewNamedPen(name string, width float64) *Pen {
	col, err := g2dc.ByName(name)
	if err != nil {
		return &Pen{image.NewUniform(g2dc.Black), NewStrokeProc(width), nil}
	}
	return &Pen{image.NewUniform(col), NewStrokeProc(width), nil}
}

// NewRandomPen returns a pen that will render a shape with the given pen
// width and a random color into an image.
func NewRandomPen(width float64) *Pen {
	return &Pen{image.NewUniform(g2dc.Random()), NewStrokeProc(width), nil}
}

// NewRandomHuePen returns a pen that will render a shape with the given pen
// width and a random hued color into an image.
func NewRandomHuePen(width float64) *Pen {
	return &Pen{image.NewUniform(g2dc.RandomHue()), NewStrokeProc(width), nil}
}
