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

	// Define points
	p1 := []float64{300, 200}
	c := []float64{200, 200}

	// Draw circle
	g2d.DrawArc(img, p1, c, g2d.TwoPi, g2d.RedPen)

	// Draw point at center in black
	g2d.DrawPoint(img, c, g2d.BlackPen)

	// Capture image output
	image.SaveImage(img, "circle")
}
