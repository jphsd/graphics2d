package main

import (
	"fmt"

	"github.com/jphsd/graphics2d/util"
)

// Print as CSV each non-linear func with 100 steps
func main() {
	table := []struct {
		nl   util.NonLinear
		name string
	}{
		{&util.NLLinear{}, "Linear"},
		{&util.NLSquare{}, "Square"},
		{&util.NLCube{}, "Cube"},
		{&util.NLCircle{}, "Circle"},
		{&util.NLSin{}, "Sin"},
		{util.NewNLExponential(1), "Exponential 1"},
		{util.NewNLExponential(10), "Exponential 2"},
		{util.NewNLExponential(100), "Exponential 3"},
		{util.NewNLLogarithmic(1), "Logarithmic 1"},
		{util.NewNLLogarithmic(10), "Logarithmic 2"},
		{util.NewNLLogarithmic(100), "Logarithmic 3"},
		{util.NewNLGauss(1), "Gauss 1"},
		{util.NewNLGauss(3), "Gauss 2"},
		{util.NewNLGauss(6), "Gauss 3"},
		{util.NewNLLogistic(1, 0.5), "Logistic 1"},
		{util.NewNLLogistic(12, 0.5), "Logistic 2"},
		{util.NewNLLogistic(60, 0.5), "Logistic 3"},
		{util.NewNLLogistic(1, 0.2), "Logistic 4"},
		{util.NewNLLogistic(12, 0.2), "Logistic 5"},
		{util.NewNLLogistic(38, 0.2), "Logistic 6"},
		{util.NewNLLogistic(1, 0.8), "Logistic 7"},
		{util.NewNLLogistic(12, 0.8), "Logistic 8"},
		{util.NewNLLogistic(100, 0.8), "Logistic 9"},
		{&util.NLP3{}, "P3"},
		{&util.NLP5{}, "P5"},
	}

	for t := 0.; t <= 1; t += 1 / 100. {
		var line string
		for _, row := range table {
			v := row.nl.Transform(t)
			line += fmt.Sprintf("%f,", v)
		}
		fmt.Printf("%s\n", line)
	}
}
