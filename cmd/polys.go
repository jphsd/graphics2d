//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math"
)

func main() {
	// Create image to write into
	width, height := 300, 300
	img := image.NewRGBA(width, height, color.White)

	n := 3
	dx, dy := width/n, height/n
	mdw := float64(dx) * 0.4
	dp1x, dp1y := float64(dx)*0.5, float64(dy)*0.9
	cx, cy := 0, 0
	shape := &g2d.Shape{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ns := i*n + j + 3
			dw := mdw * math.Tan(g2d.Pi/float64(ns))
			if dw > mdw {
				dw = mdw
			}
			shape.AddPaths(g2d.RegularPolygon([]float64{float64(cx) + dp1x + dw, float64(cy) + dp1y},
				[]float64{float64(cx) + dp1x - dw, float64(cy) + dp1y}, ns))
			cx += dx
		}
		cx = 0
		cy += dy
	}
	g2d.RenderColoredShape(img, shape, color.Red)
	g2d.DrawShape(img, shape, g2d.BlackPen)

	// Capture image output
	image.SaveImage(img, "polys")
}
