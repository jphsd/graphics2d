package color

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	// force init for go:embed
	_ "embed"
)

// colornames.bestof.csv comes from the color-names project (https://github.com/meodai/color-names)
// and is licensed under an MIT license.

//go:embed colornames.bestof.csv
var b []byte

// ColorFile retruns a []byte of the color csv file,
func ColorFile() []byte {
	return b
}

// NamedColor contains the name of the color and its color representation.
type NamedColor struct {
	Name  string
	Color color.Color
}

// NamedColors is the slice of colors loaded from the color names file.
var NamedColors []NamedColor
var namedColors = make(map[string]color.Color)

func init() {
	reader := csv.NewReader(bytes.NewReader(ColorFile()))
	if records, err := reader.ReadAll(); err == nil {
		n := len(records)
		NamedColors = make([]NamedColor, n-1)
		// Skip header
		for i := 1; i < n; i++ {
			entry := records[i]
			var colval uint32
			_, err = fmt.Sscanf(entry[1], "#%x", &colval)
			if err != nil {
				fmt.Printf("%s\n", err)
				break
			}
			b := uint8(colval & 0xff)
			colval >>= 8
			g := uint8(colval & 0xff)
			colval >>= 8
			r := uint8(colval & 0xff)
			col := color.RGBA{r, g, b, 0xff}
			namedColors[strings.ToLower(entry[0])] = col
			NamedColors[i-1] = NamedColor{entry[0], col}
		}
	}
}

// ByName returns the color given by the name. If there's no match, error will be set.
func ByName(name string) (color.Color, error) {
	name = strings.ToLower(name)
	if col, prs := namedColors[name]; prs {
		return col, nil
	}

	return nil, fmt.Errorf("Color '%s' not found", name)
}

// RandomNamed returns a random color from the list of named colors.
func RandomNamed() NamedColor {
	return NamedColors[rand.Intn(len(NamedColors))]
}

// Predefined colors.
var (
	Black   = color.RGBA{0x0,  0x0,  0x0,  0xff}
	White   = color.RGBA{0xff, 0xff, 0xff, 0xff}
	Red     = color.RGBA{0xff, 0x0,  0x0,  0xff}
	Green   = color.RGBA{0x0,  0xff, 0x0,  0xff}
	Blue    = color.RGBA{0x0,  0x0,  0xff, 0xff}
	Yellow  = color.RGBA{0xff, 0xff, 0x0,  0xff}
	Magenta = color.RGBA{0xff, 0x0,  0xff, 0xff}
	Cyan    = color.RGBA{0x0,  0xff, 0xff, 0xff}
	Orange  = color.RGBA{0xff, 0xa5, 0x0,  0xff}
)

// Random returns a randomized color in R, G and B. Alpha is set to 0xff.
func Random() color.Color {
	colval := rand.Uint32()
	b := uint8(colval & 0xff)
	colval >>= 8
	g := uint8(colval & 0xff)
	colval >>= 8
	r := uint8(colval & 0xff)
	return color.RGBA{r, g, b, 0xff}
}
