//go:build ignore

package main

import (
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/image"
	stdimg "image"
	"image/color"
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

// Demonstrate use of IrregularEllipses and fills to draw house
func main() {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)

	// Locate and scale
	xfm := g2d.NewAff3()
	xfm.Translate(250, 250)
	xfm.Scale(2, 2)

	// Draw house
	DrawHouse(img, xfm)

	// Save image
	image.SaveImage(img, "house")
}

func DrawHouse(img *stdimg.RGBA, xfm *g2d.Aff3) {
	// Order is important - furthest from eye obejcts first
	DrawChimney(img, xfm)
	DrawFirstFloor(img, xfm)
	DrawSecondFloor(img, xfm)
	DrawRoof(img, xfm)
}

func DrawFirstFloor(img *stdimg.RGBA, xfm *g2d.Aff3) {
	// Rendering offset
	ofirst := xfm.Apply(first...)
	odoor := xfm.Apply(door...)
	owindow := xfm.Apply(window...)

	// Flats
	poly := g2d.Polygon(ofirst...)
	g2d.RenderColoredShape(img, g2d.NewShape(poly), gray)
	poly = g2d.Polygon(odoor...)
	g2d.RenderColoredShape(img, g2d.NewShape(poly), brown)

	// Shadows and highlights
	shad := InkStroke([]float64{ofirst[2][0] - 10, ofirst[2][1]}, []float64{ofirst[1][0] - 10, ofirst[1][1]}, 7, 5, 0)
	g2d.RenderColoredShape(img, g2d.NewShape(shad), grayshad)
	high := InkStroke([]float64{ofirst[3][0] + 10, ofirst[3][1]}, []float64{ofirst[0][0] + 10, ofirst[0][1]}, 7, 5, 0)
	g2d.RenderColoredShape(img, g2d.NewShape(high), grayhigh)

	// Lines
	lshape := g2d.NewShape(PolyStroke(ofirst, 3, 3, 2)...)
	lshape.AddPaths(PolyStroke(owindow, 1, 2, 2)...)
	lshape.AddPaths(PolyStroke(odoor, 1, 2, 2)...)
	g2d.RenderColoredShape(img, lshape, color.Black)
}

func DrawSecondFloor(img *stdimg.RGBA, xfm *g2d.Aff3) {
	// Rendering offset
	osecond := xfm.Apply(second...)
	ostringer := xfm.Apply(stringer...)
	ostuds := xfm.Apply(studs...)

	// Flats
	poly := g2d.Polygon(osecond...)
	g2d.RenderColoredShape(img, g2d.NewShape(poly), brown)

	// Shadows
	shad1 := InkStroke([]float64{osecond[2][0] + 5, osecond[2][1]}, []float64{osecond[0][0] + 5, osecond[0][1]}, 3, 4, 0)
	shad2 := InkStroke([]float64{osecond[2][0] - 5, osecond[2][1]}, []float64{osecond[1][0] - 5, osecond[1][1]}, 3, 4, 0)
	g2d.RenderColoredShape(img, g2d.NewShape(shad1, shad2), brownshad)

	// Lines
	lshape := g2d.NewShape(PolyStroke(osecond, 3, 3, 2)...)
	lshape.AddPaths(InkStroke(ostringer[0], ostringer[1], 1, 2, 2))
	lshape.AddPaths(InkStroke(ostuds[0], ostuds[1], 1, 2, 2))
	lshape.AddPaths(InkStroke(ostuds[2], ostuds[3], 1, 2, 2))
	lshape.AddPaths(InkStroke(ostuds[4], ostuds[5], 1, 2, 2))
	lshape.AddPaths(InkStroke(ostuds[6], ostuds[7], 1, 2, 2))
	g2d.RenderColoredShape(img, lshape, color.Black)
}

func DrawRoof(img *stdimg.RGBA, xfm *g2d.Aff3) {
	// Rendering offset
	oroof := xfm.Apply(roof...)

	// Flats
	poly := g2d.Polygon(oroof...)
	g2d.RenderColoredShape(img, g2d.NewShape(poly), brown)

	// Shadows and highlights
	shad := InkStroke([]float64{oroof[5][0] - 5, oroof[5][1] + 5}, []float64{oroof[4][0] - 5, oroof[4][1]}, 3, 4, 0)
	g2d.RenderColoredShape(img, g2d.NewShape(shad), brownshad)
	high := InkStroke([]float64{oroof[5][0] + 5, oroof[5][1] + 5}, []float64{oroof[0][0] + 5, oroof[0][1]}, 3, 4, 0)
	g2d.RenderColoredShape(img, g2d.NewShape(high), brownhigh)

	// Lines
	lshape := g2d.NewShape(PolyStroke(oroof, 3, 3, 2)...)
	g2d.RenderColoredShape(img, lshape, color.Black)
}

func DrawChimney(img *stdimg.RGBA, xfm *g2d.Aff3) {
	// Rendering offset
	ochimney := xfm.Apply(chimney...)

	// Flats
	poly := g2d.Polygon(ochimney...)
	g2d.RenderColoredShape(img, g2d.NewShape(poly), gray)

	// Lines
	lshape := g2d.NewShape(PolyStroke(ochimney, 3, 3, 2)...)
	g2d.RenderColoredShape(img, lshape, color.Black)
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
	d := math.Sqrt(dx*dx + dy*dy)
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
