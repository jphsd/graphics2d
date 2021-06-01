package surface

import (
	"image"
	"image/draw"
	"math"
)

// Surface collects the ambient light, lights, a material, and normal map required to describe
// an area. If the normal map is nil then the standard normal is use {0, 0, 1}
type Surface struct {
	Ambient Light
	Lights  []Light
	Mat     Material
	Normals NormalMap
}

var blinn = false

// RenderSurface renders the supplied surface (offset by sp) into the indicated rectangle of the destination
// image.
func RenderSurface(dst draw.Image, rect image.Rectangle, surf *Surface, sp image.Point) {
	// For any point, the color rendered is the sum of the emissive, ambient and the diffuse/specular
	// contributions from all of the lights.

	material := surf.Mat
	normals := surf.Normals
	if normals == nil {
		normals = &DefaultNM{}
	}
	ambient := surf.Ambient
	view := []float64{0, 0, 1}

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		oy := y + sp.Y
		for x := rect.Min.X; x < rect.Max.X; x++ {
			ox := x + sp.X
			emm, amb, diff, spec, shine := material.At(ox, oy) // Emissive
			lemm := &FRGB{}
			if emm != nil {
				lemm = emm
			}
			acol, _, _, _ := ambient.At(ox, oy)
			lamb := amb.Prod(acol) // Ambient
			normal := normals.At(ox, oy)
			if diff == nil {
				continue
			}
			cdiff, cspec := &FRGB{}, &FRGB{}
			for _, light := range surf.Lights {
				lcol, dir, dist, pow := light.At(ox, oy)
				if lcol.IsBlack() {
					continue
				}
				lambert := dot(dir, normal)
				if lambert < 0 {
					continue
				}
				if dist > 0 {
					lcol = lcol.Scale(pow / (dist * dist))
				}
				cdiff = cdiff.Add(lcol.Prod(diff.Scale(lambert))) // Diffuse
				if spec != nil {
					if blinn {
						// Blinn-Phong
						half := norm([]float64{dir[0] + view[0], dir[1] + view[1], dir[2] + view[2]})
						dp := dot(half, normal)
						if dp > 0 {
							phong := math.Pow(dp, shine*4)
							cspec = cspec.Add(lcol.Prod(spec.Scale(phong))) // Specular
						}
					} else {
						// Phong
						dp := dot(reflect(dir, normal), view)
						if dp > 0 {
							phong := math.Pow(dp, shine)
							cspec = cspec.Add(lcol.Prod(spec.Scale(phong))) // Specular
						}
					}
				}
			}
			col := lemm
			col = col.Add(lamb)
			col = col.Add(cdiff)
			col = col.Add(cspec)
			dst.Set(x, y, col)
		}
	}
}

func dot(a, b []float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func norm(v []float64) []float64 {
	sum := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	return []float64{v[0] / sum, v[1] / sum, v[2] / sum}
}

// Return v reflected in n
func reflect(v, n []float64) []float64 {
	s := dot(v, n)
	s *= 2
	return []float64{s*n[0] - v[0], s*n[1] - v[1], s*n[2] - v[2]}
}

func cross(a, b []float64) []float64 {
	return []float64{
		a[1]*b[2] - a[2]*b[1],
		-a[0]*b[2] + a[2]*b[0],
		a[0]*b[1] - a[1]*b[0]}
}
