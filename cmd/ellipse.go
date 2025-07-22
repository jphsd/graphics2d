//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Demonstrate effect of offs and ang on circle and ellipses
func main() {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)

	r := 400.0
	c := []float64{500, 500}
	a := g2d.Pi / 4
	offs := a
	ang := 2 * (g2d.Pi - a)
	//style := g2d.ArcPie
	style := g2d.ArcOpen
	factor := 0.5
	//factor := 0.618034 // Golden
	cs := g2d.NewShape(g2d.Arc(c, r, offs, ang, style))
	g2d.DrawShape(img, cs, g2d.BlackPen)
	e1s := g2d.NewShape(g2d.EllipticalArc(c, r, factor*r, offs, ang, 0, style))
	g2d.DrawShape(img, e1s, g2d.RedPen)
	e2s := g2d.NewShape(g2d.EllipticalArc(c, factor*r, r, offs, ang, 0, style))
	g2d.DrawShape(img, e2s, g2d.GreenPen)

	image.SaveImage(img, "ellipse")
}
