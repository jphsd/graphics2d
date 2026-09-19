package color

// Premul returns an alpha premultiplied value of col as RGBA
func Premul(col Color, alpha Color) RGBA {
	r, g, b, a := col.RGBA()
	_, _, _, aa := alpha.RGBA()
	if a == aa {
		return RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	if a != 0xffff {
		// Convert to non-premultiplied values
		r *= 0xffff
		g *= 0xffff
		b *= 0xffff
		r /= a
		g /= a
		b /= a
	}
	r *= aa
	r >>= 24
	g *= aa
	g >>= 24
	b *= aa
	b >>= 24
	aa >>= 8
	return RGBA{uint8(r), uint8(g), uint8(b), uint8(aa)}
}
