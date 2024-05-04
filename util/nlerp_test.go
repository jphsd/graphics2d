package util

import "testing"

type testNLRow struct {
	nl   NonLinear
	name string
}

var tableNL []testNLRow

func init() {
	tableNL = []testNLRow{
		{&NLLinear{}, "Linear"},
		{&NLSquare{}, "Square"},
		{&NLCube{}, "Cube"},
		{&NLCircle1{}, "Circle1"},
		{&NLCircle2{}, "Circle2"},
		{NewNLLame(2, 2), "Lame 1"},
		{NewNLLame(0.5, 0.5), "Lame 2"},
		{&NLSin{}, "Sin"},
		{&NLSin1{}, "Sin1"},
		{&NLSin2{}, "Sin2"},
		{&NLCatenary{}, "Catenary"},
		{NewNLExponential(1), "Exponential 1"},
		{NewNLExponential(10), "Exponential 2"},
		{NewNLExponential(100), "Exponential 3"},
		{NewNLLogarithmic(1), "Logarithmic 1"},
		{NewNLLogarithmic(10), "Logarithmic 2"},
		{NewNLLogarithmic(100), "Logarithmic 3"},
		{NewNLGauss(1), "Gauss 1"},
		{NewNLGauss(3), "Gauss 2"},
		{NewNLGauss(6), "Gauss 3"},
		{NewNLLogistic(1, 0.5), "Logistic 1"},
		{NewNLLogistic(12, 0.5), "Logistic 2"},
		{NewNLLogistic(60, 0.5), "Logistic 3"},
		{NewNLLogistic(1, 0.2), "Logistic 4"},
		{NewNLLogistic(12, 0.2), "Logistic 5"},
		{NewNLLogistic(38, 0.2), "Logistic 6"},
		{NewNLLogistic(1, 0.8), "Logistic 7"},
		{NewNLLogistic(12, 0.8), "Logistic 8"},
		{NewNLLogistic(100, 0.8), "Logistic 9"},
		// {&NLP3{}, "P3"},
		// {&NLP5{}, "P5"},
		{&NLCompound{[]NonLinear{&NLCube{}, &NLSin{}}}, "Compound"},
		{&NLOmt{&NLCube{}}, "OneMinusT Cube"},
		{&NLOmt{&NLSin{}}, "OneMinusT Sin"},
		{&NLOmt{&NLCircle2{}}, "OneMinusT Circle2"},
		{NewNLRand(0.1, 0.01, true), "Rand Sharp"},
		{NewNLRand(0.1, 0.01, false), "Rand Flat"},
	}
}

const (
	// level of accuracy on bsInv()
	epsilon = 0.0005
)

func Equalf64(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < epsilon
}

// Big numbers so within 1 is okay
func Equalf32(a, b float32) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 1
}

func TestNLFuncs(t *testing.T) {
	for _, row := range tableNL {
		vp := -epsilon
		ttp := -epsilon
		for t1 := 0.; t1 <= 1; t1 += 1 / 100. {
			v := row.nl.Transform(t1)
			// Range check
			if v < -epsilon || v > 1 {
				t.Errorf("%s t %f -> v %f out of range [0,1]", row.name, t1, v)
			}
			// Monotonicity check
			if v < vp {
				t.Errorf("%s vp %f < v %f", row.name, vp, v)
			}
			vp = v

			tt := row.nl.InvTransform(v)
			// Monotonicity check
			if tt < ttp {
				t.Errorf("%s ttp %f < tt %f", row.name, ttp, tt)
			}
			ttp = tt
			// Range check
			if tt < -epsilon || tt > 1 {
				t.Errorf("%s v %f -> t %f out of range [0,1]", row.name, v, tt)
			}

			// InvTransfrom(Transform(t)) == t check
			if !Equalf64(t1, tt) {
				t.Errorf("%s %f -> %f -> %f", row.name, t1, v, tt)
			}
		}
	}
}

func TestNLerp(t *testing.T) {
	start, end := float64(10), float64(90)
	for _, row := range tableNL {
		vp := start
		var t1 float64
		for t1 = 0; t1 <= 1; t1 += 1 / 100. {
			v := NLerp(t1, start, end, row.nl)
			// Range check
			if v < start || v > end {
				t.Errorf("%s t %f -> v %f out of range [%f,%f]", row.name, t1, v, start, end)
			}
			// Monotonicity check
			if v < vp {
				t.Errorf("%s vp %f < v %f", row.name, vp, v)
			}
			vp = v
		}
		// t out of range [0, 1]
		t1 = -1
		v := NLerp(t1, start, end, row.nl)
		if v != start {
			t.Errorf("%s t %f -> v %f expected %f", row.name, t1, v, start)
		}
		t1 = 2
		v = NLerp(t1, start, end, row.nl)
		if v != end {
			t.Errorf("%s t %f -> v %f expected %f", row.name, t1, v, end)
		}
	}
}

