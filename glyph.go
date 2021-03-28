package graphics2d

import (
	"fmt"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

/* Boilerplate to read a font collection in:
//Read the data
fontData, err := ioutil.ReadFile(fontPath)
if err != nil {
	log.Fatalf("Error reading font\n")
}

// Load with sfnt package
fonts, err := sfnt.ParseCollection(fontData)
fmt.Printf("Found %d fonts\n", fonts.NumFonts())
var b sfnt.Buffer
for i:=0; i<fonts.NumFonts(); i++ {
	f, _ := fonts.Font(i)
	name, _ := f.Name(&b, sfnt.NameIDFull)
	fmt.Printf("Font %d: %s\n", i, name)
}
f, _ := fonts.Font(0)
*/

// Points are in font units
func GlyphToShape(font *sfnt.Font, r rune) (*Shape, error) {
	shape := &Shape{}

	var buffer sfnt.Buffer
	x, err := font.GlyphIndex(&buffer, r)
	if err != nil {
		return shape, err
	}
	if x == 0 {
		return shape, fmt.Errorf("rune %c not found in font", r)
	}
	// LoadGlyph(b *Buffer, x GlyphIndex, ppem fixed.Int26_6, opts *LoadGlyphOptions) (Segments, error)
	segments, err := font.LoadGlyph(&buffer, x, fixed.Int26_6(font.UnitsPerEm()), nil)
	if err != nil {
		return shape, err
	}

	var cp *Path
	for _, seg := range segments {
		// The divisions by 64 below are because the seg.Args values have type
		// fixed.Int26_6, a 26.6 fixed point number, and 1<<6 == 64.
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			if cp != nil {
				cp.Close()
				shape.AddPath(cp)
			}
			cp = NewPath([]float64{float64(seg.Args[0].X) / 64, float64(seg.Args[0].Y) / 64})
		case sfnt.SegmentOpLineTo:
			cp.AddStep([][]float64{{float64(seg.Args[0].X) / 64, float64(seg.Args[0].Y) / 64}})
		case sfnt.SegmentOpQuadTo:
			cp.AddStep([][]float64{{float64(seg.Args[0].X) / 64, float64(seg.Args[0].Y) / 64},
				{float64(seg.Args[1].X) / 64, float64(seg.Args[1].Y) / 64}})
		case sfnt.SegmentOpCubeTo:
			cp.AddStep([][]float64{{float64(seg.Args[0].X) / 64, float64(seg.Args[0].Y) / 64},
				{float64(seg.Args[1].X) / 64, float64(seg.Args[1].Y) / 64},
				{float64(seg.Args[2].X) / 64, float64(seg.Args[2].Y) / 64}})
		}
	}
	cp.Close()
	shape.AddPath(cp)
	return shape, nil
}
