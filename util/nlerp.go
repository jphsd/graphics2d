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

func (s *NLLinear) Transform(t float64) float64 {
	return t
}

func (s *NLLinear) InvTransform(v float64) float64 {
	return v
}

type NLSquare struct{}

func (s *NLSquare) Transform(t float64) float64 {
	return t * t
}

func (s *NLSquare) InvTransform(v float64) float64 {
	return math.Sqrt(v)
}

type NLCube struct{}

func (s *NLCube) Transform(t float64) float64 {
	return t * t * t
}

func (s *NLCube) InvTransform(v float64) float64 {
	return math.Pow(v, 1/3.0)
}

type NLExponential struct {
	k     float64
	scale float64
}

func NewNLExponential(k float64) *NLExponential {
	return &NLExponential{k, 1 / (math.Exp(k) - 1)}
}

func (s *NLExponential) Transform(t float64) float64 {
	return (math.Exp(t*s.k) - 1) * s.scale
}

func (s *NLExponential) InvTransform(v float64) float64 {
	return math.Log1p(v/s.scale) / s.k
}

type NLLogarithmic struct {
	k     float64
	scale float64
}

func NewNLLogarithmic(k float64) *NLLogarithmic {
	return &NLLogarithmic{k, 1 / math.Log1p(k)}
}

func (s *NLLogarithmic) Transform(t float64) float64 {
	return math.Log1p(t*s.k) * s.scale
}

func (s *NLLogarithmic) InvTransform(v float64) float64 {
	return (math.Exp(v/s.scale) - 1) / s.k
}

type NLSin struct{} // first derivative 0 at t=0,1

// Range [-Pi/2,Pi/2]
func (s *NLSin) Transform(t float64) float64 {
	return (math.Sin((t-0.5)*math.Pi) + 1) / 2
}

func (s *NLSin) InvTransform(v float64) float64 {
	return math.Asin((v*2)-1)/math.Pi + 0.5
}

type NLCircle struct{}

// Circle bottom right quadrant
func (s *NLCircle) Transform(t float64) float64 {
	return 1 - math.Sqrt(1-t*t)
}

func (s *NLCircle) InvTransform(v float64) float64 {
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

func (s *NLGauss) Transform(t float64) float64 {
	x := s.k * (t - 1)
	x *= -0.5 * x
	return (math.Exp(x) - s.offs) * s.scale
}

func (s *NLGauss) InvTransform(v float64) float64 {
	v /= s.scale
	v += s.offs
	v = math.Log(v)
	v *= -2
	v = math.Sqrt(v)
	return 1 - v/s.k
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

func (s *NLLogistic) Transform(t float64) float64 {
	t = (t - s.mp) * s.k
	return (logisticTransform(t) - s.offs) * s.scale
}

func (s *NLLogistic) InvTransform(v float64) float64 {
	v /= s.scale
	v += s.offs
	v = logisticInvTransform(v)
	return v/s.k + s.mp
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

func (s *NLP3) Transform(t float64) float64 {
	return t * t * (3 - 2*t)
}

func (s *NLP3) InvTransform(v float64) float64 {
	return bsInv(v, s)
}

type NLP5 struct{} // first and second derivatives 0 at t=0,1

func (s *NLP5) Transform(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func (s *NLP5) InvTransform(v float64) float64 {
	return bsInv(v, s)
}

type NLCompound struct {
	nl []NonLinear
}

func NewNLCompound(nl []NonLinear) *NLCompound {
	return &NLCompound{nl}
}

func (s *NLCompound) Transform(t float64) float64 {
	for _, f := range s.nl {
		t = f.Transform(t)
	}

	return t
}

func (s *NLCompound) InvTransform(v float64) float64 {
	for i := len(s.nl) - 1; i > -1; i-- {
		v = s.nl[i].InvTransform(v)
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
