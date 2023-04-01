//go:build ignore

package main

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math/rand"
)

var (
	origin = []float64{0, 0}
	black  = g2d.NewPen(color.Black, 4)
)

// The gopher (FKA Gordon)
// From Renee French's original design http://reneefrench.blogspot.com/
func main() {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)

	// Create the gopher renderable
	gopher := MakeGopher()
	rect := gopher.Bounds()
	fmt.Printf("%d,%d -> %d,%d (%d,%d)\n", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, rect.Dx(), rect.Dy())

	y := 125.0
	for i := 0; i < 4; i++ {
		x := 100.0
		for j := 0; j < 5; j++ {
			// Locate and scale
			xfm := g2d.NewAff3()
			xfm.Translate(x, y)
			xfm.Rotate((rand.Float64()*2 - 1) * 0.3)
			xfm.Scale(0.25, 0.25)

			// Render
			gopher.Render(img, xfm)
			x += 200
		}
		y += 250
	}

	// Capture image output
	image.SaveImage(img, "gopher2")
}

func MakeGopher() *g2d.Renderable {
	res := &g2d.Renderable{}

	// Legs
	limb := MakeLimb()
	xfm := g2d.NewAff3()
	xfm.Translate(-165, 410)
	xfm.Rotate(0.6)
	res.AddRenderable(limb, xfm)
	xfm = g2d.NewAff3()
	xfm.Translate(165, 410)
	xfm.Rotate(-0.6)
	res.AddRenderable(limb, xfm)

	// Ears
	ear := MakeEar()
	xfm = g2d.NewAff3()
	xfm.Translate(-260, -350)
	xfm.Rotate(-0.6)
	res.AddRenderable(ear, xfm)
	xfm = g2d.NewAff3()
	xfm.Translate(260, -350)
	xfm.Rotate(0.6)
	res.AddRenderable(ear, xfm)

	// Body
	body := MakeBody()
	res.AddRenderable(body, nil)

	// Eyes
	eye := MakeEye()
	xfm = g2d.NewAff3()
	xfm.Translate(-170, -200)
	res.AddRenderable(eye, xfm)
	xfm = g2d.NewAff3()
	xfm.Translate(170, -200)
	res.AddRenderable(eye, xfm)

	// Teeth
	tooth := MakeTooth()
	xfm = g2d.NewAff3()
	xfm.Translate(-14, 0)
	res.AddRenderable(tooth, xfm)
	xfm = g2d.NewAff3()
	xfm.Translate(14, 0)
	res.AddRenderable(tooth, xfm)

	// Nose
	nose := MakeNose()
	xfm = g2d.NewAff3()
	xfm.Translate(0, -80)
	res.AddRenderable(nose, xfm)

	// Arms
	xfm = g2d.NewAff3()
	xfm.Translate(-200, 100)
	xfm.Rotate(-0.8)
	res.AddRenderable(limb, xfm)
	xfm = g2d.NewAff3()
	xfm.Translate(200, 100)
	xfm.Rotate(0.8)
	res.AddRenderable(limb, xfm)

	return res
}

func MakeLimb() *g2d.Renderable {
	res := &g2d.Renderable{}

	limb := g2d.Rectangle(origin, 55, 130).Process(&g2d.RoundedProc{30})[0]
	res.AddColoredShape(g2d.NewShape(limb), color.GopherBrown, nil)
	snip := g2d.NewPathSnipProc(g2d.Line([]float64{-50, -40}, []float64{50, -40}))
	paths := limb.Process(snip)
	res.AddPennedShape(g2d.NewShape(paths[1]), black, nil)

	return res
}

func MakeEar() *g2d.Renderable {
	res := &g2d.Renderable{}

	outter := g2d.NewShape(g2d.Ellipse(origin, 45, 60, 0))
	inner := g2d.NewShape(g2d.Ellipse(origin, 25, 40, 0))
	res.AddColoredShape(outter, color.GopherBlue, nil)
	res.AddColoredShape(inner, color.GopherGray, nil)
	res.AddPennedShape(outter, black, nil)

	return res
}

func MakeBody() *g2d.Renderable {
	res := &g2d.Renderable{}

	path := g2d.Rectangle([]float64{0, 0}, 680, 860)
	body := g2d.NewShape(path.Process(&g2d.CurveProc{0.7, g2d.Bezier})...)
	res.AddColoredShape(body, color.GopherBlue, nil)
	res.AddPennedShape(body, black, nil)

	return res
}

func MakeEye() *g2d.Renderable {
	res := &g2d.Renderable{}

	eye := g2d.NewShape(g2d.Circle(origin, 130))
	res.AddColoredShape(eye, color.White, nil)
	res.AddPennedShape(eye, black, nil)

	iris := g2d.NewShape(g2d.Ellipse([]float64{10, 50}, 45, 50, 0))
	iris.AddPaths(g2d.Circle([]float64{35, 40}, 10).Reverse())
	res.AddColoredShape(iris, color.Black, nil)

	return res
}

func MakeTooth() *g2d.Renderable {
	res := &g2d.Renderable{}
	// Use a RoundedProc path processor on a rectangle to create a tooth
	tooth := g2d.NewShape(
		g2d.Rectangle([]float64{0, 0}, 28, 70).Process(&g2d.RoundedProc{15})...)
	res.AddColoredShape(tooth, color.White, nil)
	res.AddPennedShape(tooth, black, nil)

	return res
}

func MakeNose() *g2d.Renderable {
	res := &g2d.Renderable{}
	// Draw underside first
	points := [][]float64{
		{-40, 0},
		{40, 0},
		{60, 45},
		{40, 70},
		{0, 60},
		{-40, 70},
		{-60, 45},
	}
	// Use a RoundedProc path processor on a polygon rather than try and figure out control points
	bottom := g2d.NewShape(g2d.Polygon(points...).Process(&g2d.RoundedProc{30})...)
	res.AddColoredShape(bottom, color.GopherBrown, nil)
	res.AddPennedShape(bottom, black, nil)
	//g2d.RenderColoredShape(img, bottom, color.GopherBrown)
	//g2d.DrawShape(img, bottom, black)

	// Draw the top
	top := g2d.NewShape(g2d.Ellipse(origin, 40, 23, 0))
	res.AddColoredShape(top, color.GopherGray, nil)
	res.AddPennedShape(top, black, nil)

	return res
}