func TestInvNLerp(t *testing.T) {
	start, end := float64(10), float64(90)
	for _, row := range tableNL {
		vp := -epsilon
		var t1 float64
		for t1 = 10; t1 < end; t1 += 7 {
			v := InvNLerp(t1, start, end, row.nl)
			// Range check
			if v < -epsilon || v > 1 {
				t.Errorf("%s v %f -> t %f out of range [0,1]", row.name, t1, v)
			}
			// Monotonicity check
			if v < vp {
				t.Errorf("%s vp %f < t %f", row.name, vp, v)
			}
			vp = v
		}
		// t out of range [0, 1]
		t1 = 0
		v := InvNLerp(t1, start, end, row.nl)
		if v != 0 {
			t.Errorf("%s v %f -> t %f expected %f", row.name, t1, v, start)
		}
		t1 = 100
		v = InvNLerp(t1, start, end, row.nl)
		if v != 1 {
			t.Errorf("%s v %f -> t %f expected %f", row.name, t1, v, end)
		}
	}
}

func TestRemapNL(t *testing.T) {
	istart, iend := float64(10), float64(90)
	ostart, oend := float64(-10), float64(-90)
	iv := (istart + iend) / 2
	ov := (ostart + oend) / 2
	for _, row := range tableNL {
		v := RemapNL(iv, istart, iend, ostart, oend, row.nl, row.nl)
		if !Equalf64(v, ov) {
			t.Errorf("%s %f -> %f expected %f", row.name, iv, v, ov)
		}
	}
}

func TestNLerp32(t *testing.T) {
	start, end := float32(10), float32(90)
	for _, row := range tableNL {
		vp := start
		var t1 float32
		for t1 = 0; t1 <= 1; t1 += 1 / 100. {
			v := NLerp32(t1, start, end, row.nl)
			// Range check
			if v < start || v > end {
				t.Errorf("%s t %f -> v %f out of range [%f,%f]", row.name, t1, v, start, end)
			}
			// Monotonicity check
			if v < vp {
				t.Errorf("%s vp %f < v %f", row.name, vp, v)
			}
			vp = v
		}
		// t out of range [0, 1]
		t1 = -1
		v := NLerp32(t1, start, end, row.nl)
		if v != start {
			t.Errorf("%s t %f -> v %f expected %f", row.name, t1, v, start)
		}
		t1 = 2
		v = NLerp32(t1, start, end, row.nl)
		if v != end {
			t.Errorf("%s t %f -> v %f expected %f", row.name, t1, v, end)
		}
	}
}

func TestInvNLerp32(t *testing.T) {
	start, end := float32(10), float32(90)
	for _, row := range tableNL {
		vp := float32(-epsilon)
		var t1 float32
		for t1 = 10; t1 < end; t1 += 7 {
			v := InvNLerp32(t1, start, end, row.nl)
			// Range check
			if v < -epsilon || v > 1 {
				t.Errorf("%s v %f -> t %f out of range [0,1]", row.name, t1, v)
			}
			// Monotonicity check
			if v < vp {
				t.Errorf("%s vp %f < t %f", row.name, vp, v)
			}
			vp = v
		}
		// t out of range [0, 1]
		t1 = 0
		v := InvNLerp32(t1, start, end, row.nl)
		if v != 0 {
			t.Errorf("%s v %f -> t %f expected %f", row.name, t1, v, start)
		}
		t1 = 100
		v = InvNLerp32(t1, start, end, row.nl)
		if v != 1 {
			t.Errorf("%s v %f -> t %f expected %f", row.name, t1, v, end)
		}
	}
}

func TestRemapNL32(t *testing.T) {
	istart, iend := float32(10), float32(90)
	ostart, oend := float32(-10), float32(-90)
	iv := (istart + iend) / 2
	ov := (ostart + oend) / 2
	for _, row := range tableNL {
		v := RemapNL32(iv, istart, iend, ostart, oend, row.nl, row.nl)
		if !Equalf32(v, ov) {
			t.Errorf("%s %f -> %f expected %f", row.name, iv, v, ov)
		}
	}
}
