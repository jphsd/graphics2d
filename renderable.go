package graphics2d

import (
	"image"
	"image/color"
	"image/draw"
)

// Renderable represents a set of shapes and the images to fill them. In other words, enough information to be
// able to render something. This structure is used to build complex multicolored objects in a composable way.
type Renderable struct {
	Shapes  []*Shape
	Fillers []image.Image
}

// NewRenderable creates a new instance with the given shape and filler image.
func NewRenderable(shape *Shape, filler image.Image, xfm *Aff3) *Renderable {
	res := &Renderable{}
	res.AddShape(shape, filler, xfm)
	return res
}

// AddShape adds the given shape and filler to the Renderable after being transformed.
func (r *Renderable) AddShape(shape *Shape, filler image.Image, xfm *Aff3) *Renderable {
	if xfm != nil {
		r.Shapes = append(r.Shapes, shape.Transform(xfm))
	} else {
		r.Shapes = append(r.Shapes, shape)
	}
	r.Fillers = append(r.Fillers, filler)
	return r
}

// AddColoredShape adds the given shape and color to the Renderable after being transformed.
func (r *Renderable) AddColoredShape(shape *Shape, col color.Color, xfm *Aff3) *Renderable {
	if xfm != nil {
		r.Shapes = append(r.Shapes, shape.Transform(xfm))
	} else {
		r.Shapes = append(r.Shapes, shape)
	}
	r.Fillers = append(r.Fillers, image.NewUniform(col))
	return r
}

// AddPennedShape adds the given shape and pen to the Renderable after being transformed.
func (r *Renderable) AddPennedShape(shape *Shape, pen *Pen, xfm *Aff3) *Renderable {
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	if xfm != nil {
		r.Shapes = append(r.Shapes, shape.Transform(xfm))
	} else {
		r.Shapes = append(r.Shapes, shape)
	}
	r.Fillers = append(r.Fillers, pen.Filler)
	return r
}

// AddRenderable allows another renderable to be concatenated (post transform) to the current one.
func (r *Renderable) AddRenderable(rend *Renderable, xfm *Aff3) *Renderable {
	for i, shape := range rend.Shapes {
		r.AddShape(shape, rend.Fillers[i], xfm)
	}
	return r
}

// Render renders the shapes in the renderable with their respective fillers after being transformed.
func (r *Renderable) Render(img draw.Image, xfm *Aff3) {
	for i, shape := range r.Shapes {
		if xfm != nil {
			RenderShape(img, shape.Transform(xfm), r.Fillers[i], 0, 0)
		} else {
			RenderShape(img, shape, r.Fillers[i], 0, 0)
		}
	}
}

// Bounds returns the extent of the renderable.
func (r *Renderable) Bounds() image.Rectangle {
	rect := image.Rectangle{}
	for _, shape := range r.Shapes {
		rect = rect.Union(shape.Bounds())
	}
	return rect
}
