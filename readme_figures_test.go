package graphics2d_test

import (
	"fmt"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"github.com/jphsd/texture"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/sfnt"
)

// Generates a series of outlined regular shapes.
func Example_fig01() {
	paths := []*g2d.Path{
		g2d.Line([]float64{20, 20}, []float64{130, 130}),
		g2d.RegularPolygon(3, []float64{225, 75}, 110, g2d.HalfPi),
		g2d.RegularPolygon(4, []float64{375, 75}, 110, 0),
		g2d.RegularPolygon(5, []float64{525, 75}, 75, 0),
		g2d.Circle([]float64{675, 75}, 55),
		g2d.Ellipse([]float64{825, 75}, 70, 35, g2d.HalfPi/2)}
	pen := g2d.NewPen(color.Black, 3)

	img := image.NewRGBA(900, 150, color.White)
	for _, path := range paths {
		g2d.DrawPath(img, path, pen)
	}
	image.SaveImage(img, "fig1")

	fmt.Printf("See fig1.png")
	// Output: See fig1.png
}

// Generates a series of Bezier curves of increasing order.
func Example_fig02() {
	// Create curves of order 2, 3 and 4
	quad := g2d.NewPath([]float64{175, 25})
	quad.AddStep([]float64{25, 25}, []float64{25, 175})
	cube := g2d.NewPath([]float64{375, 25})
	cube.AddStep([]float64{225, 25}, []float64{375, 175}, []float64{225, 175})
	quar := g2d.NewPath([]float64{575, 25})
	quar.AddStep([]float64{500, 25}, []float64{575, 175}, []float64{425, 100}, []float64{425, 175})

	// Draw curves
	img := image.NewRGBA(600, 200, color.White)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawPath(img, quad, pen)
	g2d.DrawPath(img, cube, pen)
	g2d.DrawPath(img, quar, pen)

	// Draw controls
	box := g2d.NewShape(g2d.RegularPolygon(4, []float64{0, 0}, 4, 0))
	cproc := g2d.CapsProc{box, box, box, false}
	paths := []*g2d.Path{quad, cube, quar}
	for _, path := range paths {
		// Control lines
		cpath := path.Process(g2d.StepsToLinesProc{true})[0]
		g2d.DrawPath(img, cpath, g2d.RedPen)

		// Control points
		shape := g2d.NewShape(cpath.Process(cproc)...)
		g2d.RenderColoredShape(img, shape, color.Black)
	}

	image.SaveImage(img, "fig2")
	fmt.Printf("See fig2.png")
	// Output: See fig2.png
}

// Generates arcs with different [ArcStyle]
func Example_fig03() {
	// Arcs
	paths := []*g2d.Path{
		// Top row
		g2d.Arc([]float64{100, 100}, 85, g2d.Pi*3/4, g2d.HalfPi, g2d.ArcOpen),
		g2d.Arc([]float64{300, 100}, 85, g2d.Pi*3/4, g2d.HalfPi, g2d.ArcPie),
		g2d.Arc([]float64{500, 100}, 85, g2d.Pi*3/4, g2d.HalfPi, g2d.ArcChord),
		// Bottom row
		g2d.Arc([]float64{100, 300}, 85, g2d.Pi*3/4, -3*g2d.HalfPi, g2d.ArcOpen),
		g2d.Arc([]float64{300, 300}, 85, g2d.Pi*3/4, -3*g2d.HalfPi, g2d.ArcPie),
		g2d.Arc([]float64{500, 300}, 85, g2d.Pi*3/4, -3*g2d.HalfPi, g2d.ArcChord),
	}
	ashape := g2d.NewShape(paths...)
	fashape := g2d.NewShape(paths[1], paths[2], paths[4], paths[5])

	// Circles
	x, y := 100.0, 100.0
	dx, dy := 200.0, 200.0
	circ := g2d.Circle([]float64{x, y}, 85)
	cshape := g2d.NewShape(circ)
	cshape.AddPaths(circ.Process(g2d.Translate(dx, 0))[0])
	cshape.AddPaths(circ.Process(g2d.Translate(2*dx, 0))[0])
	cshape.AddPaths(circ.Process(g2d.Translate(0, dy))[0])
	cshape.AddPaths(circ.Process(g2d.Translate(dx, dy))[0])
	cshape.AddPaths(circ.Process(g2d.Translate(2*dx, dy))[0])

	img := image.NewRGBA(600, 400, color.White)
	g2d.DrawShape(img, cshape, g2d.RedPen)
	g2d.RenderColoredShape(img, fashape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, ashape, pen)

	image.SaveImage(img, "fig3")
	fmt.Printf("See fig3.png")
	// Output: See fig3.png
}

