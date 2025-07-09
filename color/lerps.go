package color

import "math"

// ColorRGBALerp calculates the color value at t [0,1] given a start and end color in RGB space.
func ColorRGBALerp(t float64, start, end Color) RGBA {
	rs, gs, bs, as := start.RGBA() // uint32 [0,0xffff]
	re, ge, be, ae := end.RGBA()
	omt := 1 - t
	rt := uint32(math.Floor(omt*float64(rs) + t*float64(re) + 0.5))
	gt := uint32(math.Floor(omt*float64(gs) + t*float64(ge) + 0.5))
	bt := uint32(math.Floor(omt*float64(bs) + t*float64(be) + 0.5))
	at := uint32(math.Floor(omt*float64(as) + t*float64(ae) + 0.5))
	rt >>= 8 // uint32 [0,0xff]
	gt >>= 8
	bt >>= 8
	at >>= 8
	return RGBA{uint8(rt), uint8(gt), uint8(bt), uint8(at)}
}

// ColorHSLLerp calculates the color value at t [0,1] given a start and end color in HSL space.
func ColorHSLLerp(t float64, start, end Color) HSL {
	cs, ce := NewHSL(start), NewHSL(end)
	ht := (1-t)*cs.H + t*ce.H // Will never cross 1:0
	st := (1-t)*cs.S + t*ce.S
	lt := (1-t)*cs.L + t*ce.L
	at := (1-t)*cs.A + t*ce.A
	return HSL{ht, st, lt, at}
}

// ColorHSLLerpS calculates the color value at t [0,1] given a start and end color in HSL space.
// Differs from ColorHSLLerp in that the shortest path for hue is taken.
func ColorHSLLerpS(t float64, start, end Color) HSL {
	cs, ce := NewHSL(start), NewHSL(end)
	hd := ce.H - cs.H
	ht := 0.0
	// Handle hue being circular
	if hd > 0.5 || hd < -0.5 {
		h := ce.H - 1
		ht = (1-t)*cs.H + t*h
		if ht < 0 {
			ht += 1
		}
	} else {
		ht = (1-t)*cs.H + t*ce.H
	}
	st := (1-t)*cs.S + t*ce.S
	lt := (1-t)*cs.L + t*ce.L
	at := (1-t)*cs.A + t*ce.A
	return HSL{ht, st, lt, at}
}
