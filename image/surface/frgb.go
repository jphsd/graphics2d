package surface

import "image/color"

// FRGB represents a non-premultiplied RGB tuple in floating point [0,1] such that the tuple
// can be added to, scaled and multiplied.
type FRGB struct {
	R, G, B float64
}

// NewFRGB returns a new FRGB using the supplied color.
func NewFRGB(col color.Color) *FRGB {
	r, g, b, _ := col.RGBA()
	return &FRGB{float64(r) / 0xffff, float64(g) / 0xffff, float64(b) / 0xffff}
}

// RGBA implemnts the RGBA function from the Color interface.
func (c *FRGB) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := c.R*0xffff, c.G*0xffff, c.B*0xffff
	return uint32(r), uint32(g), uint32(b), 0xffff
}

// Add returns the addition of c1 to the color (with clamping).
func (c *FRGB) Add(c1 *FRGB) *FRGB {
	r, g, b := c.R+c1.R, c.G+c1.G, c.B+c1.B
	if r < 0 {
		r = 0
	} else if r > 1 {
		r = 1
	}
	if g < 0 {
		g = 0
	} else if g > 1 {
		g = 1
	}
	if b < 0 {
		b = 0
	} else if b > 1 {
		b = 1
	}
	return &FRGB{r, g, b}
}

// Prod returns the product of c1 with the color.
func (c *FRGB) Prod(c1 *FRGB) *FRGB {
	return &FRGB{c.R * c1.R, c.G * c1.G, c.B * c1.B}
}

// Scale returns the color scaled by the value (except alpha).
func (c *FRGB) Scale(v float64) *FRGB {
	return &FRGB{c.R * v, c.G * v, c.B * v}
}

// IsBlack returns true if the color is black.
func (c *FRGB) IsBlack() bool {
	return c.R < 0.0001 && c.G < 0.0001 && c.B < 0.0001
}
