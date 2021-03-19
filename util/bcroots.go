package util

import (
	"fmt"
	"sort"
)

type BezierCurve struct {
	Weights    [][]float64
	WeightsDt  [][]float64
	WeightsDt2 [][]float64
	WeightsDt3 [][]float64
}

func NewBezierCurve(weights [][]float64) *BezierCurve {
	bc := BezierCurve{}
	bc.Weights = weights
	bc.WeightsDt = CalcDerivativeWeights(bc.Weights)
	bc.WeightsDt2 = CalcDerivativeWeights(bc.WeightsDt)
	bc.WeightsDt3 = CalcDerivativeWeights(bc.WeightsDt2)
	return &bc
}

func (bc *BezierCurve) CurveX(t float64) float64 {
	return DeCasteljau(bc.Weights, t)[0]
}

func (bc *BezierCurve) CurveY(t float64) float64 {
	return DeCasteljau(bc.Weights, t)[1]
}

func (bc *BezierCurve) CurveDtX(t float64) float64 {
	return DeCasteljau(bc.WeightsDt, t)[0]
}

func (bc *BezierCurve) CurveDtY(t float64) float64 {
	return DeCasteljau(bc.WeightsDt, t)[1]
}

func (bc *BezierCurve) CurveDt2X(t float64) float64 {
	return DeCasteljau(bc.WeightsDt2, t)[0]
}

func (bc *BezierCurve) CurveDt2Y(t float64) float64 {
	return DeCasteljau(bc.WeightsDt2, t)[1]
}

func (bc *BezierCurve) CurveDt3X(t float64) float64 {
	return DeCasteljau(bc.WeightsDt3, t)[0]
}

func (bc *BezierCurve) CurveDt3Y(t float64) float64 {
	return DeCasteljau(bc.WeightsDt3, t)[1]
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
	for k, _ := range tmap {
		fmt.Sscanf(k, "%f", &res[i])
		i++
	}
	sort.Float64s(res)
	return res
}

func calcRoots(f, df NREval, tmap map[string]bool) {
	// Find roots in range [0,1] via brute force
	dt := 1.0 / 100
	for t := 0.0; t <= 1; t += dt {
		e, r := NRM(t, f, df)
		if e != nil {
			continue
		}
		tmap[fmt.Sprintf("%.4f", r)] = true
	}
}

// NRM is a modified Newton-Raphson root search that bails if t falls outside
// of the range [0,1] since the curve isn't defined there.
func NRM(start float64, f, df func(float64) float64) (error, float64) {
	t := start

	for true {
		d := df(t)
		if Equals(d, 0) {
			return fmt.Errorf("zero derivative at %f", t), 0
		}

		dt := f(t) / d
		if Equals(dt, 0) {
			break
		}
		t = t - dt
		if t < 0 || t > 1 {
			return fmt.Errorf("t %f outside of [0,1]", t), 0
		}
	}
	return nil, t
}
