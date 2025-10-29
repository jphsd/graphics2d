//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	// Create image to write into
	width, height := 900, 900

	n := 3
	dx, dy := float64(width/n), float64(height/n)
	cx, cy := dx/2, dy/2
	shape := &g2d.Shape{}
	for i := range n {
		for j := range n {
			ns := i*n + j + 3
			path := g2d.RegularPolygon(ns, []float64{cx, cy}, 50, 0)
			shape.AddPaths(path)
			cx += dx
		}
		cx = dx / 2
		cy += dy
	}

	img := image.NewRGBA(width, height, color.White)
	g2d.FillShape(img, shape, g2d.RedPen)
	image.SaveImage(img, "polys")
}
