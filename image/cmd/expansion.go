//go:build ignore

package main

import (
	"flag"
	stdimg "image"
	"image/draw"

	"github.com/jphsd/graphics2d/image"
	"github.com/jphsd/graphics2d/util"
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

	//	imgC := image.ExtractChannel(rgba, 0)
	//	imgC := image.ExtractChannel(rgba, 1)
	imgC := image.ExtractChannel(rgba, 2)

	// Create linear contrast expansion
	hist, start, end := image.Histogram(imgC)
	lcelut := image.NLExpansionLut(end-start+1, start, end+1, &util.NLLinear{})
	rgbaLCE := image.RemapRGBSingle(rgba, lcelut)

	// Create histogram equalization
	cdf := image.CDF(hist)
	helut := image.CreateLutFromValues(cdf)
	rgbaHE := image.RemapRGBSingle(rgba, helut)

	// Output images
	image.SaveImage(rgbaLCE, "outLCE")
	image.SaveImage(rgbaHE, "outHE")
}
