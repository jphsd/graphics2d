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
	p1 := []float64{100, 200}
	c1 := []float64{200, 0}
	c2 := []float64{200, 400}
	p2 := []float64{300, 200}

	path := NewPath(p1)
	path.AddStep(c1, c2, p2)

	/* We could have written this, this way too
	path, _ := PartsToPath([][]float64{p1, c1, c2, p2}})
	*/

	// Draw curve
	DrawPath(img, path, Red)

	// Capture image output
	image.SaveImage(img, "curve")
}
