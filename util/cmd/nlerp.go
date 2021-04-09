package main

import (
	"fmt"

	. "github.com/jphsd/graphics2d/util"
)

// Print as CSV each non-linear func with 100 steps
func main() {
	table := []struct {
		nl   NonLinear
		name string
	}{
		{&NLLinear{}, "Linear"},
		{&NLSquare{}, "Square"},
		{&NLCube{}, "Cube"},
		{&NLCircle1{}, "Circle"},
		{&NLCircle2{}, "Circle"},
		{&NLSin{}, "Sin"},
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
		{&NLP3{}, "P3"},
		{&NLP5{}, "P5"},
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
