package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/sfnt"
	"image/draw"
	"math/rand"
)

// Example_splash generates a series of background images using triangles, squares, pentagons,
// circles and stars, and then draws an outlined string over them.
func Example_splash() {
	// Make background fillers
	bg := make([]*image.RGBA, 5)
	bg[0] = makeRegularBackground(300, 3, 10, 30)
	bg[1] = makeRegularBackground(200, 4, 10, 30)
	bg[2] = makeRegularBackground(200, 5, 10, 30)
	bg[3] = makeCircleBackground(200, 10, 30)
	bg[4] = makeStarBackground(200, 5, 10, 30)

	// Create image with backgrounds
	img := image.NewRGBA(1000, 200, color.Transparent)
	for i, _ := range bg {
		rect := image.Rect(i*200, 0, i*200+200, 200)
		draw.Draw(img, rect, bg[i], image.Point{}, draw.Src)
	}

	// Load font and create shapes
	ttf, err := sfnt.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}
	str := "Graphics2D"
	shape, _, err := g2d.StringToShape(ttf, str)
	if err != nil {
		panic(err)
	}

	// Figure bounding box and scaling transform
	bb := shape.BoundingBox()
	xfm := g2d.ScaleAndInset(1000, 200, 20, 20, false, bb)
	shape = shape.ProcessPaths(xfm)

	// Render string to image
	pen := g2d.NewPen(color.White, 8)
	g2d.RenderColoredShape(img, shape, color.Black)
	g2d.DrawShape(img, shape, pen)
	image.SaveImage(img, "splash")

	fmt.Printf("See splash.png")
	// Output: See splash.png
}

func makeRegularBackground(n, s int, min, max float64) *image.RGBA {
	img := image.NewRGBA(200, 200, color.Black)

	dl := max - min
	for range n {
		c := []float64{rand.Float64() * 200, rand.Float64() * 200}
		l := min + rand.Float64()*dl
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		th := rand.Float64() * g2d.TwoPi
		shape := g2d.NewShape(g2d.RegularPolygon(s, c, l, th))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}

func makeCircleBackground(n int, min, max float64) *image.RGBA {
	img := image.NewRGBA(200, 200, color.Black)

	dr := max - min
	for range n {
		c := []float64{rand.Float64() * 200, rand.Float64() * 200}
		r := min + rand.Float64()*dr
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		shape := g2d.NewShape(g2d.Circle(c, r))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}

func makeStarBackground(n, s int, min, max float64) *image.RGBA {
	img := image.NewRGBA(200, 200, color.Black)

	dr := max - min
	for range n {
		c := []float64{rand.Float64() * 200, rand.Float64() * 200}
		r := min + rand.Float64()*dr
		col := color.HSL{rand.Float64(), 1, 0.5, 1}
		th := rand.Float64() * g2d.TwoPi
		shape := g2d.NewShape(g2d.ReentrantPolygon(c, r, s, 0.5, th))
		g2d.RenderColoredShape(img, shape, col)
	}

	return img
}
