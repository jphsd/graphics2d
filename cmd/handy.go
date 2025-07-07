//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	width, height := 500, 500

	img := image.NewRGBA(width, height, color.White)

	// Create a path and decompse it into its constituent parts
	cp := []float64{250, 260}
	path := g2d.ReentrantPolygon(cp, 250, 5, 0.5, 0) // Five pointed star
	paths := path.Process(&g2d.DecomposeProc{})

	// Create path processor, Wood12 uses Catmull splines
	hpp := &g2d.HandyProc{N: 4, R: 4}
	cpp := &g2d.CurveProc{0.375, g2d.CatmullRom}
	bpp := g2d.NewCompoundProc(hpp, cpp)

	// Create the shape and apply the path processor
	shape := g2d.NewShape(paths...)
	shape = shape.ProcessPaths(bpp)

	g2d.DrawShape(img, shape, g2d.BlackPen)
	//g2d.DrawPath(img, path, g2d.RedPen)

	image.SaveImage(img, "handy")
}
