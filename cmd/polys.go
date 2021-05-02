// +build ignore

package main

import (
	"image/color"
	"image/draw"
	"math"

	// For image output only
	"fmt"
	"image/png"
	"log"
	"os"

	. "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/image"
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
	shape := &Shape{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ns := i*n + j + 3
			dw := mdw * math.Tan(math.Pi/float64(ns))
			if dw > mdw {
				dw = mdw
			}
			shape.AddPaths(RegularPolygon([]float64{float64(cx) + dp1x + dw, float64(cy) + dp1y},
				[]float64{float64(cx) + dp1x - dw, float64(cy) + dp1y}, ns))
			cx += dx
		}
		cx = 0
		cy += dy
	}
	red := color.RGBA{0xff, 0, 0, 0xff}
	RenderColoredShape(img, shape, red)

	// Capture image output
	err := saveImage(img, "polys")
	if err != nil {
		log.Fatal(err)
	}
}

func saveImage(img draw.Image, name string) error {
	fDst, err := os.Create(fmt.Sprintf("%s.png", name))
	if err != nil {
		return err
	}
	defer fDst.Close()
	err = png.Encode(fDst, img)
	if err != nil {
		return err
	}
	return nil
}
