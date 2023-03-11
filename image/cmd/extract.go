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

	imgR := image.ExtractChannel(rgba, 0)
	imgG := image.ExtractChannel(rgba, 1)
	imgB := image.ExtractChannel(rgba, 2)
	imgA := image.ExtractChannel(rgba, 3)
	rgba1 := image.CombineChannels(imgR, imgG, imgB, imgA, false)

	// Output images
	image.SaveImage(imgR, "outR")
	image.SaveImage(imgG, "outG")
	image.SaveImage(imgB, "outB")
	image.SaveImage(rgba1, "outF")
}
