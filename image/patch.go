package image

import (
	"fmt"
	"image"
	"image/color"
)

// Patch is an infinite image covered with a patch of colors.
type Patch struct {
	Colors [][]color.RGBA
	Width  int
	Height int
	OffsX  int
	OffsY  int
}

// NewPatch creates a new iamge with the supplied patch.
func NewPatch(colors [][]color.Color) (*Patch, error) {
	h := len(colors)
	w := len(colors[0])
	rgba := make([][]color.RGBA, h)
	for i := 0; i < h; i++ {
		if len(colors[i]) != w {
			return nil, fmt.Errorf("row %d has different length %d vs %d", i, len(colors[i]), w)
		}
		rgba[i] = make([]color.RGBA, w)
		for j := 0; j < w; j++ {
			rgba[i][j], _ = color.RGBAModel.Convert(colors[i][j]).(color.RGBA)
		}
	}
	return &Patch{rgba, w, h, 0, 0}, nil
}

// ColorModel implements the ColorModel function in the Image interface.
func (p *Patch) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the Bounds function in the Image interface.
func (p *Patch) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{-1e9, -1e9}, image.Point{1e9, 1e9}}
}

// At implements the At function in the Image interface.
func (p *Patch) At(x, y int) color.Color {
	i := 0
	x += p.OffsX
	y += p.OffsY
	if x < 0 {
		x = -x
		i = x % p.Width
		if i != 0 {
			i = p.Width - i
		}
	} else {
		i = x % p.Width
	}
	j := 0
	if y < 0 {
		y = -y
		j = y % p.Height
		if j != 0 {
			j = p.Height - j
		}
	} else {
		j = y % p.Height
	}
	return p.Colors[j][i]
}
