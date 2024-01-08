//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math"
	"math/rand"
)

// Demonstrate use of IrregularEllipses to draw ink brush style lines
func main() {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)

	points := [][]float64{
		{100, 100},
		{900, 100},
		{900, 900},
		{100, 900},
	}

	for i := 0; i < len(points)-1; i++ {
		path := InkStroke(points[i], points[i+1], 20, 0, 1)
		g2d.RenderColoredShape(img, g2d.NewShape(path), color.GopherBlue)
	}
	path := InkStroke(points[len(points)-1], points[0], 20, 0, 1)
	g2d.RenderColoredShape(img, g2d.NewShape(path), color.GopherBlue)

	image.SaveImage(img, "inkbox")
}

func InkStroke(p1, p2 []float64, mw, dw, ec float64) *g2d.Path {
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	th := math.Atan2(dy, dx)
	d := math.Hypot(dx, dy)
	mw /= 2

	rx2 := rand.Float64() * 0.1 * d
	rx1 := d - rx2
	ry1 := rand.Float64()*dw + mw
	ry2 := rand.Float64()*dw + mw
	disp := rand.Float64()

	t := rx2 / d
	omt := 1 - t
	c := []float64{p1[0]*omt + p2[0]*t, p1[1]*omt + p2[1]*t}

	d1 := rx1 * disp
	rx1 += ec
	rx2 += ec
	disp = d1 / rx1

	return g2d.IrregularEllipse(c, rx1, rx2, ry1, ry2, disp, th)
}
