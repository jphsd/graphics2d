package graphics2d

import (
	"image/color"
	"image/draw"
)

// Simple drawing functions

// DrawPoint renders a point into dst.
func DrawPoint(dst draw.Image, at []float64, color color.Color) {
	DrawPointW(dst, at, 1, color)
}

// DrawPointW renders a point of size width into dst.
func DrawPointW(dst draw.Image, at []float64, width float64, color color.Color) {
	RenderColoredShape(dst, NewShape(
		Point(at).Process(NewStrokeProc(width))...),
		color)
}

// DrawLine renders a line into dst.
func DrawLine(dst draw.Image, start, end []float64, color color.Color) {
	DrawLineW(dst, start, end, 1, color)
}

// DrawLineW renders a line of size width into dst.
func DrawLineW(dst draw.Image, start, end []float64, width float64, color color.Color) {
	RenderColoredShape(dst, NewShape(
		Line(start, end).Process(NewStrokeProc(width))...),
		color)
}

// DrawArc renders an arc into dst.
// radians +ve CCW, -ve CW
func DrawArc(dst draw.Image, start, center []float64, radians float64, color color.Color) {
	DrawArcW(dst, start, center, radians, 1, color)
}

// DrawArcW renders an arc of size width into dst.
func DrawArcW(dst draw.Image, start, center []float64, radians float64, width float64, color color.Color) {
	RenderColoredShape(dst, NewShape(
		ArcFromPoint(start, center, radians, ArcOpen).Process(NewStrokeProc(width))...),
		color)
}

// DrawPath renders a path into dst.
func DrawPath(dst draw.Image, path *Path, color color.Color) {
	DrawPathW(dst, path, 1, color)
}

// DrawPathW renders a path of size width into dst.
func DrawPathW(dst draw.Image, path *Path, width float64, color color.Color) {
	RenderColoredShape(dst, NewShape(
		path.Process(NewStrokeProc(width))...),
		color)
}
