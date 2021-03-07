package util

import "testing"

func TestDistanceToLine(t *testing.T) {
	a := []float64{5, 0}
	b := []float64{0, 4}
	c := []float64{5, 4}
	exp := 11.731707
	d2 := DistanceToLineSquared(a, b, c)
	if !Equals(exp, d2) {
		t.Errorf("Expected %f got %f", exp, d2)
	}
}
