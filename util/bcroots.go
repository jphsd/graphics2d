package util

import (
	"fmt"
	"math"
	"sort"
)

// BezierCurve stores the derivative curve weights.
type BezierCurve struct {
	Weights    [][]float64
	WeightsDt  [][]float64
	WeightsDt2 [][]float64
	WeightsDt3 [][]float64
}

// NewBezierCurve creates a new BezierCurve with the weights of the first, second
// and third order derivatives of the supplied curve.
func NewBezierCurve(weights [][]float64) *BezierCurve {
	bc := BezierCurve{}
	bc.Weights = weights
	bc.WeightsDt = CalcDerivativeWeights(bc.Weights)
	bc.WeightsDt2 = CalcDerivativeWeights(bc.WeightsDt)
	bc.WeightsDt3 = CalcDerivativeWeights(bc.WeightsDt2)
	return &bc
}

// CurveX returns the X value for the curve at t.
func (bc *BezierCurve) CurveX(t float64) float64 {
	return DeCasteljau(bc.Weights, t)[0]
}

// CurveY returns the Y value for the curve at t.
func (bc *BezierCurve) CurveY(t float64) float64 {
	return DeCasteljau(bc.Weights, t)[1]
}

// CurveDtX returns the X value for the derivative of the curve at t.
func (bc *BezierCurve) CurveDtX(t float64) float64 {
	return DeCasteljau(bc.WeightsDt, t)[0]
}

// CurveDtY returns the Y value for the derivative of the curve at t.
func (bc *BezierCurve) CurveDtY(t float64) float64 {
	return DeCasteljau(bc.WeightsDt, t)[1]
}

// CurveDt2X returns the X value for the second order derivative of the curve at t.
func (bc *BezierCurve) CurveDt2X(t float64) float64 {
	return DeCasteljau(bc.WeightsDt2, t)[0]
}

// CurveDt2Y returns the Y value for the second order derivative of the curve at t.
func (bc *BezierCurve) CurveDt2Y(t float64) float64 {
	return DeCasteljau(bc.WeightsDt2, t)[1]
}

// CurveDt3X returns the X value for the third order derivative of the curve at t.
func (bc *BezierCurve) CurveDt3X(t float64) float64 {
	return DeCasteljau(bc.WeightsDt3, t)[0]
}

// CurveDt3Y returns the Y value for the third order derivative of the curve at t.
func (bc *BezierCurve) CurveDt3Y(t float64) float64 {
	return DeCasteljau(bc.WeightsDt3, t)[1]
}

// Kappa calculates the curvature at t. Radius of curvature at t is 1/kappa(t)
func (bc *BezierCurve) Kappa(t float64) float64 {
	dpt := DeCasteljau(bc.WeightsDt, t)
	d2pt := DeCasteljau(bc.WeightsDt2, t)
	return (dpt[0]*d2pt[1] - d2pt[0]*dpt[1]) / math.Pow(dpt[0]*dpt[0] + dpt[1]*dpt[1], 1.5)
}

// CalcExtremities finds the extremes of a curve in terms of t.
func CalcExtremities(points [][]float64) []float64 {
	n := len(points)

	if n == 2 {
		return []float64{0, 1}
	}

	bc := NewBezierCurve(points)
	tmap := make(map[string]bool) // Use "%.4f"
	tmap["0.0000"], tmap["1.0000"] = true, true

	// Find local minima and maxima with Dt and Dt2
	calcRoots(bc.CurveDtX, bc.CurveDt2X, tmap)
	calcRoots(bc.CurveDtY, bc.CurveDt2Y, tmap)

	if n > 3 {
		// Find inflection points with Dt2 and Dt3
		calcRoots(bc.CurveDt2X, bc.CurveDt3X, tmap)
		calcRoots(bc.CurveDt2Y, bc.CurveDt3Y, tmap)
	}

	// Convert t values back to float64
	res := make([]float64, len(tmap))
	i := 0
	for k := range tmap {
		fmt.Sscanf(k, "%f", &res[i])
		i++
	}
	sort.Float64s(res)
	return res
}

func calcRoots(f, df func(float64) float64, tmap map[string]bool) {
	// Find roots in range [0,1] via brute force
	dt := 1.0 / 100
	for t := 0.0; t <= 1; t += dt {
		r, e := NRM(t, f, df)
		if e != nil {
			continue
		}
		tmap[fmt.Sprintf("%.4f", r)] = true
	}
}

// NRM is a modified Newton-Raphson root search that bails if t falls outside
// of the range [0,1] since the curve isn't defined there.
func NRM(start float64, f, df func(float64) float64) (float64, error) {
	t := start

	for true {
		d := df(t)
		if Equals(d, 0) {
			return 0, fmt.Errorf("zero derivative at %f", t)
		}

		dt := f(t) / d
		if Equals(dt, 0) {
			break
		}
		t = t - dt
		if t < 0 || t > 1 {
			return 0, fmt.Errorf("t %f outside of [0,1]", t)
		}
	}
	return t, nil
}
