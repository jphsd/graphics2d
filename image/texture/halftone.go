package texture

import (
	g2d "github.com/jphsd/graphics2d"
	"image/color"
	"image/draw"
	"math"
)

// Halftone renders colored halftone dots into the destination image. Line separation, rotation
// and percentage fill with a point offset (from {0, 0}) control the dot locations.
func Halftone(dst draw.Image, c color.Color, perc, sep, rot float64, offs []float64) {
	r := dst.Bounds()

	if perc <= 0 {
		return
	} else if perc > 1 {
		perc = 1
	}
	rad := perc2radius(perc)
	// Circle path template centered on 0, 0
	circp := g2d.Circle([]float64{0, 0}, rad*sep)

	// Calculate grid of centers
	l := r.Max.X
	if l < r.Max.Y {
		l = r.Max.Y
	}
	lf := 1.5 * float64(l) // [-lf,lf]^2
	n := int(math.Ceil(2*lf/sep)) + 1
	points := make([][]float64, 1, n*n)
	points[0] = []float64{0, 0}
	// Axis
	for i := sep; i < lf; i += sep {
		points = append(points, []float64{0, i}, []float64{0, -i}, []float64{i, 0}, []float64{-i, 0})
	}
	for y := sep; y < lf; y += sep {
		for x := sep; x < lf; x += sep {
			points = append(points, []float64{x, y}, []float64{x, -y}, []float64{-x, y}, []float64{-x, -y})
		}
	}

	// Handle offs and rot
	for offs[0] < 0 {
		offs[0] += sep
	}
	for offs[0] > sep {
		offs[0] -= sep
	}
	for offs[1] < 0 {
		offs[1] += sep
	}
	for offs[1] > sep {
		offs[1] -= sep
	}
	for rot < 0 {
		rot += math.Pi
	}
	for rot > math.Pi {
		rot -= math.Pi
	}

	xfm := g2d.NewAff3()
	xfm.Translate(offs[0], offs[1])
	xfm.Rotate(rot)
	points = xfm.Apply(points...)

	// Construct shape with dots in it at each point location
	shape := &g2d.Shape{}
	for _, pt := range points {
		xfm := g2d.NewAff3()
		xfm.Translate(pt[0], pt[1])
		shape.AddPaths(circp.Transform(xfm))
	}

	// Write to image
	g2d.RenderColoredShape(dst, shape, c)
}

// Precalculated radius to percentage for once dots start to overlap
var lut = [][]float64{
	{0.5, 0.785398163},
	{0.51, 0.811757717},
	{0.52, 0.834192027},
	{0.53, 0.854184781},
	{0.54, 0.872243896},
	{0.55, 0.888652751},
	{0.56, 0.903596183},
	{0.57, 0.917205474},
	{0.58, 0.929579102},
	{0.59, 0.940793811},
	{0.6, 0.950911131},
	{0.61, 0.959981487},
	{0.62, 0.968046934},
	{0.63, 0.975143051},
	{0.64, 0.981300303},
	{0.65, 0.986545039},
	{0.66, 0.990900248},
	{0.67, 0.994386142},
	{0.68, 0.997020609},
	{0.69, 0.998819573},
	{0.7, 0.999797286},
	{0.71, 1},
}

func perc2radius(perc float64) float64 {
	// For r <= 0.5 there's no overlap so simple calculation
	if perc < lut[0][1] {
		return math.Sqrt(perc / math.Pi)
	}
	var i int
	for i = 0; lut[i][1] < perc; i++ {
	}
	// Linear interp
	r0, r1 := lut[i-1][0], lut[i][0]
	p0, p1 := lut[i-1][1], lut[i][1]
	t := (perc - p0) / (p1 - p0)
	return r0*(1-t) + r1*t
}
