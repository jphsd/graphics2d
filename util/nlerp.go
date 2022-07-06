package util

import (
	"math"
	"math/rand"
)

/*
 * Non-linear interpolations between 0 and 1. Clamping is enforced lest the result not
 * be defined outside of [0,1]
 */

// NLerp returns the value of the supplied non-linear function at t. Note t is clamped to [0,1]
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

// InvNLerp performs the inverse of NLerp and returns the value of t for a value v (clamped to [start, end]).
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

// RemapNL converts v from one space to another by applying InvNLerp to find t in the initial range, and
// then using t to find v' in the new range.
func RemapNL(v, istart, iend, ostart, oend float64, fi, fo NonLinear) float64 {
	return NLerp(InvNLerp(v, istart, iend, fi), ostart, oend, fo)
}

// NLerp32 is a float32 version of NLerp for Path and x/image/vector
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

// InvNLerp32 is a float32 version of InvNLerp for Path and x/image/vector
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

// RemapNL32 is a float32 version of RemapNL for Path and x/image/vector
func RemapNL32(v, istart, iend, ostart, oend float32, fi, fo NonLinear) float32 {
	return NLerp32(InvNLerp32(v, istart, iend, fi), ostart, oend, fo)
}

// NonLinear interface defines the transform and its inverse as used by NLerp etc.
// For mapping 0 -> 1 non-linearly. No checks! Only valid in range [0,1]
type NonLinear interface {
	Transform(t float64) float64
	InvTransform(v float64) float64
}

// NLLinear v = t
type NLLinear struct{}

func (nl *NLLinear) Transform(t float64) float64 {
	return t
}

func (nl *NLLinear) InvTransform(v float64) float64 {
	return v
}

// NLSquare v = t^2
type NLSquare struct{}

func (nl *NLSquare) Transform(t float64) float64 {
	return t * t
}

func (nl *NLSquare) InvTransform(v float64) float64 {
	return math.Sqrt(v)
}

// NLCube v = t^3
type NLCube struct{}

func (nl *NLCube) Transform(t float64) float64 {
	return t * t * t
}

func (nl *NLCube) InvTransform(v float64) float64 {
	return math.Pow(v, 1/3.0)
}

// NLExponential v = (exp(t*k) - 1) * scale
type NLExponential struct {
	K     float64
	Scale float64
}

func NewNLExponential(k float64) *NLExponential {
	return &NLExponential{k, 1 / (math.Exp(k) - 1)}
}

func (nl *NLExponential) Transform(t float64) float64 {
	return (math.Exp(t*nl.K) - 1) * nl.Scale
}

func (nl *NLExponential) InvTransform(v float64) float64 {
	return math.Log1p(v/nl.Scale) / nl.K
}

// NLLogarithmic v = log(1+t*k) * scale
type NLLogarithmic struct {
	K     float64
	Scale float64
}

func NewNLLogarithmic(k float64) *NLLogarithmic {
	return &NLLogarithmic{k, 1 / math.Log1p(k)}
}

func (nl *NLLogarithmic) Transform(t float64) float64 {
	return math.Log1p(t*nl.K) * nl.Scale
}

func (nl *NLLogarithmic) InvTransform(v float64) float64 {
	return (math.Exp(v/nl.Scale) - 1) / nl.K
}

// NLSin v = sin(t) with t mapped to [-Pi/2,Pi/2]
type NLSin struct{} // first derivative 0 at t=0,1

func (nl *NLSin) Transform(t float64) float64 {
	return (math.Sin((t-0.5)*math.Pi) + 1) / 2
}

func (nl *NLSin) InvTransform(v float64) float64 {
	return math.Asin((v*2)-1)/math.Pi + 0.5
}

// NLSin1 v = sin(t) with t mapped to [0,Pi/2]
type NLSin1 struct{} // first derivative 0 at t=1

