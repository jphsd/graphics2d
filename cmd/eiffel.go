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
	width, height := 360, 400

	img := image.NewRGBA(width, height, color.White)

	// Make Eiffel shape
	path := NewPath([]float64{160, 80})
	path.AddStep([]float64{200, 80})
	path.AddStep([]float64{200, 200}, []float64{280, 320})
	path.AddStep([]float64{220, 320})
	path.AddStep([]float64{220, 240}, []float64{140, 240}, []float64{140, 320})
	path.AddStep([]float64{80, 320})
	path.AddStep([]float64{160, 200}, []float64{160, 80})
	path.Close()

	shape := NewShape(path)
	shape1 := shape.Transform(CreateTransform(-20, -20, 1, 0))
	path1 := path.Transform(CreateTransform(20, 20, 1, 0))

	// Render the shape in blue
	blue := color.RGBA{0, 0, 0xff, 0xff}
	RenderColoredShape(img, shape1, blue)

	// and again offset in green
	green := color.RGBA{0, 0xff, 0, 0xff}
	RenderColoredShape(img, shape, green)

	// and again offset in red
	red := color.RGBA{0xff, 0, 0, 0xff}
	RenderColoredPath(img, path1, red)

	err := saveImage(img, "eiffel")
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
