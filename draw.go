package graphics2d

import (
	"image"
	"image/draw"
)

// Simple drawing functions with pens

// DrawPoint renders a point with the pen into the destination image.
func DrawPoint(dst draw.Image, at []float64, pen *Pen) {
	shape := NewShape(Point(at).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, []float32{0, 0}, pen.Filler, image.Point{}, nil, draw.Over)
}

// DrawLine renders a line with the pen into the destination image.
func DrawLine(dst draw.Image, start, end []float64, pen *Pen) {
	shape := NewShape(Line(start, end).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, []float32{0, 0}, pen.Filler, image.Point{}, nil, draw.Over)
}

// DrawArc renders an arc with the pen into the destination image.
// radians +ve CCW, -ve CW
func DrawArc(dst draw.Image, start, center []float64, radians float64, pen *Pen) {
	shape := NewShape(ArcFromPoint(start, center, radians, ArcOpen).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, []float32{0, 0}, pen.Filler, image.Point{}, nil, draw.Over)
}

// DrawPath renders a path with the pen into the destination image.
func DrawPath(dst draw.Image, path *Path, pen *Pen) {
	shape := NewShape(path)
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, []float32{0, 0}, pen.Filler, image.Point{}, nil, draw.Over)
}

// DrawShape renders a shape with the pen into the destination image.
func DrawShape(dst draw.Image, shape *Shape, pen *Pen) {
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, []float32{0, 0}, pen.Filler, image.Point{}, nil, draw.Over)
}