// Generates a series of reentrant shapes.
func Example_fig04() {
	shape := g2d.NewShape(g2d.ReentrantPolygon([]float64{100, 100}, 90, 3, 0.5, 0))
	shape.AddPaths(g2d.ReentrantPolygon([]float64{300, 100}, 90, 4, 0.5, 0))
	shape.AddPaths(g2d.ReentrantPolygon([]float64{500, 100}, 90, 5, 0.5, 0))
	shape.AddPaths(g2d.ReentrantPolygon([]float64{700, 100}, 90, 6, 0.5, 0))

	img := image.NewRGBA(800, 200, color.White)
	g2d.RenderColoredShape(img, shape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)
	image.SaveImage(img, "fig4")
	fmt.Printf("See fig4.png")
	// Output: See fig4.png
}

// Displays the different CurveProc path processor styles.
func Example_fig05() {
	// A closed and open path
	closed := g2d.ReentrantPolygon([]float64{100, 100}, 90, 5, 0.5, 0)
	parts := closed.Parts()
	open := g2d.PartsToPath(parts[0 : len(parts)-2]...).Process(g2d.Translate(0, 200))[0]
	pshape := g2d.NewShape(closed, open)

	// Constructions
	c1 := closed.Process(g2d.Translate(200, 0))[0]
	o1 := open.Process(g2d.Translate(200, 0))[0]
	c2 := closed.Process(g2d.Translate(400, 0))[0]
	o2 := open.Process(g2d.Translate(400, 0))[0]
	c3 := closed.Process(g2d.Translate(600, 0))[0]
	o3 := open.Process(g2d.Translate(600, 0))[0]
	cshape := g2d.NewShape(c1, o1, c2, o2, c3, o3)

	// CurveProcs for each curve style
	qcproc := g2d.CurveProc{Scale: 0.5, Style: g2d.Quad}
	bcproc := g2d.CurveProc{Scale: 0.5, Style: g2d.Bezier}
	ccproc := g2d.CurveProc{Scale: 0.3, Style: g2d.CatmullRom}

	// Run the path processors
	pshape.AddPaths(c1.Process(qcproc)...)
	pshape.AddPaths(o1.Process(qcproc)...)
	pshape.AddPaths(c2.Process(bcproc)...)
	pshape.AddPaths(o2.Process(bcproc)...)
	pshape.AddPaths(c3.Process(ccproc)...)
	pshape.AddPaths(o3.Process(ccproc)...)

	img := image.NewRGBA(800, 400, color.White)
	g2d.DrawShape(img, cshape, g2d.RedPen)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, pshape, pen)
	image.SaveImage(img, "fig5")
	fmt.Printf("See fig5.png")
	// Output: See fig5.png
}

// Displays the RoundedProc path processor with different radii.
func Example_fig06() {
	// A closed and open path
	closed := g2d.ReentrantPolygon([]float64{100, 100}, 90, 5, 0.5, 0)
	parts := closed.Parts()
	open := g2d.PartsToPath(parts[0 : len(parts)-2]...).Process(g2d.Translate(0, 200))[0]
	pshape := g2d.NewShape(closed, open)

	// Constructions
	c1 := closed.Process(g2d.Translate(200, 0))[0]
	o1 := open.Process(g2d.Translate(200, 0))[0]
	c2 := closed.Process(g2d.Translate(400, 0))[0]
	o2 := open.Process(g2d.Translate(400, 0))[0]
	c3 := closed.Process(g2d.Translate(600, 0))[0]
	o3 := open.Process(g2d.Translate(600, 0))[0]
	cshape := g2d.NewShape(c1, o1, c2, o2, c3, o3)

	// CurveProcs for each curve style
	r1proc := g2d.RoundedProc{5}
	r2proc := g2d.RoundedProc{10}
	r3proc := g2d.RoundedProc{50}

	// Run the path processors
	pshape.AddPaths(c1.Process(r1proc)...)
	pshape.AddPaths(o1.Process(r1proc)...)
	pshape.AddPaths(c2.Process(r2proc)...)
	pshape.AddPaths(o2.Process(r2proc)...)
	pshape.AddPaths(c3.Process(r3proc)...)
	pshape.AddPaths(o3.Process(r3proc)...)

	img := image.NewRGBA(800, 400, color.White)
	g2d.DrawShape(img, cshape, g2d.RedPen)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, pshape, pen)
	image.SaveImage(img, "fig6")
	fmt.Printf("See fig6.png")
	// Output: See fig6.png
}

