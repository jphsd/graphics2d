// go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	width, height := 500, 500

	img := image.NewRGBA(width, height, color.White)

	cp := []float64{250, 260}
	path := g2d.ReentrantPolygon(cp, 250, 5, 0.5, 0)

	hpp := &g2d.HandyProc{N: 3, R: 4}
	paths := path.Process(hpp)
	shape := g2d.NewShape(paths...)

	// Wood12 uses Catmull splines...
	rpp := &g2d.RoundedProc{10000}
	shape = shape.ProcessPaths(rpp)

	g2d.DrawShape(img, shape, g2d.BlackPen)

	image.SaveImage(img, "handy")
}
