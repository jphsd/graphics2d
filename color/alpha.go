package color

// Premul provides an easy way to apply an alpha value to an existing color.
func Premul(col Color, alpha Color) RGBA64 {
	r, g, b, a := col.RGBA()
	_, _, _, aa := alpha.RGBA()
	return RGBA64{
		R: uint16((r * aa) / 0xffff),
		G: uint16((g * aa) / 0xffff),
		B: uint16((b * aa) / 0xffff),
		A: uint16((a * aa) / 0xffff)}
}
