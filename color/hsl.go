package color

import (
	"math"
)

// HSL describes a color in Hue Saturation Lightness space. All values are in range [0,1].
type HSL struct {
	H, S, L, A float64
}

// HSL conversions (see https://www.w3.org/TR/css-color-3/#hsl-color)

// RGBA implements the RGBA function from the color.Color interface.
func (c HSL) RGBA() (uint32, uint32, uint32, uint32) {
	// Sanitize H, S, L and A - allow H to wrap
	c.H = math.Mod(c.H, 1)
	if c.H < 0 {
		c.H += 1
	}
	c.S = max(c.S, 0)
	c.S = min(c.S, 1)
	c.L = max(c.L, 0)
	c.L = min(c.L, 1)
	c.A = max(c.A, 0)
	c.A = min(c.A, 1)

	// Convert HSL to RGB (all in [0,1])
	// From https://en.wikipedia.org/wiki/HSL_and_HSV
	l := math.Abs(2*c.L - 1)
	chroma := (1 - l) * c.S
	hue := c.H * 6
	x := chroma * (1 - math.Abs(math.Mod(hue, 2)-1))

	var r, g, b float64
	if hue < 1 {
		r = chroma
		g = x
	} else if hue < 2 {
		r = x
		g = chroma
	} else if hue < 3 {
		g = chroma
		b = x
	} else if hue < 4 {
		g = x
		b = chroma
	} else if hue < 5 {
		r = x
		b = chroma
	} else {
		r = chroma
		b = x
	}

	// Scale by A and convert to uint32
	return uint32(r * c.A * 0xffff), uint32(g * c.A * 0xffff), uint32(b * c.A * 0xffff), uint32(c.A * 0xffff)
}

// HSLModel standard HSL color type with all values in range [0,1]
var HSLModel Model = ModelFunc(hslModel)

func hslModel(col Color) Color {
	hsl, ok := col.(HSL)
	if !ok {
		hsl = NewHSL(col)
	}
	return hsl
}

// NewHSL returns the color as an HSL triplet.
func NewHSL(col Color) HSL {
	if hsl, ok := col.(HSL); ok {
		return HSL{hsl.H, hsl.S, hsl.L, hsl.A}
	}
	ir, ig, ib, ia := col.RGBA()
	if ia == 0 {
		return HSL{0, 0, 0, 0}
	}

	// Convert to [0,1]
	r, g, b, a := float64(ir)/0xffff, float64(ig)/0xffff, float64(ib)/0xffff, float64(ia)/0xffff
	r /= a
	g /= a
	b /= a

	min := min(min(r, g), b)
	max := max(max(r, g), b)

	l := (max + min) / 2

	var s, h float64
	if equals(min, max) {
		s = 0
		h = 0
	} else {
		d := max - min
		if l < 0.5 {
			s = d / (max + min)
		} else {
			s = d / (2.0 - max - min)
		}

		if max == r {
			h = (g - b) / d
		} else if max == g {
			h = 2 + (b-r)/d
		} else {
			h = 4 + (r-g)/d
		}

		h /= 6

		if h < 0 {
			h += 1
		}
	}

	return HSL{h, s, l, a}
}

// HSVToHSL is a convenience function that maps HSV to HSL (all in [0,1])
func HSVToHSL(h, sv, v float64) HSL {
	// Sanitize H, S, V and A
	h = math.Mod(h, 1)
	if h < 0 {
		h += 1
	}
	sv = max(sv, 0)
	sv = min(sv, 1)
	v = max(v, 0)
	v = min(v, 1)

	// Convert
	l := v * (1 - sv/2)
	var sl float64
	if l > 0 && l < 1 {
		sl = (v - l) / min(l, 1-l)
	}
	return HSL{h, sl, l, 1}
}

// Complement returns the color's complement.
func Complement(col Color) HSL {
	hsl := NewHSL(col)
	hsl.H += 0.5
	if hsl.H > 1 {
		hsl.H -= 1
	}
	return hsl
}

// Monochrome returns the color's monochrome palette (excluding black and white).
// Note the palette may not contain the original color since the values are equally
// spaced over L.
func Monochrome(col Color, n int) []HSL {
	if n < 2 {
		return []HSL{NewHSL(col)}
	}
	res := make([]HSL, n)
	dl := 1.0 / float64(n-1)
	l := dl
	for i := range n {
		hsl := NewHSL(col)
		hsl.L = l
		res[i] = hsl
		l += dl
	}
	return res
}

// Analogous returns the color's analogous colors.
func Analogous(col Color) []HSL {
	a1, a2 := NewHSL(col), NewHSL(col)
	d := 1 / float64(12)
	a1.H += d
	if a1.H > 1 {
		a1.H -= 1
	}
	a2.H -= d
	if a2.H < 0 {
		a2.H += 1
	}
	return []HSL{a1, a2}
}

