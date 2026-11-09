package graphics2d

import (
	"image/draw"
)

type Image draw.Image
type Op draw.Op

const (
	Over = Op(draw.Over)
	Src  = Op(draw.Src)
)
