//go:build ignore

package main

import (
	"image/draw"

	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	width, height := 600, 600

	rbpatch, _ := image.NewPatch([][]color.Color{{color.Red, color.Blue}, {color.Blue, color.Red}})
	brpatch, _ := image.NewPatch([][]color.Color{{color.Green, color.Red}, {color.Red, color.Green}})

	img := image.NewRGBA(width, height, color.White)
	draw.Draw(img, img.Bounds(), rbpatch, image.Point{0, 0}, draw.Src)

	pen := &g2d.Pen{brpatch, nil, nil}
	circle := g2d.Circle([]float64{300, 300}, 290)
	g2d.FillPath(img, circle, pen)

	image.SaveImage(img, "patch")
}
