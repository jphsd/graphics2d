package graphics2d

import (
	"fmt"
	"golang.org/x/image/font"
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
f := fonts.Font(0)
*/

// GlyphToShape returns a shape containing the paths for rune r as found in the font. The path is in
// font units. Use font.UnitsPerEm() to calculate scale factors.
func GlyphToShape(font *sfnt.Font, r rune) (*Shape, error) {
	var buffer sfnt.Buffer
	x, err := font.GlyphIndex(&buffer, r)
	if err != nil {
		return nil, err
	}
	// x == 0 means use the unfound glyph
	return GlyphIndexToShape(font, x)
}

// GlyphIndexToShape returns a shape containing the paths for glyph index x as found in the font. The path is in
// font units. Use font.UnitsPerEm() to calculate scale factors.
func GlyphIndexToShape(font *sfnt.Font, x sfnt.GlyphIndex) (*Shape, error) {
	var buffer sfnt.Buffer
	segments, err := font.LoadGlyph(&buffer, x, fixed.I(int(font.UnitsPerEm())), nil)
	if err != nil {
		return nil, err
	}

	var cp *Path
	shape := &Shape{}
	for _, seg := range segments {
		// The divisions by 64 below are because the seg.Args values have type
		// fixed.Int26_6, a 26.6 fixed point number, and 1<<6 == 64.
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			if cp != nil {
				cp.Close()
				shape.AddPaths(cp)
			}
			cp = NewPath([]float64{I266ToF64(seg.Args[0].X), I266ToF64(seg.Args[0].Y)})
		case sfnt.SegmentOpLineTo:
			cp.AddStep([]float64{I266ToF64(seg.Args[0].X), I266ToF64(seg.Args[0].Y)})
		case sfnt.SegmentOpQuadTo:
			cp.AddStep([]float64{I266ToF64(seg.Args[0].X), I266ToF64(seg.Args[0].Y)},
				[]float64{I266ToF64(seg.Args[1].X), I266ToF64(seg.Args[1].Y)})
		case sfnt.SegmentOpCubeTo:
			cp.AddStep([]float64{I266ToF64(seg.Args[0].X), I266ToF64(seg.Args[0].Y)},
				[]float64{I266ToF64(seg.Args[1].X), I266ToF64(seg.Args[1].Y)},
				[]float64{I266ToF64(seg.Args[2].X), I266ToF64(seg.Args[2].Y)})
		}
	}
	if cp != nil {
		cp.Close()
		shape.AddPaths(cp)
	}
	return shape, nil
}

// StringToShape returns the string rendered as both a single shape, and as individual shapes, correctly
// offset. Glyphs with no paths are not returned (e.g. space etc.).
func StringToShape(tfont *sfnt.Font, str string) (*Shape, []*Shape, error) {
	r2gi := make(map[rune]sfnt.GlyphIndex)
	gi2s := make(map[sfnt.GlyphIndex]*Shape)
	gi2adv := make(map[sfnt.GlyphIndex]float64)
	upem := fixed.I(int(tfont.UnitsPerEm()))

	var buffer sfnt.Buffer
	var err error
	x := 0.0

	shape := &Shape{}
	shapes := []*Shape{}
	for i, r := range str {
		gi, ok := r2gi[r]
		var s *Shape
		if !ok {
			// Find the glyph index
			gi, err = tfont.GlyphIndex(&buffer, r)
			if err != nil {
				return nil, nil, fmt.Errorf("error at rune %d (%s)", i, err.Error())
			}
			r2gi[r] = gi
			// Create a shape for it
			s, err = GlyphIndexToShape(tfont, gi)
			if err != nil {
				return nil, nil, err
			}
			gi2s[gi] = s
			// Lookup its advance and convert it to float64
			adv, err := tfont.GlyphAdvance(&buffer, gi, upem, font.HintingNone)
			if err != nil {
				return nil, nil, err
			}
			gi2adv[gi] = I266ToF64(adv)
		} else {
			s = gi2s[gi]
		}
		if len(s.Paths()) > 0 {
			// Add to result shape
			xfm := Translate(x, 0)
			s = s.Transform(xfm)
			shapes = append(shapes, s)
			shape.AddShapes(s)
		}
		x += gi2adv[gi]
	}
	return shape, shapes, nil
}

// I266ToF64 converts a fixed.Int26_6 to float64
func I266ToF64(fi fixed.Int26_6) float64 {
	return float64(fi>>6) + 0.015625*float64(fi&0x3f)
}
