package image

import (
	"image"
	"image/color"
	"image/draw"
)

// Tile is an infinite image covered with a tile.
type Tile struct {
	Tile   *image.RGBA
	Width  int
	Height int
	OffsX  int
	OffsY  int
}

// NewTile creates a new image with the supplied image tile.
func NewTile(img image.Image) *Tile {
	rect := img.Bounds()
	w, h := rect.Dx(), rect.Dy()
	tile := image.NewRGBA(image.Rectangle{image.Point{}, rect.Size()})
	draw.Draw(tile, tile.Bounds(), img, rect.Min, draw.Src)
	return &Tile{tile, w, h, 0, 0}
}

// ColorModel implements the ColorModel function in the Image interface.
func (t *Tile) ColorModel() color.Model {
	return t.Tile.ColorModel()
}

// Bounds implements the Bounds function in the Image interface.
func (t *Tile) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{-1e9, -1e9}, image.Point{1e9, 1e9}}
}

// At implements the At function in the Image interface.
func (t *Tile) At(x, y int) color.Color {
	x += t.OffsX
	y += t.OffsY
	x %= t.Width
	if x < 0 {
		x = t.Width - x
	}
	y %= t.Height
	if y < 0 {
		y = t.Height - y
	}
	return t.Tile.At(x, y)
}