// Triad returns the color's other two triadics.
func Triad(col Color) []HSL {
	a1, a2 := NewHSL(col), NewHSL(col)
	d := 1 / float64(3)
	a1.H += d
	if a1.H > 1 {
		a1.H -= 1
	}
	a2.H -= d
	if a2.H < 0 {
		a2.H += 1
	}
	return []HSL{a1, a2}
}

// Tetrad returns the color's other three tetradics.
func Tetrad(col Color) []HSL {
	a1, a2, a3 := NewHSL(col), NewHSL(col), NewHSL(col)
	d := 0.25
	a1.H += d
	if a1.H > 1 {
		a1.H -= 1
	}
	a2.H -= d
	if a2.H < 0 {
		a2.H += 1
	}
	a3.H += 0.5
	if a3.H > 1 {
		a3.H -= 1
	}
	return []HSL{a1, a3, a2}
}

// Warmer returns the color shifted towards red.
func Warmer(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.H, 0) {
		return hsl
	}
	if hsl.H < 0.5 {
		if hsl.H < 0.1 {
			hsl.H = 0
			return hsl
		}
		hsl.H -= 0.1
		return hsl
	}
	if hsl.H > 0.9 {
		hsl.H = 0
		return hsl
	}
	hsl.H += 0.1
	return hsl
}

// Cooler returns the color shifted toward cyan.
func Cooler(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.H, 0.5) {
		return hsl
	}
	if hsl.H < 0.5 {
		if hsl.H > 0.4 {
			hsl.H = 0.5
			return hsl
		}
		hsl.H += 0.1
		return hsl
	}
	if hsl.H < 0.6 {
		hsl.H = 0.5
		return hsl
	}
	hsl.H -= 0.1
	return hsl
}

// Tint returns the color shifted towards white.
func Tint(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.L, 1) {
		return hsl
	}
	if hsl.L > 0.9 {
		hsl.L = 1
		return hsl
	}
	hsl.L += 0.1
	return hsl
}

// Shade returns the color shifted towards black.
func Shade(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.L, 0) {
		return hsl
	}
	if hsl.L < 0.1 {
		hsl.L = 0
		return hsl
	}
	hsl.L -= 0.1
	return hsl
}

// Boost returns the color shifted away from gray.
func Boost(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.S, 1) {
		return hsl
	}
	if hsl.S > 0.9 {
		hsl.S = 1
		return hsl
	}
	hsl.S += 0.1
	return hsl
}

// Tone returns the color shifted towards gray.
func Tone(col Color) HSL {
	hsl := NewHSL(col)
	if equals(hsl.S, 0) {
		return hsl
	}
	if hsl.S < 0.1 {
		hsl.S = 0
		return hsl
	}
	hsl.S -= 0.1
	return hsl
}

// Compound returns the colors analogous to the color's complement.
func Compound(col Color) []HSL {
	return Analogous(Complement(col))
}

// HuePair returns a pair of colors d away on either side of the color. d is in range [0,1]
func HuePair(col Color, d float64) []HSL {
	a1, a2 := NewHSL(col), NewHSL(col)
	if d < 0 {
		d = 0
	} else if d > 1 {
		d = 1
	}
	a1.H += d
	if a1.H > 1 {
		a1.H -= 1
	}
	a2.H -= d
	if a2.H < 0 {
		a2.H += 1
	}
	return []HSL{a1, a2}
}

// SatPair returns a pair of colors d away on either side of the color. d is in range [0,1]
func SatPair(col Color, d float64) []HSL {
	a1, a2 := NewHSL(col), NewHSL(col)
	if d < 0 {
		d = 0
	} else if d > 1 {
		d = 1
	}
	a1.S += d
	if a1.S > 1 {
		a1.S = 1
	}
	a2.S -= d
	if a2.S < 0 {
		a2.S = 0
	}
	return []HSL{a1, a2}
}

// LightPair returns a pair of colors d away on either side of the color. d is in range [0,1]
func LightPair(col Color, d float64) []HSL {
	a1, a2 := NewHSL(col), NewHSL(col)
	if d < 0 {
		d = 0
	} else if d > 1 {
		d = 1
	}
	a1.L += d
	if a1.L > 1 {
		a1.L = 1
	}
	a2.L -= d
	if a2.L < 0 {
		a2.L = 0
	}
	return []HSL{a1, a2}
}

func equals(v1, v2 float64) bool {
	d := v1 - v2
	if d < 0 {
		d = -d
	}

	return d < 0.000001
}
