package main

import (
	//	glfw
	//  glui
	"image/color"

	// For image output only
	"image/png"
	"log"
	"os"

	. "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/image"
)

func main() {
	width, height := 360, 400

	// Initialize the GLFW system

	img := image.NewRGBA(width, height, color.White)

	// Make shape
	path := NewPath([]float64{160, 80})
	path.AddStep([][]float64{{200, 80}})
	path.AddStep([][]float64{{200, 200}, {280, 320}})
	path.AddStep([][]float64{{220, 320}})
	path.AddStep([][]float64{{220, 240}, {140, 240}, {140, 320}})
	path.AddStep([][]float64{{80, 320}})
	path.AddStep([][]float64{{160, 200}, {160, 80}})
	path.Close()

	shape := &Shape{}
	shape.AddPath(path)

	// Render the shape in blue
	blue := color.RGBA{0, 0, 0xff, 0xff}
	RenderColoredShape(img, shape, []float32{-20, -20}, blue)

	// and again offset in green
	green := color.RGBA{0, 0xff, 0, 0xff}
	RenderColoredShape(img, shape, []float32{0, 0}, green)

	// and again offset in red
	red := color.RGBA{0xff, 0, 0, 0xff}
	RenderColoredPath(img, path, []float32{20, 20}, red)

	// Display it on screen
	// Well, in an image for the time being then
	fDst, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fDst.Close()
	err = png.Encode(fDst, img)
	if err != nil {
		log.Fatal(err)
	}
}
