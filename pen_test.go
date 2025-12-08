package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Generates a series of regular shapes with a dashed pen.
func ExamplePen() {
	shape := g2d.NewShape(
		g2d.Line([]float64{20, 20}, []float64{130, 130}),
		g2d.RegularPolygon(3, []float64{225, 75}, 110, g2d.HalfPi),
		g2d.RegularPolygon(4, []float64{375, 75}, 110, 0),
		g2d.RegularPolygon(5, []float64{525, 75}, 75, 0),
		g2d.Circle([]float64{675, 75}, 55),
		g2d.Ellipse([]float64{825, 75}, 70, 35, g2d.HalfPi/2))

	img := image.NewRGBA(900, 150, color.White)
	dash := g2d.NewDashProc([]float64{8, 2, 2, 2}, 0)
	pen := g2d.NewProcessorPen(color.Black, 3, dash)
	g2d.DrawShape(img, shape, pen)
	image.SaveImage(img, "pen")

	fmt.Printf("See pen.png")
	// Output: See pen.png
}
