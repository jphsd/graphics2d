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
	Xfm    *Aff3
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
	col, err := color.ByName(name)
	if err != nil {
		return &Pen{image.NewUniform(color.Black), NewStrokeProc(width), nil}
	}
	return &Pen{image.NewUniform(col), NewStrokeProc(width), nil}
}

// NewRandomPen returns a pen that will render a shape with the given pen
// width and a random color into an image.
func NewRandomPen(width float64) *Pen {
	return &Pen{image.NewUniform(color.Random()), NewStrokeProc(width), nil}
}

// NewRandomHuePen returns a pen that will render a shape with the given pen
// width and a random hued color into an image.
func NewRandomHuePen(width float64) *Pen {
	return &Pen{image.NewUniform(color.RandomHue()), NewStrokeProc(width), nil}
}

// NewFilledPen returns a pen that will render a shape with the given pen
// width and a random hued color into an image.
func NewFilledPen(grad image.Image, width float64) *Pen {
	return &Pen{grad, NewStrokeProc(width), nil}
}

// Width returns the width of this pen.
// Note this only works for pens with a StrokeProc path processor.
func (p *Pen) Width() float64 {
	sp, ok := p.Stroke.(*StrokeProc)
	if !ok {
		return -1
	}
	return sp.RTraceProc.Width - sp.LTraceProc.Width
}

// ChangeWidth returns a new pen the width while preserving the other aspects of the pen's stroke.
// Note this only works for pens with a StrokeProc path processor.
func (p *Pen) ChangeWidth(width float64) *Pen {
	sp, ok := p.Stroke.(*StrokeProc)
	if !ok {
		return nil
	}

	tsp := NewStrokeProc(width)
	tsp.PointFunc = sp.PointFunc
	tsp.CapStartFunc = sp.CapStartFunc
	tsp.CapEndFunc = sp.CapEndFunc
	tsp.RTraceProc.JoinFunc = sp.RTraceProc.JoinFunc
	tsp.LTraceProc.JoinFunc = sp.LTraceProc.JoinFunc

	if p.Xfm == nil {
		return &Pen{p.Filler, tsp, nil}
	}
	return &Pen{p.Filler, tsp, p.Xfm.Copy()}
}

// ScaleWidth returns a new pen the scaled width while preserving the other aspects of the pen's stroke.
// Note this only works for pens with a StrokeProc path processor.
func (p *Pen) ScaleWidth(scale float64) *Pen {
	sp, ok := p.Stroke.(*StrokeProc)
	if !ok {
		return nil
	}

	tsp := NewStrokeProc(p.Width() * scale)
	tsp.PointFunc = sp.PointFunc
	tsp.CapStartFunc = sp.CapStartFunc
	tsp.CapEndFunc = sp.CapEndFunc
	tsp.RTraceProc.JoinFunc = sp.RTraceProc.JoinFunc
	tsp.LTraceProc.JoinFunc = sp.LTraceProc.JoinFunc

	if p.Xfm == nil {
		return &Pen{p.Filler, tsp, nil}
	}
	return &Pen{p.Filler, tsp, p.Xfm.Copy()}
}
