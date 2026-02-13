//go:build ignore

package main

import (
	"bytes"
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/svg"
	"math"
	"math/rand"
)

var (
	// Colors
	brown     = color.RGBA{156, 126, 99, 255}
	brownhigh = color.RGBA{225, 191, 158, 255}
	brownshad = color.RGBA{97, 60, 31, 255}
	gray      = color.RGBA{176, 170, 157, 255}
	grayhigh  = color.RGBA{221, 216, 207, 255}
	grayshad  = color.RGBA{120, 111, 93, 255}

	// Structures
	first = [][]float64{
		{30, 210},
		{230, 210},
		{230, 110},
		{30, 110},
	}
	second = [][]float64{
		{30, 110},
		{230, 110},
		{130, 20},
	}
	stringer = [][]float64{
		{40, 100},
		{220, 100},
	}
	studs = [][]float64{
		// These are in pairs
		{70, 100}, {70, 70},
		{110, 100}, {110, 30},
		{150, 100}, {150, 30},
		{190, 100}, {190, 70},
	}
	roof = [][]float64{
		{0, 120},
		{20, 120},
		{130, 20},
		{240, 120},
		{260, 120},
		{130, 0},
	}
	chimney = [][]float64{
		{200, 110},
		{230, 110},
		{230, 30},
		{200, 30},
	}
	door = [][]float64{
		{150, 210},
		{200, 210},
		{200, 130},
		{150, 130},
	}
	window = [][]float64{
		{60, 180},
		{120, 180},
		{120, 130},
		{60, 130},
	}
)

// Demonstrate use of IrregularEllipses and Renderables to draw house
func main() {
	// Make house
	house := MakeHouse()

	// Locate and scale
	xfm := g2d.NewAff3()
	xfm.Translate(250, 250)
	xfm.Scale(2, 2)

	b := &bytes.Buffer{}
	enc := svg.NewEncoder(b)
	svg.RenderRenderable(enc, house, xfm)
	svg.Complete(b)
	fmt.Println(b.String())
}

func MakeHouse() *g2d.Renderable {
	// Order is important - furthest from eye objects first
	res := MakeChimney()
	res.AddRenderable(MakeFirstFloor(), nil)
	res.AddRenderable(MakeSecondFloor(), nil)
	res.AddRenderable(MakeRoof(), nil)

	return res
}

func MakeFirstFloor() *g2d.Renderable {
	res := &g2d.Renderable{}

	// Flats
	poly := g2d.Polygon(first...)
	res.AddColoredShape(g2d.NewShape(poly), gray, nil)
	poly = g2d.Polygon(door...)
	res.AddColoredShape(g2d.NewShape(poly), brown, nil)

	// Shadows and highlights
	shad := InkStroke([]float64{first[2][0] - 10, first[2][1]}, []float64{first[1][0] - 10, first[1][1]}, 7, 5, 0)
	res.AddColoredShape(g2d.NewShape(shad), grayshad, nil)
	high := InkStroke([]float64{first[3][0] + 10, first[3][1]}, []float64{first[0][0] + 10, first[0][1]}, 7, 5, 0)
	res.AddColoredShape(g2d.NewShape(high), grayhigh, nil)

	// Lines
	lshape := g2d.NewShape(PolyStroke(first, 3, 3, 2)...)
	lshape.AddPaths(PolyStroke(window, 1, 2, 2)...)
	lshape.AddPaths(PolyStroke(door, 1, 2, 2)...)
	res.AddColoredShape(lshape, color.Black, nil)

	return res
}

func MakeSecondFloor() *g2d.Renderable {
	res := &g2d.Renderable{}

	// Flats
	poly := g2d.Polygon(second...)
	res.AddColoredShape(g2d.NewShape(poly), brown, nil)

	// Shadows
	shad1 := InkStroke([]float64{second[2][0] + 5, second[2][1]}, []float64{second[0][0] + 5, second[0][1]}, 3, 4, 0)
	shad2 := InkStroke([]float64{second[2][0] - 5, second[2][1]}, []float64{second[1][0] - 5, second[1][1]}, 3, 4, 0)
	res.AddColoredShape(g2d.NewShape(shad1, shad2), brownshad, nil)

	// Lines
	lshape := g2d.NewShape(PolyStroke(second, 3, 3, 2)...)
	lshape.AddPaths(InkStroke(stringer[0], stringer[1], 1, 2, 2))
	lshape.AddPaths(InkStroke(studs[0], studs[1], 1, 2, 2))
	lshape.AddPaths(InkStroke(studs[2], studs[3], 1, 2, 2))
	lshape.AddPaths(InkStroke(studs[4], studs[5], 1, 2, 2))
	lshape.AddPaths(InkStroke(studs[6], studs[7], 1, 2, 2))
	res.AddColoredShape(lshape, color.Black, nil)

	return res
}

func MakeRoof() *g2d.Renderable {
	res := &g2d.Renderable{}

	// Flats
	poly := g2d.Polygon(roof...)
	res.AddColoredShape(g2d.NewShape(poly), brown, nil)

	// Shadows and highlights
	shad := InkStroke([]float64{roof[5][0] - 5, roof[5][1] + 5}, []float64{roof[4][0] - 5, roof[4][1]}, 3, 4, 0)
	res.AddColoredShape(g2d.NewShape(shad), brownshad, nil)
	high := InkStroke([]float64{roof[5][0] + 5, roof[5][1] + 5}, []float64{roof[0][0] + 5, roof[0][1]}, 3, 4, 0)
	res.AddColoredShape(g2d.NewShape(high), brownhigh, nil)

	// Lines
	lshape := g2d.NewShape(PolyStroke(roof, 3, 3, 2)...)
	res.AddColoredShape(lshape, color.Black, nil)

	return res
}

func MakeChimney() *g2d.Renderable {
	res := &g2d.Renderable{}

	// Flats
	poly := g2d.Polygon(chimney...)
	res.AddColoredShape(g2d.NewShape(poly), gray, nil)

	// Lines
	lshape := g2d.NewShape(PolyStroke(chimney, 3, 3, 2)...)
	res.AddColoredShape(lshape, color.Black, nil)

	return res
}

func PolyStroke(points [][]float64, mw, dw, ec float64) []*g2d.Path {
	np := len(points)
	res := make([]*g2d.Path, np)
	for i := 0; i < np-1; i++ {
		res[i] = InkStroke(points[i], points[i+1], mw, dw, ec)
	}
	res[np-1] = InkStroke(points[np-1], points[0], mw, dw, ec)
	return res
}

func InkStroke(p1, p2 []float64, mw, dw, ec float64) *g2d.Path {
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	th := math.Atan2(dy, dx)
	d := math.Hypot(dx, dy)
	mw /= 2

	rx2 := rand.Float64() * 0.1 * d
	rx1 := d - rx2
	ry1 := rand.Float64()*dw + mw
	ry2 := rand.Float64()*dw + mw
	disp := rand.Float64()

	t := rx2 / d
	omt := 1 - t
	c := []float64{p1[0]*omt + p2[0]*t, p1[1]*omt + p2[1]*t}

	d1 := rx1 * disp
	rx1 += ec
	rx2 += ec
	disp = d1 / rx1

	return g2d.IrregularEllipse(c, rx1, rx2, ry1, ry2, disp, th)
}
