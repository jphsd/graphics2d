package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/sfnt"
	"math/rand"
)

// Example_splash generates a series of background images using triangles, squares, pentagons,
// circles and stars, and uses them as fillers for the letters in a string.
func Example_splash() {
	// Make background fillers
	bg := make([]*image.RGBA, 5)
	bg[0] = makeRegularBackground(1500, 3, 10, 30)
	bg[1] = makeRegularBackground(1000, 4, 10, 30)
	bg[2] = makeRegularBackground(1000, 5, 10, 30)
	bg[3] = makeCircleBackground(1000, 10, 30)
	bg[4] = makeStarBackground(1000, 5, 10, 30)

	// Load font and create shapes
	ttf, err := sfnt.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}
	str := "Graphics2D"
	shape, shapes, err := g2d.StringToShape(ttf, str)
	if err != nil {
		panic(err)
	}

	// Figure bounding box and scaling transform
	bb := shape.BoundingBox()
	xfm := g2d.ScaleAndInset(1000, 200, 20, 20, false, bb)

	// Render to image
	img := image.NewRGBA(1000, 200, color.Black)
	for i, ss := range shapes {
		ss := ss.ProcessPaths(xfm)
		g2d.RenderShape(img, ss, bg[i%5])
	}
	image.SaveImage(img, "splash")

	fmt.Printf("See splash.png")
	// Output: See splash.png
}

func makeRegularBackground(n, s int, min, max float64) *image.RGBA {
	img := image.NewRGBA(1000, 200, color.White)

	dl := max - min
	for range n {
		c := []float64{rand.Float64() * 1000, rand.Float64() * 200}
		l := min + rand.Float64()*dl
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		th := rand.Float64() * g2d.TwoPi
		shape := g2d.NewShape(g2d.RegularPolygon(s, c, l, th))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}

func makeCircleBackground(n int, min, max float64) *image.RGBA {
	img := image.NewRGBA(1000, 200, color.White)

	dr := max - min
	for range n {
		c := []float64{rand.Float64() * 1000, rand.Float64() * 200}
		r := min + rand.Float64()*dr
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		shape := g2d.NewShape(g2d.Circle(c, r))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}

func makeStarBackground(n, s int, min, max float64) *image.RGBA {
	img := image.NewRGBA(1000, 200, color.White)

	dr := max - min
	for range n {
		c := []float64{rand.Float64() * 1000, rand.Float64() * 200}
		r := min + rand.Float64()*dr
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		th := rand.Float64() * g2d.TwoPi
		shape := g2d.NewShape(g2d.ReentrantPolygon(c, r, s, 0.5, th))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}