// Creates a string from a font file and displays the control points too.
func Example_fig07() {
	// Load font and create shapes
	ttf, err := sfnt.Parse(goitalic.TTF)
	if err != nil {
		panic(err)
	}
	str := "G2D"
	shape, _, err := g2d.StringToShape(ttf, str)
	if err != nil {
		panic(err)
	}

	// Figure bounding box and scaling transform
	bb := shape.BoundingBox()
	xfm := g2d.ScaleAndInset(500, 300, 20, 20, false, bb)
	shape = shape.ProcessPaths(xfm)

	// Render string
	img := image.NewRGBA(500, 300, color.White)
	g2d.RenderColoredShape(img, shape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)

	// Render construction
	box := g2d.NewShape(g2d.RegularPolygon(4, []float64{0, 0}, 4, 0))
	cproc := g2d.CapsProc{box, box, box, false}
	for _, path := range shape.Paths() {
		// Control lines
		cpath := path.Process(g2d.StepsToLinesProc{true})[0]
		g2d.DrawPath(img, cpath, g2d.RedPen)

		// Control points
		shape := g2d.NewShape(cpath.Process(cproc)...)
		g2d.RenderColoredShape(img, shape, color.Black)
	}
	image.SaveImage(img, "fig7")

	fmt.Printf("See fig7.png")
	// Output: See fig7.png
}

// Generates a series of regular shapes with dashed outlines.
func Example_fig08() {
	paths := []*g2d.Path{
		g2d.Line([]float64{20, 20}, []float64{130, 130}),
		g2d.RegularPolygon(3, []float64{225, 75}, 110, g2d.HalfPi),
		g2d.RegularPolygon(4, []float64{375, 75}, 110, 0),
		g2d.RegularPolygon(5, []float64{525, 75}, 75, 0),
		g2d.Circle([]float64{675, 75}, 55),
		g2d.Ellipse([]float64{825, 75}, 70, 35, g2d.HalfPi/2)}

	// Path processors
	proc1 := g2d.NewDashProc([]float64{4, 2}, 0)
	proc2 := g2d.NewDashProc([]float64{8, 2, 2, 2}, 0)
	proc3 := g2d.NewDashProc([]float64{10, 4}, 0)
	head := g2d.NewShape(g2d.PolyLine([]float64{-2, 2}, []float64{0, 0}, []float64{-2, -2}))
	cproc := g2d.CapsProc{nil, head, nil, true}

	shape := &g2d.Shape{}
	xfm := g2d.Translate(0, 150)
	for _, path := range paths {
		shape.AddPaths(path.Process(proc1)...)
		path = path.Process(xfm)[0]
		shape.AddPaths(path.Process(proc2)...)
		path = path.Process(xfm)[0]
		lshape := g2d.NewShape(path.Process(proc3)...)
		shape.AddShapes(lshape)
		// Add arrow heads to dashes from proc3
		shape.AddShapes(lshape.ProcessPaths(cproc))
	}

	img := image.NewRGBA(900, 450, color.White)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)
	image.SaveImage(img, "fig8")

	fmt.Printf("See fig8.png")
	// Output: See fig8.png
}

// Generates a series of path traces using different join functions.
func Example_fig09() {
	path := g2d.PolyLine(
		[]float64{20, 50},
		[]float64{120, 150},
		[]float64{220, 50},
		[]float64{320, 150},
		[]float64{420, 50},
		[]float64{520, 150})

	proc1 := g2d.NewTraceProc(20)
	proc2 := g2d.NewTraceProc(20)
	proc2.JoinFunc = g2d.JoinRound
	proc3 := g2d.NewTraceProc(20)
	proc3.JoinFunc = g2d.NewMiterJoin().JoinMiter

	img := image.NewRGBA(560, 600, color.White)

	g2d.DrawPath(img, path, g2d.RedPen)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawPath(img, path.Process(proc1)[0], pen)

	path = path.Process(g2d.Translate(0, 200))[0]
	g2d.DrawPath(img, path, g2d.RedPen)
	g2d.DrawPath(img, path.Process(proc2)[0], pen)

	path = path.Process(g2d.Translate(0, 200))[0]
	g2d.DrawPath(img, path, g2d.RedPen)
	g2d.DrawPath(img, path.Process(proc3)[0], pen)

	image.SaveImage(img, "fig9")
	fmt.Printf("See fig9.png")
	// Output: See fig9.png
}

