package graphics2d

import (
	"github.com/jphsd/graphics2d/color"
	"image"
)

// Pen describes the color/image, stroke and shape to image transform to
// use when rendering shapes. If Stroke is nil then the shape's paths are
// used as is and forced closed (i.e. this is a fill). If Xfm is nil then
// the identity xfm is assumed.
type Pen struct {
	Filler image.Image
	Stroke PathProcessor
	Xfm    Transform
}

// Predefined pens.
var (
	BlackPen     = NewPen(color.Black, 1)
	DarkGrayPen  = NewPen(color.DarkGray, 1)
	GrayPen      = NewPen(color.MidGray, 1)
	LightGrayPen = NewPen(color.LightGray, 1)
	WhitePen     = NewPen(color.White, 1)
	RedPen       = NewPen(color.Red, 1)
	GreenPen     = NewPen(color.Green, 1)
	BluePen      = NewPen(color.Blue, 1)
	YellowPen    = NewPen(color.Yellow, 1)
	MagentaPen   = NewPen(color.Magenta, 1)
	CyanPen      = NewPen(color.Cyan, 1)
	OrangePen    = NewPen(color.Orange, 1)
	BrownPen     = NewPen(color.Brown, 1)
)

// NewPen returns a pen that will render a shape with the given pen
// width and color into an image.
// It uses the Stroke path processor to accomplish this with the bevel join and butt cap functions.
func NewPen(color color.Color, width float64) *Pen {
	return &Pen{image.NewUniform(color), NewStrokeProc(width), nil}
}

// NewProcessorPen returns a pen that will render a shape with the given pen
// width and color into an image after applying the supplied path processor.
func NewProcessorPen(color color.Color, width float64, proc PathProcessor) *Pen {
	return &Pen{image.NewUniform(color), NewCompoundProc(proc, NewStrokeProc(width)), nil}
}

// NewStrokedPen returns a pen that will render a shape with the given pen
// width and color into an image using the supplied join and cap functions.
func NewStrokedPen(color color.Color, width float64,
	join func(Part, []float64, Part) []Part,
	cap func(Part, []float64, Part) []Part) *Pen {
	if width < 0 {
		width = -width
	}
	width /= 2
	if join == nil {
		join = JoinBevel
	}
	if cap == nil {
		cap = CapButt
	}
	return &Pen{image.NewUniform(color), NewStrokeProcExt(width, -width, join, 0.5, cap), nil}
}
