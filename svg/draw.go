package svg

import (
	"encoding/xml"
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
)

type xshape struct {
	Color string `xml:"fill,attr"`
	Paths []*g2d.Path
}

// RenderColoredShape writes SVG describing the shape and its color to the encoder.
func RenderColoredShape(enc *xml.Encoder, shape *g2d.Shape, col color.Color) error {
	r, g, b, _ := col.RGBA() // [0, 0xffff]
	r >>= 8
	g >>= 8
	b >>= 8
	cstr := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return enc.EncodeElement(xshape{cstr, shape.Paths()}, xml.StartElement{Name: xml.Name{"", "g"}})
}

// DrawShape writes SVG describing the shape as rendered by the pen to the encoder.
func DrawShape(enc *xml.Encoder, shape *g2d.Shape, pen *g2d.Pen) error {
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	col := pen.Filler.At(0, 0) // Assumes the filler image has constant color
	return RenderColoredShape(enc, shape, col)
}

// RenderRenderable renders the shapes in the renderable with their respective fillers after being transformed.
func RenderRenderable(enc *xml.Encoder, r *g2d.Renderable, xfm g2d.Transform) error {
	for i, shape := range r.Shapes {
		clip := r.Clips[i]
		if xfm != nil {
			if clip == nil {
				shape = shape.Transform(xfm)
				RenderColoredShape(enc, shape, r.Fillers[i].At(0, 0))
			} else {
				return fmt.Errorf("Clipping not supported")
			}
		} else if clip == nil {
			RenderColoredShape(enc, shape, r.Fillers[i].At(0, 0))
		} else {
			return fmt.Errorf("Clipping not supported")
		}
	}
	return nil
}
