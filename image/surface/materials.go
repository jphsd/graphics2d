package surface

import (
	gc "github.com/jphsd/graphics2d/color"
	"image/color"
)

// Material provides the At function to determine the emissive light, various reflectances and shininess at a location.
// Reflectances are ordered as ambient, diffuse and specular.
type Material interface {
	At(x, y int) (*gc.FRGBA, *gc.FRGBA, *gc.FRGBA, *gc.FRGBA, float64)
}

type defaultMaterial struct {
	Ambient, Diffuse *gc.FRGBA
}

// DefaultMaterial describes a material with 0 emissivity, full white ambient and directional, and no specular
// components.
var DefaultMaterial = &defaultMaterial{gc.NewFRGBA(color.White), gc.NewFRGBA(color.White)}

// At implements the At function in the Material interface.
func (d *defaultMaterial) At(x, y int) (*gc.FRGBA, *gc.FRGBA, *gc.FRGBA, *gc.FRGBA, float64) {
	return nil, d.Ambient, d.Diffuse, nil, 0
}
