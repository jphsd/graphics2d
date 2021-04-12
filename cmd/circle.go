package main

import (
	"image"
	"image/color"
	"math"

	// For image output only
	"fmt"
	"image/png"
	"log"
	"os"

	. "github.com/jphsd/graphics2d"
	g2dimg "github.com/jphsd/graphics2d/image"
)

func main() {
	// Create image to write into
	width, height := 400, 400
	img := g2dimg.NewRGBA(width, height, color.White)

	// Define points
	p1 := []float64{300, 200}
	c := []float64{200, 200}
	red := color.RGBA{0xff, 0, 0, 0xff}

	// Draw lines
	DrawArc(img, p1, c, math.Pi*2, red)

	// Capture image output
	err := saveImage(img, "out")
	if err != nil {
		log.Fatal(err)
	}
}

func saveImage(img *image.RGBA, name string) error {
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
