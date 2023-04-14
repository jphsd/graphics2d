package image

import (
	g2dcol "github.com/jphsd/graphics2d/color"
	"image"
	"image/color"
	"sort"
)

// Colorizer uses the gray or red channel of an image to create an image tinted with a lerp'd gradient
// between two or more colors.
type Colorizer struct {
	Img  image.Image
	Rect image.Rectangle
	Lut  []color.RGBA
	G16  bool
}

type cstop struct {
	s int
	c color.RGBA
}

// NewColorizer creates a new Colorizer and creates the internal lut. c1 and c2 are the colors at the start and end
// of the colorizer. The stops (in range [1,254] and colors values determine any intermediate points and can be nil.
// The post flag turns off the lerping between colors and yields a posterized image.
func NewColorizer(img image.Image, c1, c2 color.Color, stops []int, colors []color.Color, post bool) *Colorizer {
	lut := make([]color.RGBA, 256)

	sc := len(colors)
	if sc > 254 {
		sc = 254
	}
	if stops != nil && len(stops) < sc {
		sc = len(stops)
	}

	csl := make([]cstop, sc+2)
	c, _ := color.RGBAModel.Convert(c1).(color.RGBA)
	csl[0] = cstop{0, c}
	c, _ = color.RGBAModel.Convert(c2).(color.RGBA)
	csl[sc+1] = cstop{255, c}

	if stops != nil {
		for i := 0; i < sc; i++ {
			c, _ = color.RGBAModel.Convert(colors[i]).(color.RGBA)
			csl[i+1] = cstop{stops[i], c}
		}
		// Check stops are ascending
		sort.Slice(csl, func(i, j int) bool {
			return csl[i].s < csl[j].s
		})
	} else {
		// Figure out stop values
		bx := 256.0 / float64(sc+1)
		b := bx
		for i := 0; i < sc; i++ {
			c, _ = color.RGBAModel.Convert(colors[i]).(color.RGBA)
			csl[i+1] = cstop{int(b), c}
			b += bx
		}
	}

	// Build lut - simple lerp between c1, stops and c2, unless post set
	if post {
		ci := 0
		for i := 0; i < 256; i++ {
			lut[i] = csl[ci].c
			if i != 255 && i+1 == csl[ci+1].s {
				ci++
			}
		}
	} else {
		ci := 0
		ls := 0
		for i := 0; i < 256; i++ {
			if i == csl[ci].s {
				lut[i] = csl[ci].c
				if i != 255 {
					ls = csl[ci+1].s - csl[ci].s
				}
			} else {
				ds := i - csl[ci].s
				t := float64(ds) / float64(ls)
				lut[i] = g2dcol.ColorRGBALerp(t, csl[ci].c, csl[ci+1].c)
			}
			if i != 255 && i+1 == csl[ci+1].s {
				ci++
			}
		}
	}

	cm := img.ColorModel()
	g16 := false
	if cm == color.Gray16Model {
		g16 = true
	}
	return &Colorizer{img, img.Bounds(), lut, g16}
}

// ColorModel implements the ColorModel function in the Image interface.
func (c *Colorizer) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the Bounds function in the Image interface.
func (c *Colorizer) Bounds() image.Rectangle {
	return c.Rect
}

// At implements the At function in the Image interface.
func (c *Colorizer) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(c.Rect)) {
		return color.RGBA{}
	}
	var v uint32
	cc := c.Img.At(x, y)
	if c.G16 {
		gc := cc.(color.Gray16)
		v = uint32(gc.Y)
	} else {
		v, _, _, _ = cc.RGBA()
	}
	v = (v & 0xff00) >> 8
	return c.Lut[v]
}
