//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	// Create image to write into
	width, height := 400, 400
	img := image.NewRGBA(width, height, color.White)

	// Create a path describing the box
	//box := g2d.RegularPolygon(4, []float64{200, 200}, 200, 0)
	box := g2d.Rectangle([]float64{200, 200}, 200, 200)
	g2d.DrawPath(img, box, g2d.RedPen)

	// Capture image output
	image.SaveImage(img, "box")
}
