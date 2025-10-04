package color

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	// force init for go:embed
	_ "embed"
)

// colornames.bestof.csv comes from the color-names project (https://github.com/meodai/color-names)
// and is licensed under an MIT license. Includes the CSS/HTML standard colors.

// There's no associated color model since this isn't a color space.

//go:embed colornames.bestof.csv
var b0 []byte

// ColorFile returns a []byte of the colornames.bestof csv file,
func ColorFile0() []byte {
	return b0
}

//go:embed colornames.css.csv
var b1 []byte

// ColorFile returns a []byte of the colornames.css csv file,
func ColorFile1() []byte {
	return b1
}

// NamedRGB contains the name of the color and its RGB color representation.
type NamedRGB struct {
	Name  string
	Color RGBA
}

var (
	// NamedRGBs is the slice of colors loaded from the best color names file.
	BestNamedRGBs []*NamedRGB

	// CSSNamedRGBs is the slice of colors loaded from the W3 CSS color names file.
	CSSNamedRGBs []*NamedRGB

	bestMap = make(map[string]*NamedRGB)
	cssMap  = make(map[string]*NamedRGB)
)

func init() {
	BestNamedRGBs = NewNamedRGBSlice(ColorFile0(), bestMap)
	CSSNamedRGBs = NewNamedRGBSlice(ColorFile1(), cssMap)
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
			col := RGBA{r, g, b, 0xff}
			lcn := strings.ToLower(entry[0])
			nc := &NamedRGB{entry[0], col}
			nmap[lcn] = nc
			res[i-1] = nc
		}
		return res
	}
	return nil
}

// NewNamedRGBSlice parses the data provided to it into a slice of named colors.
// It uses csv.Reader to perform the parsing and expects the name in first entry and the color
// in the second as #RRGGBB hex.
// The first line of data is assumed to be a header and is skipped.
func NewNamedRGBSlice(data []byte, names map[string]*NamedRGB) []*NamedRGB {
	reader := csv.NewReader(bytes.NewReader(data))
	return parse(reader, names)
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

// String returns a string represntation of NamedRGB.
func (nc *NamedRGB) String() string {
	return fmt.Sprintf("%s #%02x%02x%02xff", nc.Name, nc.Color.R, nc.Color.G, nc.Color.B)
}

// RGBA implements the color.Color interface.
func (nc *NamedRGB) RGBA() (uint32, uint32, uint32, uint32) {
	return nc.Color.RGBA()
}
