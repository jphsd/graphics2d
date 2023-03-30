package graphics2d

import (
	"image"
	"image/draw"
)

// Simple immediate drawing and filling functions with pens

// DrawPoint renders a point with the pen into the destination image.
func DrawPoint(dst draw.Image, at []float64, pen *Pen) {
	r := dst.Bounds()
	shape := NewShape(Point(at).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// DrawLine renders a line with the pen into the destination image.
func DrawLine(dst draw.Image, start, end []float64, pen *Pen) {
	r := dst.Bounds()
	shape := NewShape(Line(start, end).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// DrawArc renders an arc with the pen into the destination image.
// radians +ve CCW, -ve CW
func DrawArc(dst draw.Image, start, center []float64, radians float64, pen *Pen) {
	r := dst.Bounds()
	shape := NewShape(ArcFromPoint(start, center, radians, ArcOpen).Process(pen.Stroke)...)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// DrawPath renders a path with the pen into the destination image.
func DrawPath(dst draw.Image, path *Path, pen *Pen) {
	r := dst.Bounds()
	shape := NewShape(path)
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// DrawShape renders a shape with the pen into the destination image.
func DrawShape(dst draw.Image, shape *Shape, pen *Pen) {
	r := dst.Bounds()
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// Fill functions ignore the pen stroke and if any path isn't closed, it's forced so.

// FillPath renders a path with the pen filler image and transform into the destination image.
func FillPath(dst draw.Image, path *Path, pen *Pen) {
	r := dst.Bounds()
	shape := NewShape(path)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}

// FillShape renders a shape with the pen filler and transform into the destination image.
func FillShape(dst draw.Image, shape *Shape, pen *Pen) {
	r := dst.Bounds()
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShapeExt(dst, shape, pen.Filler, r.Min, nil, image.Point{}, draw.Over)
}
