package image

import (
	"image"
	//"image/color"
	"github.com/jphsd/graphics2d/color"
	"image/draw"
)

// Tile is an infinite image covered with a tile.
type Tile struct {
	TileImg *RGBA
	Width   int
	Height  int
	OffsX   int
	OffsY   int
	StagX   int
	StagY   int
}

// NewTile creates a new image with the supplied image tile.
func NewTile(img Image) *Tile {
	rect := img.Bounds()
	w, h := rect.Dx(), rect.Dy()
	tile := image.NewRGBA(image.Rectangle{image.Point{}, rect.Size()})
	draw.Draw(tile, tile.Bounds(), img, rect.Min, draw.Src)
	return &Tile{tile, w, h, 0, 0, 0, 0}
}

// ColorModel implements the ColorModel function in the Image interface.
func (t *Tile) ColorModel() color.Model {
	return t.TileImg.ColorModel()
}

// Bounds implements the Bounds function in the Image interface.
func (t *Tile) Bounds() Rectangle {
	return Rectangle{Point{-1e9, -1e9}, Point{1e9, 1e9}}
}

// At implements the At function in the Image interface.
func (t *Tile) At(x, y int) color.Color {
	if t.StagX != 0 {
		x += t.OffsX + t.StagX
		x %= t.Width
		if x < 0 {
			x = t.Width + x
		}
		y += t.OffsY
		y %= t.Height
		if y < 0 {
			y = t.Height + y
		}
		return t.TileImg.At(x, y)
	}
	x += t.OffsX
	x %= t.Width
	if x < 0 {
		x = t.Width + x
	}
	y += t.OffsY + t.StagY
	y %= t.Height
	if y < 0 {
		y = t.Height + y
	}
	return t.TileImg.At(x, y)
}
