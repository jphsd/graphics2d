package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Demonstrates various eggs with different waists.
func ExampleEgg() {
	width, height := 1500, 300

	n := 5
	dx := float64(width / n)
	cx, cy := dx/2, float64(height)/2
	dt := 1.0 / float64(n+1)
	t := dt
	shape := &g2d.Shape{}
	for range n {
		c := []float64{cx, cy}
		shape.AddPaths(g2d.Egg(c, 150, 150, t, 0))
		cx += dx
		t += dt
	}

	img := image.NewRGBA(width, height, color.White)
	g2d.RenderColoredShape(img, shape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)

	image.SaveImage(img, "egg")
	fmt.Println("Check egg.png")
	// Output: Check egg.png
}
