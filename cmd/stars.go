//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	// Create image to write into
	width, height := 400, 400
	img := image.NewRGBA(width, height, color.White)

	n := 5
	dx, dy := width/n, height/n
	cx, cy := dx/2, dy/2
	r := float64(dx) / 2 * 0.9
	a := 0.0
	shape := &g2d.Shape{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cp := []float64{float64(cx), float64(cy)}
			shape.AddPaths(g2d.ReentrantPolygon(cp, r, i+3, float64(j*25+1)/100.0, a))
			cx += dx
		}
		cx = dx / 2
		cy += dy
	}
	g2d.RenderColoredShape(img, shape, color.Red)
	g2d.DrawShape(img, shape, g2d.BlackPen)

	// Capture image output
	image.SaveImage(img, "stars")
}
