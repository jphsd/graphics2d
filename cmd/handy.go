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

	// Create path processor
	// Wood12 uses Catmull splines
	cpp := &g2d.CurveProc{0.375, g2d.CatmullRom}
	hpp := &g2d.HandyProc{N: 3, R: 4}
	bpp := g2d.NewCompoundProc(hpp, cpp)

	shape := g2d.NewShape(path.Process(bpp)...)

	g2d.DrawShape(img, shape, g2d.BlackPen)

	image.SaveImage(img, "handy")
}
