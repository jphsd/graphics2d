package surface

import (
	"image"
	"image/color"
)

// NormalMap provides an At function to determine the unit normal at a location.
type NormalMap interface {
	At(x, y int) []float64
}

// DefaultNM is a map where all normals are {0, 0, 1}.
type DefaultNM struct{}

// At implements the At function of the NormalMap interface.
func (n *DefaultNM) At(x, y int) []float64 {
	return []float64{0, 0, 1}
}

// ImageNM describes a normal map derived from an image. The width and height of the map are the same
// as the image and a flag indicates if the map should be tiled for values outside of that range.
type ImageNM struct {
	Width, Height int
	Tiled         bool
	Normals       [][][]float64
}

// NewImageNM creates a normal map from an image. The input mage is converted to Grey16 and then
// the finite difference method is used to construct the normals. The x and y components of the
// normal are scaled by scale prior to unitization.
func NewImageNM(img image.Image, scale float64, inv, tiled bool) *ImageNM {
	rect := img.Bounds()
	w, h := rect.Dx(), rect.Dy()
	g16m := color.Gray16Model
	ox, oy := rect.Min.X, rect.Min.Y
	values := make([][]int, h)
	normals := make([][][]float64, h)
	for r := 0; r < h; r++ {
		values[r] = make([]int, w)
		normals[r] = make([][]float64, w)
		for c := 0; c < w; c++ {
			col := img.At(c+ox, r+oy)
			col = g16m.Convert(col)
			gcol, _ := col.(color.Gray16)
			if inv {
				gcol.Y = 0xffff - gcol.Y
			}
			values[r][c] = int(gcol.Y)
		}
	}

	// values contains the gray scale pixel values as int
	denom := float64(0xffff) / scale
	for y := 0; y < h-1; y++ {
		cur := values[y][0]
		for x := 0; x < w-1; x++ {
			nx := values[y][x+1]
			ny := values[y+1][x]
			dx, dy := cur-nx, cur-ny
			fx, fy := float64(dx)/denom, float64(dy)/denom
			normals[y][x] = norm([]float64{fx, fy, 1}) // unit normal
			cur = nx
		}
	}
	// Handle last column - repeat dx
	for y := 0; y < h-1; y++ {
		x := w - 1
		cur := values[y][x]
		nx := values[y][x-1]
		ny := values[y+1][x]
		dx, dy := nx-cur, cur-ny
		fx, fy := float64(dx)/denom, float64(dy)/denom
		normals[y][x] = norm([]float64{fx, fy, 1}) // unit normal
	}
	// Handle last row - repeat dy
	for x := 0; x < w-1; x++ {
		y := h - 1
		cur := values[y][x]
		nx := values[y][x+1]
		ny := values[y-1][x]
		dx, dy := cur-nx, ny-cur
		fx, fy := float64(dx)/denom, float64(dy)/denom
		normals[y][x] = norm([]float64{fx, fy, 1}) // unit normal
	}
	// Handle last point
	cur := values[h-1][w-1]
	dx := values[h-1][w-2] - cur
	dy := values[h-2][w-1] - cur
	fx, fy := float64(dx)/denom, float64(dy)/denom
	normals[h-1][w-1] = norm([]float64{fx, fy, 1}) // unit normal

	return &ImageNM{w, h, tiled, normals}
}

// At implements the At function of the NormalMap interface. If the map isn't tiled, values out of
// range return {0, 0, 1}.
func (n *ImageNM) At(x, y int) []float64 {
	if n.Tiled {
		x %= n.Width
		if x < 0 {
			x = n.Width - x
		}
		y %= n.Height
		if y < 0 {
			y = n.Height - y
		}
	} else if x < 0 || x >= n.Width || y < 0 || y >= n.Height {
		return []float64{0, 0, 1}
	}
	return n.Normals[y][x]
}

// NormalsToImage converts unit normals to an RGBA image with x in R, y in G and z in B. Note, the range
// [-1,1] is mapped to [0,1] for use by FRGB so {-1, 0, 1} => {0x0, 0x7f, 0xff}.
func NormalsToImage(normals [][][]float64) *image.RGBA {
	h, w := len(normals), len(normals[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			norm := normals[y][x]
			img.Set(x, y, &FRGB{(norm[0] + 1) / 2, (norm[1] + 1) / 2, (norm[2] + 1) / 2})
		}
	}
	return img
}