// Generates a variable width trace of a path.
func Example_fig10() {
	// Line, MPD it, round it - a wriggle
	path := g2d.Line([]float64{30, 150}, []float64{530, 150})
	path = path.Process(&g2d.MPDProc{.3, 3, 0.5, false})[0]
	path = path.Process(&g2d.RoundedProc{1000})[0]

	proc := &g2d.VWTraceProc{
		Width:   -20,
		Flatten: g2d.RenderFlatten,
	}
	proc.Func = func(t, w float64) float64 {
		return (1-t)*proc.Width + 1
	}

	img := image.NewRGBA(560, 300, color.White)

	g2d.DrawPath(img, path, g2d.RedPen)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawPath(img, path.Process(proc)[0], pen)

	image.SaveImage(img, "fig10")
	fmt.Printf("See fig10.png")
	// Output: See fig10.png
}

// Demonstrates the variety of stroke end caps.
func Example_fig11() {
	img := image.NewRGBA(900, 300, color.White)

	// Butt
	path := g2d.Line([]float64{50, 75}, []float64{250, 75})
	proc := g2d.NewStrokeProc(40)
	shape := g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	pen := g2d.NewPen(color.Black, 3)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	// Square
	path = g2d.Line([]float64{350, 75}, []float64{550, 75})
	proc.CapStartFunc = g2d.CapSquare
	proc.CapEndFunc = g2d.CapSquare
	shape = g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	// Rounded Square
	path = g2d.Line([]float64{650, 75}, []float64{850, 75})
	rsc := g2d.RSCap{0.5}
	proc.CapStartFunc = rsc.CapRoundedSquare
	proc.CapEndFunc = rsc.CapRoundedSquare
	shape = g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	// Round
	path = g2d.Line([]float64{50, 225}, []float64{250, 225})
	proc.CapStartFunc = g2d.CapInvRound
	proc.CapEndFunc = g2d.CapRound
	shape = g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	// Oval
	path = g2d.Line([]float64{350, 225}, []float64{550, 225})
	oc := g2d.OvalCap{2, 0}
	proc.CapStartFunc = oc.CapInvOval
	proc.CapEndFunc = oc.CapOval
	shape = g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	// Point
	path = g2d.Line([]float64{650, 225}, []float64{850, 225})
	proc.CapStartFunc = g2d.CapInvPoint
	proc.CapEndFunc = g2d.CapPoint
	shape = g2d.NewShape(path.Process(proc)...)
	g2d.RenderColoredShape(img, shape, color.Green)
	g2d.DrawShape(img, shape, pen)
	g2d.DrawPath(img, path, g2d.RedPen)

	image.SaveImage(img, "fig11")
	fmt.Printf("See fig11.png")
	// Output: See fig11.png
}

// Demonstrates the different gradient types: linear, radial, elliptical and conic.
func Example_fig12() {
	width, height := 1200, 300

	shape1 := g2d.NewShape(g2d.Circle([]float64{150, 150}, 100))
	filler1 := texture.NewLinearRGBA(width, height,
		[]float64{50, 150}, []float64{250, 150},
		color.Green, color.Red,
		texture.NewNLLinear(),
		false, false)
	shape2 := g2d.NewShape(g2d.Circle([]float64{450, 150}, 100))
	filler2 := texture.NewRadialRGBA(width, height,
		[]float64{450, 150}, 100,
		color.Red, color.Green,
		texture.NewNLLinear(),
		false, false)
	shape3 := g2d.NewShape(g2d.Circle([]float64{750, 150}, 100))
	filler3 := texture.NewEllipticalRGBA(width, height,
		[]float64{750, 150},
		100, 50, g2d.Pi/4,
		color.Red, color.Green,
		texture.NewNLLinear(),
		true, false)
	shape4 := g2d.NewShape(g2d.Circle([]float64{1050, 150}, 100))
	filler4 := texture.NewConicRGBA(width, height,
		[]float64{1050, 150}, 0,
		color.Green, color.Red,
		texture.NewNLLinear())

	img := image.NewRGBA(width, height, color.White)
	pen := g2d.NewPen(color.Black, 3)
	g2d.RenderShape(img, shape1, filler1)
	g2d.DrawShape(img, shape1, pen)
	g2d.RenderShape(img, shape2, filler2)
	g2d.DrawShape(img, shape2, pen)
	g2d.RenderShape(img, shape3, filler3)
	g2d.DrawShape(img, shape3, pen)
	g2d.RenderShape(img, shape4, filler4)
	g2d.DrawShape(img, shape4, pen)
	image.SaveImage(img, "fig12")
	fmt.Printf("See fig12.png")
	// Output: See fig12.png
}
