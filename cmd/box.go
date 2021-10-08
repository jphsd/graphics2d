// +build ignore

package main

import (
	"image/color"
	"image/draw"

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
	width, height := 400, 400
	img := image.NewRGBA(width, height, color.White)

	// Define points
	p1 := []float64{100, 100}
	p2 := []float64{300, 100}
	p3 := []float64{300, 300}
	p4 := []float64{100, 300}
	red := NewPen(1, color.RGBA{0xff, 0, 0, 0xff})

	// Draw lines
	DrawLineP(img, p1, p2, red)
	DrawLineP(img, p2, p3, red)
	DrawLineP(img, p3, p4, red)
	DrawLineP(img, p4, p1, red)

	// Capture image output
	err := saveImage(img, "box")
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
