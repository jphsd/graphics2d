// +build ignore

package main

import (
	"flag"
	"fmt"
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
	_ "image/png"
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

	// Convert to Gray
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, image.Point{}, draw.Src)

	// Invert and undo anti-aliasing and compression
	lut := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		if i < 128 {
			lut[i] = 0xff
		} else {
			lut[i] = 0
		}
	}
	gray = RemapGray(gray, lut)
	saveImg(gray, "out")

	// Perform a selection of operators
	suppt := Z8
	out := Open(gray, suppt)
	saveImg(out, "out-open")
	out = Close(gray, suppt)
	saveImg(out, "out-close")
	out = TopHat(gray, suppt)
	saveImg(out, "out-top")
	out = BotHat(gray, suppt)
	saveImg(out, "out-bot")
	out = gray
	prev := &image.Gray{}
	n := 0
	for !Equal(prev, out) {
		prev = out
		out = Thin(out)
		n++
		//		if n%10 == 0 {
		//			saveImg(out, fmt.Sprintf("out%d", n))
		//			v, n := Variance(prev, out)
		//			fmt.Printf("Variance %f over %d\n", v*float64(out.Bounds().Dx()*out.Bounds().Dy())/float64(n), n)
		//		}
	}
	saveImg(Not(out), fmt.Sprintf("out-skel%d", n))

	outs := LJSkeleton(gray, Z4, 32)
	saveImg(Not(outs[32]), fmt.Sprintf("out-ljskel%d", 32))

	out1 := LJReconstitute(outs[:32], Z4)
	saveImg(Not(out1), "out-ljrecon")
}

func saveImg(img image.Image, name string) {
	fDst, err := os.Create(name + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer fDst.Close()
	err = png.Encode(fDst, img)
	if err != nil {
		log.Fatal(err)
	}
}
