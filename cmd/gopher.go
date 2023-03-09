//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	stdimg "image"
)

// The face of a gopher (FKA Gordon)
// From Renee French's original design http://reneefrench.blogspot.com/
func main() {
	// Create image to write into
	width, height := 400, 370
	img := image.NewRGBA(width, height, color.GopherBlue)

	// Face is made up of eyes, teeth and nose
	DrawEye(img, []float64{100, 150})
	DrawEye(img, []float64{300, 150})
	DrawTeeth(img, []float64{200, 240})
	DrawNose(img, []float64{200, 225})

	// Capture image output
	image.SaveImage(img, "gopher")
}

func DrawEye(img *stdimg.RGBA, offs []float64) {
	eye1 := g2d.NewShape(g2d.Circle(offs, 75))
	g2d.RenderColoredShape(img, eye1, color.White)
	g2d.DrawShape(img, eye1, g2d.BlackPen)

	iris := g2d.NewShape(g2d.Ellipse(offs, 25, 27, 0))
	iris.AddPaths(g2d.Circle([]float64{offs[0] + 10, offs[1] + 10}, 5).Reverse())
	xfm := g2d.NewAff3()
	xfm.Translate(-25, 25)
	eye2 := iris.Transform(xfm)
	g2d.RenderColoredShape(img, eye2, color.Black)
}

func DrawTeeth(img *stdimg.RGBA, offs []float64) {
	// Use a RoundedProc path processor on a rectangle to create a tooth
	tooth := g2d.NewShape(
		g2d.Rectangle([]float64{offs[0] - 10, offs[1] + 25}, 20, 50).Process(&g2d.RoundedProc{8})...)
	g2d.RenderColoredShape(img, tooth, color.White)
	g2d.DrawShape(img, tooth, g2d.BlackPen)

	// Use a transform to draw the second tooth using the first
	xfm := g2d.NewAff3()
	xfm.Translate(20, 0)
	tooth = tooth.Transform(xfm)
	g2d.RenderColoredShape(img, tooth, color.White)
	g2d.DrawShape(img, tooth, g2d.BlackPen)
}

func DrawNose(img *stdimg.RGBA, offs []float64) {
	// Draw underside first
	points := [][]float64{
		{offs[0] - 20, offs[1]},
		{offs[0] + 20, offs[1]},
		{offs[0] + 40, offs[1] + 25},
		{offs[0] + 20, offs[1] + 35},
		{offs[0], offs[1] + 30},
		{offs[0] - 20, offs[1] + 35},
		{offs[0] - 40, offs[1] + 25},
	}
	// Use a RoundedProc path processor on a polygon rather than try and figure out control points
	bottom := g2d.NewShape(g2d.Polygon(points...).Process(&g2d.RoundedProc{30})...)
	g2d.RenderColoredShape(img, bottom, color.GopherBrown)
	g2d.DrawShape(img, bottom, g2d.BlackPen)

	// Draw the top
	top := g2d.NewShape(g2d.Ellipse(offs, 30, 15, 0))
	g2d.RenderColoredShape(img, top, color.GopherGray)
	g2d.DrawShape(img, top, g2d.BlackPen)
}
