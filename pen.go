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

// NewPens constructs a slice of n pens with the given color starting at width and increasing by winc.
func NewPens(color color.Color, n int, width, winc float64) []*Pen {
	res := make([]*Pen, n)
	uc := image.NewUniform(color)
	for i := 0; i < n; i++ {
		res[i] = &Pen{uc, NewStrokeProc(width), nil}
		width += winc
	}
	return res
}

// NewColoredPens constructs a slice of n pens with the given colors and width.
func NewColoredPens(colors []color.Color, width float64) []*Pen {
	n := len(colors)
	res := make([]*Pen, n)
	for i := 0; i < n; i++ {
		uc := image.NewUniform(colors[i])
		res[i] = &Pen{uc, NewStrokeProc(width), nil}
	}
	return res
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

// NewFilledPen returns a pen that will render a shape with the given pen
// width and a random hued color into an image.
func NewFilledPen(grad image.Image, width float64) *Pen {
	return &Pen{grad, NewStrokeProc(width), nil}
}
