package image

import (
	"image"
	"image/color"
)

type Image = image.Image
type Alpha = image.Alpha
type Alpha16 = image.Alpha16
type RGBA = image.RGBA
type RGBA64 = image.RGBA64
type NRGBA = image.NRGBA
type NRGBA64 = image.NRGBA64
type Gray = image.Gray
type Gray16 = image.Gray16
type Uniform = image.Uniform
type Point = image.Point
type Rectangle = image.Rectangle

var Rect func(int, int, int, int) Rectangle = image.Rect
var NewUniform func(color.Color) *Uniform = image.NewUniform
var Opaque = image.Opaque
var Transparent = image.Transparent
