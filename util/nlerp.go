package util

import "math"

// t is clamped to [0,1]
func NLerp(t, start, end float64, f NonLinear) float64 {
	if t < 0 {
		return start
	}
	if t > 1 {
		return end
	}
	t = f.Transform(t)
	return (1-t)*start + t*end
}

// v is clamped to [start, end]
func InvNLerp(v, start, end float64, f NonLinear) float64 {
	t := (v - start) / (end - start)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return f.InvTransform(t)
}

func RemapNL(v, istart, iend, ostart, oend float64, fi, fo NonLinear) float64 {
	return NLerp(InvNLerp(v, istart, iend, fi), ostart, oend, fo)
}

// Float32 versions for Path and x/image/vector
// t is clamped to [0,1]
func NLerp32(t, start, end float32, f NonLinear) float32 {
	if t < 0 {
		return start
	}
	if t > 1 {
		return end
	}
	t = float32(f.Transform(float64(t)))
	return (1-t)*start + t*end
}

// v is clamped to [start, end]
func InvNLerp32(v, start, end float32, f NonLinear) float32 {
	t := (v - start) / (end - start)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return float32(f.InvTransform(float64(t)))
}

func RemapNL32(v, istart, iend, ostart, oend float32, fi, fo NonLinear) float32 {
	return NLerp32(InvNLerp32(v, istart, iend, fi), ostart, oend, fo)
}

// For mapping 0 -> 1 non-linearly
// No checks! Only valid in range [0,1]
type NonLinear interface {
	Transform(t float64) float64
	InvTransform(v float64) float64
}

type NLLinear struct{}

func (nl *NLLinear) Transform(t float64) float64 {
	return t
}

func (nl *NLLinear) InvTransform(v float64) float64 {
	return v
}

type NLSquare struct{}

func (nl *NLSquare) Transform(t float64) float64 {
	return t * t
}

func (nl *NLSquare) InvTransform(v float64) float64 {
	return math.Sqrt(v)
}

type NLCube struct{}

func (nl *NLCube) Transform(t float64) float64 {
	return t * t * t
}

func (nl *NLCube) InvTransform(v float64) float64 {
	return math.Pow(v, 1/3.0)
}

type NLExponential struct {
	k     float64
	scale float64
}

func NewNLExponential(k float64) *NLExponential {
	return &NLExponential{k, 1 / (math.Exp(k) - 1)}
}

func (nl *NLExponential) Transform(t float64) float64 {
	return (math.Exp(t*nl.k) - 1) * nl.scale
}

func (nl *NLExponential) InvTransform(v float64) float64 {
	return math.Log1p(v/nl.scale) / nl.k
}

type NLLogarithmic struct {
	k     float64
	scale float64
}

func NewNLLogarithmic(k float64) *NLLogarithmic {
	return &NLLogarithmic{k, 1 / math.Log1p(k)}
}

func (nl *NLLogarithmic) Transform(t float64) float64 {
	return math.Log1p(t*nl.k) * nl.scale
}

func (nl *NLLogarithmic) InvTransform(v float64) float64 {
	return (math.Exp(v/nl.scale) - 1) / nl.k
}

type NLSin struct{} // first derivative 0 at t=0,1

// Range [-Pi/2,Pi/2]
func (nl *NLSin) Transform(t float64) float64 {
	return (math.Sin((t-0.5)*math.Pi) + 1) / 2
}

func (nl *NLSin) InvTransform(v float64) float64 {
	return math.Asin((v*2)-1)/math.Pi + 0.5
}

type NLCircle struct{}

// Circle bottom right quadrant
func (nl *NLCircle) Transform(t float64) float64 {
	return 1 - math.Sqrt(1-t*t)
}

func (nl *NLCircle) InvTransform(v float64) float64 {
	return math.Sqrt(1 - (v-1)*(v-1))
}

type NLGauss struct {
	k, offs, scale float64
}

func NewNLGauss(k float64) *NLGauss {
	offs := math.Exp(-k * k * 0.5)
	scale := 1 / (1 - offs)
	return &NLGauss{k, offs, scale}
}

func (nl *NLGauss) Transform(t float64) float64 {
	x := nl.k * (t - 1)
	x *= -0.5 * x
	return (math.Exp(x) - nl.offs) * nl.scale
}

func (nl *NLGauss) InvTransform(v float64) float64 {
	v /= nl.scale
	v += nl.offs
	v = math.Log(v)
	v *= -2
	v = math.Sqrt(v)
	return 1 - v/nl.k
}

type NLLogistic struct {
	k, mp, offs, scale float64
}

// k > 0 and mp (0,1) - not checked
func NewNLLogistic(k, mp float64) *NLLogistic {
	v0 := -mp * k
	v0 = logisticTransform(v0)
	v1 := (1 - mp) * k
	v1 = logisticTransform(v1)
	return &NLLogistic{k, mp, v0, 1 / (v1 - v0)}
}

func (nl *NLLogistic) Transform(t float64) float64 {
	t = (t - nl.mp) * nl.k
	return (logisticTransform(t) - nl.offs) * nl.scale
}

func (nl *NLLogistic) InvTransform(v float64) float64 {
	v /= nl.scale
	v += nl.offs
	v = logisticInvTransform(v)
	return v/nl.k + nl.mp
}

// L = 1, k = 1, mp = 0
func logisticTransform(t float64) float64 {
	return 1 / (1 + math.Exp(-t))
}

// L = 1, k = 1, mp = 0
func logisticInvTransform(v float64) float64 {
	return -math.Log(1/v - 1)
}

type NLP3 struct{} // first derivative 0 at t=0,1

func (nl *NLP3) Transform(t float64) float64 {
	return t * t * (3 - 2*t)
}

func (nl *NLP3) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

type NLP5 struct{} // first and second derivatives 0 at t=0,1

func (nl *NLP5) Transform(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func (nl *NLP5) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

type NLCompound struct {
	nl []NonLinear
}

func NewNLCompound(nl []NonLinear) *NLCompound {
	return &NLCompound{nl}
}

func (nl *NLCompound) Transform(t float64) float64 {
	for _, f := range nl.nl {
		t = f.Transform(t)
	}

	return t
}

func (nl *NLCompound) InvTransform(v float64) float64 {
	for i := len(nl.nl) - 1; i > -1; i-- {
		v = nl.nl[i].InvTransform(v)
	}
	return v
}

// Numerical method to find inverse
func bsInv(v float64, f NonLinear) float64 {
	n := 16
	t := 0.5
	s := 0.25

	for ; n > 0; n-- {
		if f.Transform(t) > v {
			t -= s
		} else {
			t += s
		}
		s /= 2
	}
	return t
}
