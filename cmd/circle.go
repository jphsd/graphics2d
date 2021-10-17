//go:build ignore

package main

import (
	"math"

	. "github.com/jphsd/graphics2d"
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
	DrawArc(img, p1, c, math.Pi*2, Red)

	// Draw point at center in black
	DrawPoint(img, c, Black)

	// Capture image output
	image.SaveImage(img, "circle")
}
