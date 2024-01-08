package graphics2d

import (
	"image/draw"
)

// Simple immediate drawing and filling functions with pens

// DrawPoint renders a point with the pen into the destination image.
func DrawPoint(dst draw.Image, at []float64, pen *Pen) {
	DrawShape(dst, NewShape(Point(at)), pen)
}

// DrawLine renders a line with the pen into the destination image.
func DrawLine(dst draw.Image, start, end []float64, pen *Pen) {
	DrawShape(dst, NewShape(Line(start, end)), pen)
}

// DrawArc renders an arc with the pen into the destination image.
// radians +ve CCW, -ve CW
func DrawArc(dst draw.Image, start, center []float64, radians float64, pen *Pen) {
	DrawShape(dst, NewShape(ArcFromPoint(start, center, radians, ArcOpen)), pen)
}

// DrawPath renders a path with the pen into the destination image.
func DrawPath(dst draw.Image, path *Path, pen *Pen) {
	DrawShape(dst, NewShape(path), pen)
}

// DrawShape renders a shape with the pen into the destination image.
func DrawShape(dst draw.Image, shape *Shape, pen *Pen) {
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShape(dst, shape, pen.Filler)
}

// DrawClippedShape renders a shape with the pen against a clip shape into the destination image.
func DrawClippedShape(dst draw.Image, shape, clip *Shape, pen *Pen) {
	if pen.Stroke != nil {
		shape = shape.ProcessPaths(pen.Stroke)
	}
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderClippedShape(dst, shape, clip, pen.Filler)
}

// Fill functions ignore the pen stroke and if any path isn't closed, it's forced so.

// FillPath renders a path with the pen filler image and transform into the destination image.
func FillPath(dst draw.Image, path *Path, pen *Pen) {
	shape := NewShape(path)
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShape(dst, shape, pen.Filler)
}

// FillShape renders a shape with the pen filler and transform into the destination image.
func FillShape(dst draw.Image, shape *Shape, pen *Pen) {
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderShape(dst, shape, pen.Filler)
}

// FillClippedShape renders a shape with the pen filler against a clipe shape and transform into the destination image.
func FillClippedShape(dst draw.Image, shape, clip *Shape, pen *Pen) {
	if pen.Xfm != nil {
		shape = shape.Transform(pen.Xfm)
	}
	RenderClippedShape(dst, shape, clip, pen.Filler)
}
