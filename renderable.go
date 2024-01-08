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
	Clips   []*Shape
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
	return r.AddClippedShape(shape, nil, filler, xfm)
}

// AddClippedShape adds the given shape, clip and filler to the Renderable after being transformed.
func (r *Renderable) AddClippedShape(shape, clip *Shape, filler image.Image, xfm *Aff3) *Renderable {
	if xfm != nil {
		r.Shapes = append(r.Shapes, shape.Transform(xfm))
		if clip != nil {
			r.Clips = append(r.Clips, clip.Transform(xfm))
		} else {
			r.Clips = append(r.Clips, nil)
		}
	} else {
		r.Shapes = append(r.Shapes, shape)
		if clip != nil {
			r.Clips = append(r.Clips, clip)
		} else {
			r.Clips = append(r.Clips, nil)
		}
	}
	r.Fillers = append(r.Fillers, filler)
	return r
}

// AddColoredShape adds the given shape and color to the Renderable after being transformed.
func (r *Renderable) AddColoredShape(shape *Shape, col color.Color, xfm *Aff3) *Renderable {
	return r.AddClippedColoredShape(shape, nil, col, xfm)
}

// AddClippedColoredShape adds the given shape, clip and color to the Renderable after being transformed.
func (r *Renderable) AddClippedColoredShape(shape, clip *Shape, col color.Color, xfm *Aff3) *Renderable {
	return r.AddClippedShape(shape, clip, image.NewUniform(col), xfm)
}

// AddPennedShape adds the given shape and pen to the Renderable after being transformed.
func (r *Renderable) AddPennedShape(shape *Shape, pen *Pen, xfm *Aff3) *Renderable {
	return r.AddClippedPennedShape(shape, nil, pen, xfm)
}

// AddClippedPennedShape adds the given shape, clip and pen to the Renderable after being transformed.
func (r *Renderable) AddClippedPennedShape(shape, clip *Shape, pen *Pen, xfm *Aff3) *Renderable {
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}

	return r.AddClippedShape(shape, clip, pen.Filler, xfm)
}

// AddRenderable allows another renderable to be concatenated (post transform) to the current one.
func (r *Renderable) AddRenderable(rend *Renderable, xfm *Aff3) *Renderable {
	for i, shape := range rend.Shapes {
		r.AddClippedShape(shape, rend.Clips[i], rend.Fillers[i], xfm)
	}
	return r
}

// Render renders the shapes in the renderable with their respective fillers after being transformed.
func (r *Renderable) Render(img draw.Image, xfm *Aff3) {
	for i, shape := range r.Shapes {
		clip := r.Clips[i]
		if xfm != nil {
			if clip == nil {
				RenderShape(img, shape.Transform(xfm), r.Fillers[i])
			} else {
				RenderClippedShape(img, shape.Transform(xfm), clip.Transform(xfm), r.Fillers[i])
			}
		} else if clip == nil {
			RenderShape(img, shape, r.Fillers[i])
		} else {
			RenderClippedShape(img, shape, clip, r.Fillers[i])
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
