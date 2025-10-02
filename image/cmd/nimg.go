//go:build ignore

package main

import (
	"flag"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
)

// Create a white image with default size 1000 by 1000.
func main() {
	wf := flag.Int("w", 1000, "width")
	hf := flag.Int("h", 1000, "height")
	flag.Parse()
	w, h := *wf, *hf
	args := flag.Args()
	var str string
	if len(args) == 0 {
		str = "untitled"
	} else {
		str = args[0]
	}

	img := image.NewRGBA(w, h, color.White)
	image.SaveImage(img, str)
}
