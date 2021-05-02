// +build ignore

package main

import (
	"flag"
	"image"
	"image/draw"
	"os"

	. "github.com/jphsd/graphics2d/image"

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
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	imgR := ExtractChannel(rgba, 0)
	imgG := ExtractChannel(rgba, 1)
	imgB := ExtractChannel(rgba, 2)
	imgA := ExtractChannel(rgba, 3)
	rgba1 := CombineChannels(imgR, imgG, imgB, imgA, false)

	// Output images
	fRDst, err := os.Create("outR.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fRDst.Close()
	err = png.Encode(fRDst, imgR)
	if err != nil {
		log.Fatal(err)
	}
	fGDst, err := os.Create("outG.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fGDst.Close()
	err = png.Encode(fGDst, imgG)
	if err != nil {
		log.Fatal(err)
	}
	fBDst, err := os.Create("outB.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fBDst.Close()
	err = png.Encode(fBDst, imgB)
	if err != nil {
		log.Fatal(err)
	}
	fFDst, err := os.Create("outF.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fFDst.Close()
	err = png.Encode(fFDst, rgba1)
	if err != nil {
		log.Fatal(err)
	}
}
