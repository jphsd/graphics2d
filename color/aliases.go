package color

import "image/color"

type Color = color.Color

type Alpha = color.Alpha
type Alpha16 = color.Alpha16
type Gray = color.Gray
type Gray16 = color.Gray16
type RGBA = color.RGBA
type NRGBA = color.NRGBA
type RGBA64 = color.RGBA64
type NRGBA64 = color.NRGBA64

type Model = color.Model

var (
	GrayModel    Model = color.GrayModel
	Gray16Model  Model = color.Gray16Model
	RGBAModel    Model = color.RGBAModel
	NRGBAModel   Model = color.NRGBAModel
	RGBA64Model  Model = color.RGBA64Model
	NRGBA64Model Model = color.NRGBA64Model

	ModelFunc func(func(Color) Color) Model = color.ModelFunc
)
