// +build ignore

package main

import (
	"image/color"
	"image/draw"
	"math"

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
	p1 := []float64{300, 200}
	c := []float64{200, 200}
	red := color.RGBA{0xff, 0, 0, 0xff}

	// Draw circle
	DrawArc(img, p1, c, math.Pi*2, red)

	// Draw point at center in black
	DrawPoint(img, c, color.Black)

	// Capture image output
	err := saveImage(img, "circle")
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
