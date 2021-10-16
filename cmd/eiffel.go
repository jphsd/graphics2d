//go:build ignore
// +build ignore

package main

import (
	. "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	width, height := 360, 400

	img := image.NewRGBA(width, height, color.White)

	// Make Eiffel shape
	path := NewPath([]float64{160, 80})
	path.AddStep([]float64{200, 80})
	path.AddStep([]float64{200, 200}, []float64{280, 320})
	path.AddStep([]float64{220, 320})
	path.AddStep([]float64{220, 240}, []float64{140, 240}, []float64{140, 320})
	path.AddStep([]float64{80, 320})
	path.AddStep([]float64{160, 200}, []float64{160, 80})
	path.Close()

	shape := NewShape(path)
	shape1 := shape.Transform(CreateTransform(-20, -20, 1, 0))
	path1 := path.Transform(CreateTransform(20, 20, 1, 0))

	// Render the shape in blue
	DrawShape(img, shape1, Blue)

	// and again offset in green
	DrawShape(img, shape, Green)

	// and again offset in red
	DrawPath(img, path1, Red)

	image.SaveImage(img, "eiffel")
}
