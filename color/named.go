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
// and is licensed under an MIT license. Includes the CSS/HTML standard colors.

// There's no associated color model since this isn't a color space.

//go:embed colornames.bestof.csv
var b0 []byte

// ColorFile returns a []byte of the color csv file,
func ColorFile0() []byte {
	return b0
}

//go:embed colornames.css.csv
var b1 []byte

// ColorFile returns a []byte of the color csv file,
func ColorFile1() []byte {
	return b1
}

// NamedRGB contains the name of the color and its RGB color representation.
type NamedRGB struct {
	Name  string
	Color color.RGBA
}

// NamedRGBs is the slice of colors loaded from the color names file.
var BestNamedRGBs []*NamedRGB
var CSSNamedRGBs []*NamedRGB
var bestMap = make(map[string]*NamedRGB)
var cssMap = make(map[string]*NamedRGB)

func init() {
	reader := csv.NewReader(bytes.NewReader(ColorFile0()))
	BestNamedRGBs = parse(reader, bestMap)
	reader = csv.NewReader(bytes.NewReader(ColorFile1()))
	CSSNamedRGBs = parse(reader, cssMap)
}

func parse(reader *csv.Reader, nmap map[string]*NamedRGB) []*NamedRGB {
	if records, err := reader.ReadAll(); err == nil {
		n := len(records)
		res := make([]*NamedRGB, n-1)
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
			lcn := strings.ToLower(entry[0])
			nc := &NamedRGB{entry[0], col}
			nmap[lcn] = nc
			res[i-1] = nc
		}
		return res
	}
	return nil
}

// ByName returns the color given by the name. If there's no match, error will be set.
func ByName(name string) (*NamedRGB, error) {
	name = strings.ToLower(name)
	if col, prs := bestMap[name]; prs {
		return col, nil
	}

	return nil, fmt.Errorf("Color '%s' not found", name)
}

// ByCSSName returns the color given by the name. If there's no match, error will be set.
func ByCSSName(name string) (*NamedRGB, error) {
	name = strings.ToLower(name)
	if col, prs := cssMap[name]; prs {
		return col, nil
	}

	return nil, fmt.Errorf("Color '%s' not found", name)
}

// NamedRGBPalette performs a concrete to interface conversion
func NamedRGBPalette() []color.Color {
	nc := len(BestNamedRGBs)
	res := make([]color.Color, nc)
	for i, c := range BestNamedRGBs {
		res[i] = c.Color // Remove a level of indirection
	}
	return res
}

// RandomNamedRGB returns a random color from the list of named colors.
func RandomNamedRGB() *NamedRGB {
	return BestNamedRGBs[rand.Intn(len(BestNamedRGBs))]
}

// String returns a string represntation of NamedRGB.
func (nc *NamedRGB) String() string {
	return fmt.Sprintf("%s #%02x%02x%02xff", nc.Name, nc.Color.R, nc.Color.G, nc.Color.B)
}

// RGBA implements the color.Color interface.
func (nc *NamedRGB) RGBA() (uint32, uint32, uint32, uint32) {
	return nc.Color.RGBA()
}
