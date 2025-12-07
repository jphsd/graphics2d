package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Demonstrates using EllipticalArc to create a logo
func ExampleEllipticalArc() {
	width, height := 400, 400

	// Make first half
	c := []float64{200, 220}
	r := 150.0
	rxy := 0.4
	path1 := g2d.EllipticalArc(c, r, r, 0, g2d.Pi, 0, g2d.ArcOpen)
	path2 := g2d.EllipticalArc(c, r, r*rxy, g2d.Pi, -g2d.Pi, 0, g2d.ArcOpen)
	path3 := g2d.EllipticalArc(c, r, r*rxy, g2d.Pi, -g2d.TwoPi*0.7, 0, g2d.ArcOpen)
	p1 := []float64{c[0] - r, c[1]}
	// Pull p2 from path3
	p2 := path3.Current()
	r2 := r * 0.85
	path4 := g2d.EllipticalArcFromPoints2(p1, p2, r2, r2*rxy, 0, true, false, g2d.ArcOpen)

	// Combine paths into a shape
	shape := g2d.NewShape(path1, path2, path3, path4)

	// Rotate by 180 to get second half and add it
	xfm := g2d.RotateAbout(g2d.Pi, 200, 200)
	shape.AddShapes(shape.ProcessPaths(xfm))

	// Rotate shape CCW by 60
	xfm = g2d.RotateAbout(-g2d.Pi/6, 200, 200)
	shape = shape.ProcessPaths(xfm)

	img := image.NewRGBA(width, height, color.White)
	g2d.DrawShape(img, shape, g2d.BlackPen)
	image.SaveImage(img, "logo")

	fmt.Println("Check logo.png")
	// Output: Check logo.png
}
