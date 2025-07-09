package color

import "image/color"

type Color = color.Color

type Gray = color.Gray
type Gray16 = color.Gray16
type RGBA = color.RGBA
type NRGBA = color.NRGBA
type RGBA64 = color.RGBA64
type NRGBA64 = color.NRGBA64
type Model = color.Model

var ModelFunc func(func(Color) Color) Model = color.ModelFunc
