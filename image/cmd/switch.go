//go:build ignore

package main

import (
	"flag"
	stdimg "image"
	"image/draw"

	"github.com/jphsd/graphics2d/image"
)

func main() {
	// Read in image file indicated in command line
	flag.Parse()
	args := flag.Args()
	img, err := image.ReadImage(args[0])
	if err != nil {
		panic(err)
	}

	// Convert to RGBA
	rgba := stdimg.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, stdimg.Point{}, draw.Src)

	imgR := image.SwitchChannels(rgba, 1, 2)
	imgG := image.SwitchChannels(rgba, 0, 2)
	imgB := image.SwitchChannels(rgba, 0, 1)

	image.SaveImage(imgR, "outR.png")
	image.SaveImage(imgG, "outG.png")
	image.SaveImage(imgB, "outB.png")
}
