package image

import (
	"fmt"
	g2dcol "github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/util"
	"image/color"
	"math"
)

// CreateLutFromValues maps a series of values in [0,1] to [0,255]
// Note - no checking on values range
func CreateLutFromValues(values []float64) []uint8 {
	n := len(values)
	res := make([]uint8, n)
	for i := 0; i < n; i++ {
		res[i] = uint8(values[i] * 255)
	}
	return res
}

// PremulLut creates a slice of luts from lut for every possible value of alpha
func PremulLut(lut []uint8) [][]uint8 {
	n := len(lut)
	res := make([][]uint8, n)
	for i := 0; i < n; i++ {
		nres := make([]uint8, 256)
		v := uint32(lut[i])
		v |= v << 8
		for a := 0; a < 256; a++ {
			nv := uint32(a) * v
			nv /= 0xff
			nres[a] = uint8(nv >> 8)
		}
		res[i] = nres
	}
	return res
}

// NLExpansionLut generates a lut [start,end), normalized and mapped through f.
// For example NLExpansionLut(256, 0, 256, &NLSin{}) will generate a Sin ramp.
func NLExpansionLut(n, start, end int, f util.NonLinear) []uint8 {
	if start > end {
		start, end = end, start
	}
	if start < 0 || start > n-1 || end < 1 || end > n {
		panic(fmt.Errorf("start or end not in range"))
	}
	res := make([]uint8, n)

	var i int
	for ; i < start; i++ {
		res[i] = 0
	}
	delta := end - start
	if delta > 0 {
		div := 1 / float64(delta)
		for ; i < end; i++ {
			v := float64(i-start) * div
			v = f.Transform(v)
			res[i] = uint8(v * 0xff)
		}
	}
	for ; i < n; i++ {
		res[i] = 0xff
	}

	return res
}

// NLColorLut generates a color lut using the tvals ((0,1) ascending, strict) and colors given in either HSL or RGB space.
func NLColorLut(n int, f util.NonLinear, start, end color.Color, hsl bool, tvals []float64, colors []color.Color) []color.Color {
	if n < 3 {
		return []color.Color{start, end}
	}

	nv := len(tvals)
	if len(colors) < nv {
		nv = len(colors)
	}
	dt := 1 / float64(n-1)

	// Initial is just an NLerp between start and end
	res := make([]color.Color, n)
	t, ts, te := 0.0, 0.0, 1.0
	fs, fe, fd := ts, te, te-ts
	ci, cs, ce := 0, start, end
	if ci < nv {
		te = tvals[ci]
		ce = colors[ci]
		fe = f.Transform(te)
		fd = fe - fs
		ci++
	}

	for i := 0; i < n; i++ {
		ft := f.Transform(t)
		ftp := (ft - fs) / fd
		if hsl {
			res[i] = ColorHSLLerp(ftp, cs, ce)
		} else {
			res[i] = ColorRGBALerp(ftp, cs, ce)
		}
		t += dt
		for t > te {
			ts, fs, cs = te, fe, ce
			if ci < nv {
				te = tvals[ci]
				ce = colors[ci]
				fe = f.Transform(te)
				fd = fe - fs
				ci++
			} else {
				te = 1
				ce = end
				fe = 1
				fd = fe - fs
			}
		}
	}
	return res
}

// ColorRGBALerp calulates the color value at t [0,1] given a start and end color in RGB space.
func ColorRGBALerp(t float64, start, end color.Color) color.Color {
	rs, gs, bs, as := start.RGBA() // uint32 [0,0xffff]
	re, ge, be, ae := end.RGBA()
	rt := uint32(math.Floor((1-t)*float64(rs) + t*float64(re) + 0.5))
	gt := uint32(math.Floor((1-t)*float64(gs) + t*float64(ge) + 0.5))
	bt := uint32(math.Floor((1-t)*float64(bs) + t*float64(be) + 0.5))
	at := uint32(math.Floor((1-t)*float64(as) + t*float64(ae) + 0.5))
	rt >>= 8 // uint32 [0,0xff]
	gt >>= 8
	bt >>= 8
	at >>= 8
	return &color.RGBA{uint8(rt), uint8(gt), uint8(bt), uint8(at)}
}

// ColorHSLLerp calulates the color value at t [0,1] given a start and end color in HSL space.
func ColorHSLLerp(t float64, start, end color.Color) color.Color {
	cs, ce := g2dcol.NewHSL(start), g2dcol.NewHSL(end)
	ht := (1-t)*cs.H + t*ce.H
	st := (1-t)*cs.S + t*ce.S
	lt := (1-t)*cs.L + t*ce.L
	at := (1-t)*cs.A + t*ce.A
	return &g2dcol.HSL{ht, st, lt, at}
}
