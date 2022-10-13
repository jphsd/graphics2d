//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math"
)

// Demonstrate effect of offs and ang on circle and ellipses
func main() {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)

	r := 200.0
	c := []float64{500, 500}
	a := math.Pi / 4
	offs := a
	ang := 2 * (math.Pi - a)
	cs := g2d.NewShape(g2d.Arc(c, r, offs, ang, g2d.ArcPie))
	g2d.DrawShape(img, cs, g2d.BlackPen)
	e1s := g2d.NewShape(g2d.EllipticalArc(c, r, 2*r, offs, ang, 0, g2d.ArcPie))
	g2d.DrawShape(img, e1s, g2d.RedPen)
	e2s := g2d.NewShape(g2d.EllipticalArc(c, 2*r, r, offs, ang, 0, g2d.ArcPie))
	g2d.DrawShape(img, e2s, g2d.GreenPen)

	image.SaveImage(img, "ellipse")
}
