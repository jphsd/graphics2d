package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math"
)

// Demonstrates various eggs with different waists.
func ExampleEgg() {
	width, height := 1500, 300

	n := 5
	dx := float64(width / n)
	bb := [][]float64{{10, 10}, {dx - 10, float64(height) - 10}}
	c := []float64{0, 0}
	dt := 1.0 / float64(n+1)
	t := dt
	shape := &g2d.Shape{}
	for range n {
		path := g2d.Egg(c, 100, 200, t, 0)
		pbb := path.BoundingBox()
		// Fix pbb width
		pbb[0][0] = -100
		pbb[1][0] = 100
		xfm := g2d.BBTransform(pbb, bb)
		path = path.Process(xfm)[0]
		shape.AddPaths(path)
		bb[0][0] += dx
		bb[1][0] += dx
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

// Demonstrates various right eggs.
func ExampleRightEgg() {
	width, height := 1500, 300

	h := 200.0
	c := []float64{0, 0}
	dvals := []float64{
		// Values from Mathographics, Robert Dixon, 1987
		1.0 / (4.0 - math.Sqrt2),      // Moss
		2.0 / (3.0 + g2d.Sqrt3),       // Cundy Rollett
		7.0 / 16.0,                    // Thom1
		3.0 / (2.0 + math.Sqrt(10.0)), // Thom2
		1.0 / (1.0 + math.Phi),        // Golden
	}
	dx := float64(width / len(dvals))
	bb := [][]float64{{10, 10}, {dx - 10, float64(height) - 10}}

	shape := &g2d.Shape{}
	for _, d := range dvals {
		path := g2d.RightEgg(c, h, d, g2d.Pi) // Flip Y
		pbb := path.BoundingBox()
		// Fix pbb width
		pbb[0][0] = -100
		pbb[1][0] = 100
		xfm := g2d.BBTransform(pbb, bb)
		path = path.Process(xfm)[0]
		shape.AddPaths(path)
		bb[0][0] += dx
		bb[1][0] += dx
	}

	img := image.NewRGBA(width, height, color.White)
	g2d.RenderColoredShape(img, shape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)

	image.SaveImage(img, "regg")
	fmt.Println("Check regg.png")
	// Output: Check regg.png
}
