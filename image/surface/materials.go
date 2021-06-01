package surface

import "image/color"

// Material provides the At function to determine the emissive light, various reflectances and shininess at a location.
// Reflectances are ordered as ambient, diffuse and specular.
type Material interface {
	At(x, y int) (*FRGB, *FRGB, *FRGB, *FRGB, float64)
}

type defaultMaterial struct {
	Ambient, Diffuse *FRGB
}

// DefaultMaterial describes a material with 0 emissivity, full white ambient and directional, and no specular
// components.
var DefaultMaterial = &defaultMaterial{NewFRGB(color.White), NewFRGB(color.White)}

// At implements the At function in the Material interface.
func (d *defaultMaterial) At(x, y int) (*FRGB, *FRGB, *FRGB, *FRGB, float64) {
	return nil, d.Ambient, d.Diffuse, nil, 0
}
