package color

import (
	"image/color"
	"math/rand"
)

// Predefined colors.
var (
	Black   = color.RGBA{0x00, 0x00, 0x00, 0xff}
	Red     = color.RGBA{0xff, 0x00, 0x00, 0xff}
	Green   = color.RGBA{0x00, 0xff, 0x00, 0xff}
	Blue    = color.RGBA{0x00, 0x00, 0xff, 0xff}
	Yellow  = color.RGBA{0xff, 0xff, 0x00, 0xff}
	Magenta = color.RGBA{0xff, 0x00, 0xff, 0xff}
	Cyan    = color.RGBA{0x00, 0xff, 0xff, 0xff}
	White   = color.RGBA{0xff, 0xff, 0xff, 0xff}

	DarkGray  = color.RGBA{0x63, 0x66, 0x6a, 0xff}
	MidGray   = color.RGBA{0x7f, 0x7f, 0x7f, 0xff}
	LightGray = color.RGBA{0xd9, 0xd9, 0xd6, 0xff}

	Brown  = color.RGBA{0xa4, 0x75, 0x51, 0xff}
	Orange = color.RGBA{0xff, 0xa5, 0x00, 0xff}
	Purple = color.RGBA{0x40, 0x00, 0x80, 0xff}

	GopherBlue  = color.RGBA{0x9d, 0xe8, 0xfd, 0xff}
	GopherBrown = color.RGBA{0xf3, 0xe2, 0xc9, 0xff}
	GopherGray  = color.RGBA{0xbd, 0xb9, 0xaf, 0xff}

	StandardPalette = []color.Color{
		Black,
		Red,
		Orange,
		Brown,
		Yellow,
		Green,
		Blue,
		Cyan,
		Purple,
		Magenta,
		White,
		MidGray,
	}
)

// Random returns a randomized color in R, G and B. Alpha is set to 0xff.
func Random() color.RGBA {
	colval := rand.Uint32()
	b := uint8(colval & 0xff)
	colval >>= 8
	g := uint8(colval & 0xff)
	colval >>= 8
	r := uint8(colval & 0xff)
	return color.RGBA{r, g, b, 0xff}
}

// RandomFromPalette selects a color from the supplied palette.
func RandomFromPalette(palette []color.Color) color.Color {
	nc := len(palette)
	return palette[rand.Intn(nc)]
}
