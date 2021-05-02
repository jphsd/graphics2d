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
	p1 := []float64{100, 200}
	c1 := []float64{200, 0}
	c2 := []float64{200, 400}
	p2 := []float64{300, 200}
	red := color.RGBA{0xff, 0, 0, 0xff}

	path := NewPath(p1)
	err := path.AddStep(c1, c2, p2)
	if err != nil {
		log.Fatal(err)
	}

	/* We could have written this, this way too
	path, err := PartsToPath([][]float64{p1, c1, c2, p2}})
	*/

	// Draw curve
	DrawPath(img, path, red)

	// Capture image output
	err = saveImage(img, "curve")
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
