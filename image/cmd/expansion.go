package main

import (
	"flag"
	"image"
	"image/draw"
	"os"

	. "github.com/jphsd/graphics2d/image"
	"github.com/jphsd/graphics2d/util"

	// For image output only
	"image/png"
	"log"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	//	_ "image/png"
)

func main() {
	// Read in image file indicated in command line
	flag.Parse()
	args := flag.Args()
	f, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	// Convert to RGBA
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	//	imgC := ExtractChannel(rgba, 0)
	//	imgC := ExtractChannel(rgba, 1)
	imgC := ExtractChannel(rgba, 2)

	// Create linear contrast expansion
	hist, start, end := Histogram(imgC)
	lcelut := NLExpansionLut(start, end+1, &util.NLLinear{})
	rgbaLCE := RemapRGBSingle(rgba, lcelut)

	// Create histogram equalization
	cdf := CDF(hist)
	helut := CreateLutFromValues(cdf)
	rgbaHE := RemapRGBSingle(rgba, helut)

	// Output images
	fLCEDst, err := os.Create("outLCE.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fLCEDst.Close()
	err = png.Encode(fLCEDst, rgbaLCE)
	if err != nil {
		log.Fatal(err)
	}
	fHEDst, err := os.Create("outHE.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fHEDst.Close()
	err = png.Encode(fHEDst, rgbaHE)
	if err != nil {
		log.Fatal(err)
	}
}
