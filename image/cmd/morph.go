//go:build ignore

package main

import (
	"flag"
	"fmt"
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

	// Convert to Gray
	gray := stdimg.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, stdimg.Point{}, draw.Src)

	// Invert and undo anti-aliasing and compression
	lut := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		if i < 128 {
			lut[i] = 0xff
		} else {
			lut[i] = 0
		}
	}
	gray = image.RemapGray(gray, lut)
	image.SaveImage(gray, "out")

	// Perform a selection of operators
	suppt := image.Z8
	out := image.Open(gray, suppt)
	image.SaveImage(out, "out-open")
	out = image.Close(gray, suppt)
	image.SaveImage(out, "out-close")
	out = image.TopHat(gray, suppt)
	image.SaveImage(out, "out-top")
	out = image.BotHat(gray, suppt)
	image.SaveImage(out, "out-bot")
	out = gray
	prev := &stdimg.Gray{}
	n := 0
	for !image.Equal(out, prev, stdimg.Point{}) {
		prev = out
		out = image.Thin(out)
		n++
		//		if n%10 == 0 {
		//			image.SaveImage(out, fmt.Sprintf("out%d", n))
		//			v, n := Variance(prev, out)
		//			fmt.Printf("Variance %f over %d\n", v*float64(out.Bounds().Dx()*out.Bounds().Dy())/float64(n), n)
		//		}
	}
	image.SaveImage(image.Not(out), fmt.Sprintf("out-skel%d", n))

	outs := image.LJSkeleton(gray, image.Z4, 32)
	image.SaveImage(image.Not(outs[32]), fmt.Sprintf("out-ljskel%d", 32))

	out1 := image.LJReconstitute(outs[:32], image.Z4)
	image.SaveImage(image.Not(out1), "out-ljrecon")
}
