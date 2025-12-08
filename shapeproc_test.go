package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Demonstrates using a shape processor by plotting a series of different colored stars along
// a circular path.
func ExampleBucketProc() {
	path := g2d.Circle([]float64{250, 250}, 200)
	star := g2d.NewShape(g2d.ReentrantPolygon([]float64{0, 0}, 20, 5, 0.5, 0))

	// Use ShapesProc path processor to draw a star every 50 pixels along the path
	pp := g2d.NewShapesProc([]*g2d.Shape{star}, 50, g2d.RotRandom)
	stars := g2d.NewShape(path.Process(pp)...)

	// Use BucketProc shape processor to take each path in stars and add it to one of n output shapes
	n := 6
	sp := g2d.BucketProc{N: n, Style: g2d.RoundRobin} // Other styles are g2d.Chunk and g2d.Random
	shapes := stars.Process(sp)

	img := image.NewRGBA(500, 500, color.White)
	for i, shape := range shapes {
		// Color shape by index into shapes
		pen := g2d.NewPen(color.HSL{float64(i+1) / float64(n), 1, 0.5, 1}, 3)
		g2d.DrawShape(img, shape, pen)
	}
	image.SaveImage(img, "stars")

	fmt.Printf("See stars.png")
	// Output: See stars.png
}
