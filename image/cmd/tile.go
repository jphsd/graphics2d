//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	twidth, theight := 30, 30
	timg := image.NewRGBA(twidth, theight, color.Black)

	g2d.FillPath(timg, g2d.Circle([]float64{7.5, 7.5}, 7), g2d.RedPen)
	g2d.FillPath(timg, g2d.Circle([]float64{22.5, 7.5}, 7), g2d.GreenPen)
	g2d.FillPath(timg, g2d.Circle([]float64{7.5, 22.5}, 7), g2d.BluePen)
	g2d.FillPath(timg, g2d.Circle([]float64{22.5, 22.5}, 7), g2d.YellowPen)

	tile := image.NewTile(timg)

	width, height := 600, 600
	img := image.NewRGBA(width, height, color.White)

	pen := &g2d.Pen{Filler: tile}
	circle := g2d.Circle([]float64{300, 300}, 290)
	g2d.FillPath(img, circle, pen)

	image.SaveImage(img, "tile")
}
