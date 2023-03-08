//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Demonstrate different render/drawing methods
func main() {
	width, height := 360, 400

	img := image.NewRGBA(width, height, color.White)

	// Make Eiffel shape
	path := g2d.NewPath([]float64{160, 80})
	path.AddStep([]float64{200, 80})
	path.AddStep([]float64{200, 200}, []float64{280, 320})
	path.AddStep([]float64{220, 320})
	path.AddStep([]float64{220, 240}, []float64{140, 240}, []float64{140, 320})
	path.AddStep([]float64{80, 320})
	path.AddStep([]float64{160, 200}, []float64{160, 80})
	path.Close()

	shape := g2d.NewShape(path)
	shape1 := shape.Transform(g2d.CreateTransform(-20, -20, 1, 0))
	path1 := path.Transform(g2d.CreateTransform(20, 20, 1, 0))

	// Fill the shape with blue
	g2d.RenderColoredShape(img, shape1, color.Blue)

	// and again offset with the Green pen
	g2d.DrawShape(img, shape, g2d.GreenPen)

	// Draw the path offset with the Red pen
	g2d.DrawPath(img, path1, g2d.RedPen)

	image.SaveImage(img, "eiffel")
}
