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

var (
	fontCache = make(map[*sfnt.Font]map[rune]*Shape)
	giCache   = make(map[*sfnt.Font]map[rune]sfnt.GlyphIndex)
)

// GlyphToShape returns a shape containing the paths for rune r as found in the font.
// The path is in font units.
// Use font.UnitsPerEm() to calculate scale factors.
func GlyphToShape(font *sfnt.Font, r rune) (*Shape, error) {
	rcache, ok := fontCache[font]
	if !ok {
		rcache = make(map[rune]*Shape)
		fontCache[font] = rcache
		giCache[font] = make(map[rune]sfnt.GlyphIndex)
	}
	shape, ok := rcache[r]
	if !ok {
		var buffer sfnt.Buffer
		x, err := font.GlyphIndex(&buffer, r)
		if err != nil {
			return nil, err
		}
		giCache[font][r] = x
		// x == 0 means use the unfound glyph
		shape, err = GlyphIndexToShape(font, x)
		if err != nil {
			return nil, err
		}
		rcache[r] = shape
	}
	return shape, nil
}

// GlyphIndexToShape returns a shape containing the paths for glyph index x as found in the font. The path is in
// font units. Use font.UnitsPerEm() to calculate scale factors.
func GlyphIndexToShape(font *sfnt.Font, x sfnt.GlyphIndex) (*Shape, error) {
	var buffer sfnt.Buffer
	segments, err := font.LoadGlyph(&buffer, x, fixed.I(int(font.UnitsPerEm())), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading glyph for index %d, (%s)", x, err.Error())
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

// StringToShape returns the string rendered as both a single shape,
// and as individual shapes, correctly offset in font units.
// Glyphs with no paths are not returned (e.g. space etc.).
func StringToShape(tfont *sfnt.Font, str string) (*Shape, []*Shape, error) {
	r2gi, ok := giCache[tfont]
	if !ok {
		r2gi = make(map[rune]sfnt.GlyphIndex)
		giCache[tfont] = r2gi
	}
	r2s, ok := fontCache[tfont]
	if !ok {
		r2s = make(map[rune]*Shape)
		fontCache[tfont] = r2s
	}
	r2adv := make(map[rune]float64)
	upem := fixed.I(int(tfont.UnitsPerEm()))

	var buffer sfnt.Buffer
	x := 0.0

	shape := &Shape{}
	shapes := []*Shape{}
	pgi := sfnt.GlyphIndex(0xffff)
	var pr rune
	for i, r := range str {
		s, ok := r2s[r]
		if !ok {
			// Find the glyph index
			gi, err := tfont.GlyphIndex(&buffer, r)
			if err != nil {
				return nil, nil, fmt.Errorf("gi error at rune %d (%s)", i, err.Error())
			}
			r2gi[r] = gi
			// Create a shape for it
			s, err = GlyphIndexToShape(tfont, gi)
			if err != nil {
				return nil, nil, err
			}
			r2s[r] = s
			// Lookup its advance and convert it to float64
			adv, err := tfont.GlyphAdvance(&buffer, gi, upem, font.HintingNone)
			if err != nil {
				return nil, nil, fmt.Errorf("error finding advance for index %d(%c), (%s)", gi, r, err.Error())
			}
			r2adv[r] = I266ToF64(adv)
		}
		gi := r2gi[r]
		if pgi != 0xffff {
			// Apply any kerning
			kern, err := tfont.Kern(&buffer, pgi, gi, upem, font.HintingNone)
			k := 0.0
			if err != nil && err != sfnt.ErrNotFound {
				return nil, nil, fmt.Errorf("error finding kerning for %d(%c) and %d(%c), (%s)", gi, r, pgi, pr, err.Error())
			} else {
				k = I266ToF64(kern)
			}
			x += k
		}
		pgi = gi
		pr = r
		if len(s.Paths()) > 0 {
			// Add to result shape
			xfm := Translate(x, 0)
			s = s.Transform(xfm)
			shapes = append(shapes, s)
			shape.AddShapes(s)
		}
		x += r2adv[r]
	}
	return shape, shapes, nil
}

// I266ToF64 converts a fixed.Int26_6 to float64
func I266ToF64(fi fixed.Int26_6) float64 {
	return float64(fi>>6) + 0.015625*float64(fi&0x3f)
}
