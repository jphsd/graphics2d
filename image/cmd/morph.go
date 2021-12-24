//go:build ignore

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"

	. "github.com/jphsd/graphics2d/image"
)

func main() {
	// Read in image file indicated in command line
	flag.Parse()
	args := flag.Args()
	img, err := ReadImage(args[0])
	if err != nil {
		panic(err)
	}

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
	SaveImage(gray, "out")

	// Perform a selection of operators
	suppt := Z8
	out := Open(gray, suppt)
	SaveImage(out, "out-open")
	out = Close(gray, suppt)
	SaveImage(out, "out-close")
	out = TopHat(gray, suppt)
	SaveImage(out, "out-top")
	out = BotHat(gray, suppt)
	SaveImage(out, "out-bot")
	out = gray
	prev := &image.Gray{}
	n := 0
	for !Equal(out, prev, image.Point{}) {
		prev = out
		out = Thin(out)
		n++
		//		if n%10 == 0 {
		//			SaveImage(out, fmt.Sprintf("out%d", n))
		//			v, n := Variance(prev, out)
		//			fmt.Printf("Variance %f over %d\n", v*float64(out.Bounds().Dx()*out.Bounds().Dy())/float64(n), n)
		//		}
	}
	SaveImage(Not(out), fmt.Sprintf("out-skel%d", n))

	outs := LJSkeleton(gray, Z4, 32)
	SaveImage(Not(outs[32]), fmt.Sprintf("out-ljskel%d", 32))

	out1 := LJReconstitute(outs[:32], Z4)
	SaveImage(Not(out1), "out-ljrecon")
}
