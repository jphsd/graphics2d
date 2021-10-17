//go:build ignore

package main

import (
	. "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	// Create image to write into
	width, height := 400, 400
	img := image.NewRGBA(width, height, color.White)

	// Define points
	p1 := []float64{100, 100}
	p2 := []float64{300, 100}
	p3 := []float64{300, 300}
	p4 := []float64{100, 300}

	// Draw lines with the red pen
	DrawLine(img, p1, p2, Red)
	DrawLine(img, p2, p3, Red)
	DrawLine(img, p3, p4, Red)
	DrawLine(img, p4, p1, Red)

	// Capture image output
	image.SaveImage(img, "box")
}
