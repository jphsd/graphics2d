package graphics2d

import (
	"image/color"
	"image/draw"
)

// Simple drawing functions

// DrawPoint renders a point into img.
func DrawPoint(img draw.Image, at []float64, color color.Color) {
	DrawPointW(img, at, 1, color)
}

// DrawPointW renders a point of size width into img.
func DrawPointW(img draw.Image, at []float64, width float64, color color.Color) {
	RenderColoredShape(img, NewShape(
		Point(at).Process(NewStrokeProc(width))...),
		color)
}

// DrawLine renders a line into img.
func DrawLine(img draw.Image, start, end []float64, color color.Color) {
	DrawLineW(img, start, end, 1, color)
}

// DrawLineW renders a line of size width into img.
func DrawLineW(img draw.Image, start, end []float64, width float64, color color.Color) {
	RenderColoredShape(img, NewShape(
		Line(start, end).Process(NewStrokeProc(width))...),
		color)
}

// DrawArc renders an arc into img.
// radians +ve CCW, -ve CW
func DrawArc(img draw.Image, start, center []float64, radians float64, color color.Color) {
	DrawArcW(img, start, center, radians, 1, color)
}

// DrawArcW renders an arc of size width into img.
func DrawArcW(img draw.Image, start, center []float64, radians float64, width float64, color color.Color) {
	RenderColoredShape(img, NewShape(
		ArcFromPoint(start, center, radians, ArcOpen).Process(NewStrokeProc(width))...),
		color)
}

// DrawPath renders a path into img.
func DrawPath(img draw.Image, path *Path, color color.Color) {
	DrawPathW(img, path, 1, color)
}

// DrawPathW renders a path of size width into img.
func DrawPathW(img draw.Image, path *Path, width float64, color color.Color) {
	RenderColoredShape(img, NewShape(
		path.Process(NewStrokeProc(width))...),
		color)
}