func (nl *NLSin1) Transform(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

func (nl *NLSin1) InvTransform(v float64) float64 {
	return math.Asin(v) / math.Pi * 2
}

// NLSin2 v = sin(t) with t mapped to [-Pi/2,0]
type NLSin2 struct{} // first derivative 0 at t=0,1

func (nl *NLSin2) Transform(t float64) float64 {
	return math.Sin((t-1)*math.Pi/2) + 1
}

func (nl *NLSin2) InvTransform(v float64) float64 {
	return math.Asin(v-1)*2/math.Pi + 1
}

// NLCircle1 v = 1 - sqrt(1-t^2)
type NLCircle1 struct{}

func (nl *NLCircle1) Transform(t float64) float64 {
	return 1 - math.Sqrt(1-t*t)
}

func (nl *NLCircle1) InvTransform(v float64) float64 {
	return math.Sqrt(1 - (v-1)*(v-1))
}

// NLCircle2 v = sqrt(2t-t^2)
type NLCircle2 struct{}

func (nl *NLCircle2) Transform(t float64) float64 {
	return math.Sqrt(t * (2 - t))
}

func (nl *NLCircle2) InvTransform(v float64) float64 {
	return 1 - math.Sqrt(1-v*v)
}

// NLCatenary v = cosh(t)
type NLCatenary struct{}

func (nl *NLCatenary) Transform(t float64) float64 {
	return (math.Cosh(t) - 1) / (math.Cosh(1) - 1)
}

func (nl *NLCatenary) InvTransform(v float64) float64 {
	return math.Acosh(v*(math.Cosh(1)-1) + 1)
}

// NLGauss v = gauss(t, k)
type NLGauss struct {
	K, Offs, Scale float64
}

func NewNLGauss(k float64) *NLGauss {
	offs := math.Exp(-k * k * 0.5)
	scale := 1 / (1 - offs)
	return &NLGauss{k, offs, scale}
}

func (nl *NLGauss) Transform(t float64) float64 {
	x := nl.K * (t - 1)
	x *= -0.5 * x
	return (math.Exp(x) - nl.Offs) * nl.Scale
}

func (nl *NLGauss) InvTransform(v float64) float64 {
	v /= nl.Scale
	v += nl.Offs
	v = math.Log(v)
	v *= -2
	v = math.Sqrt(v)
	return 1 - v/nl.K
}

// NLLogistic v = logistic(t, k, mp)
type NLLogistic struct {
	K, Mp, Offs, Scale float64
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
	t = (t - nl.Mp) * nl.K
	return (logisticTransform(t) - nl.Offs) * nl.Scale
}

func (nl *NLLogistic) InvTransform(v float64) float64 {
	v /= nl.Scale
	v += nl.Offs
	v = logisticInvTransform(v)
	return v/nl.K + nl.Mp
}

// L = 1, k = 1, mp = 0
func logisticTransform(t float64) float64 {
	return 1 / (1 + math.Exp(-t))
}

// L = 1, k = 1, mp = 0
func logisticInvTransform(v float64) float64 {
	return -math.Log(1/v - 1)
}

// NLP3 v = t^2 * (3-2t)
type NLP3 struct{} // first derivative 0 at t=0,1

func (nl *NLP3) Transform(t float64) float64 {
	return t * t * (3 - 2*t)
}

func (nl *NLP3) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

// NLP5 v = t^3 * (t*(6t-15) + 10)
type NLP5 struct{} // first and second derivatives 0 at t=0,1

func (nl *NLP5) Transform(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func (nl *NLP5) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

// NLCompound v = nl[0](nl[1](nl[2](...nl[n-1](t))))
type NLCompound struct {
	Fs []NonLinear
}

func NewNLCompound(fs []NonLinear) *NLCompound {
	return &NLCompound{fs}
}

func (nl *NLCompound) Transform(t float64) float64 {
	for _, f := range nl.Fs {
		t = f.Transform(t)
	}

	return t
}

func (nl *NLCompound) InvTransform(v float64) float64 {
	for i := len(nl.Fs) - 1; i > -1; i-- {
		v = nl.Fs[i].InvTransform(v)
	}
	return v
}

// NLOmt v = 1-f(1-t)
type NLOmt struct {
	F NonLinear
}

func NewNLOmt(f NonLinear) *NLOmt {
	return &NLOmt{f}
}

func (nl *NLOmt) Transform(t float64) float64 {
	return 1 - nl.F.Transform(1-t)
}

func (nl *NLOmt) InvTransform(v float64) float64 {
	return 1 - nl.F.InvTransform(1-v)
}

// NLRand uses random incremental steps from a normal distribution, smoothed with a cubic.
type NLRand struct {
	Steps []float64
	Dt    float64
}

// NewRandNL calculates a collection of ascending values [0,1] with an average increment of mean
// and standard deviation of std.
func NewNLRand(mean, std float64, sharp bool) *NLRand {
	sum := 0.0
	steps := []float64{0, 0} // 0 prefix
	for sum < 1 {
		v := rand.NormFloat64()*std + mean
		if v < 0.00001 {
			v = 0.00001
		}
		sum += v
		if sum > 1 {
			sum = 1
		}
		steps = append(steps, sum)
	}
	if sharp {
		steps[0] = -steps[2]
		steps = append(steps, 2-steps[len(steps)-2])
	} else {
		steps = append(steps, 1.0) // 1 postfix
	}
	return &NLRand{steps, 1.0 / float64(len(steps)-3)}
}

func (nl *NLRand) Transform(t float64) float64 {
	// Find steps t is between and use cubic interpolation
	// to calc v
	fn := math.Floor(t / nl.Dt)
	fr := t - fn*nl.Dt
	frt := fr / nl.Dt
	n := int(fn)
	if n+4 > len(nl.Steps) {
		return 1
	}
	p := nl.Steps[n : n+4]
	return Cubic(frt, p)
}

func (nl *NLRand) InvTransform(v float64) float64 {
	return bsInv(v, nl)
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

// Cubic calculates the value of f(t) for t in range [0,1] given the values of t at -1, 0, 1, 2 in p[]
// fitted to a cubic polynomial: f(t) = at^3 + bt^2 + ct + d. Clamped because it over/undershoots.
func Cubic(t float64, p []float64) float64 {
	v := p[1] + 0.5*t*(p[2]-p[0]+t*(2.0*p[0]-5.0*p[1]+4.0*p[2]-p[3]+t*(3.0*(p[1]-p[2])+p[3]-p[0])))
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	return v
}
