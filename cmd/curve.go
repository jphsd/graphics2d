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
	points := [][]float64{
		{100, 200},
		{200, 0},
		{200, 400},
		{300, 200},
	}

	path := g2d.PartsToPath(points)

	// Draw curve
	g2d.DrawPath(img, path, g2d.RedPen)

	// Capture image output
	image.SaveImage(img, "curve")
}
