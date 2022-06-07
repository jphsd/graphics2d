package util

import "testing"

func TestDistanceToLine(t *testing.T) {
	a := []float64{5, 0}
	b := []float64{0, 5}
	c := []float64{5, 5}
	d2, ip, it := DistanceToLineSquared(a, b, c)
	exp := 12.5
	if !Equals(exp, d2) {
		t.Errorf("Expected d2 %f got %f", exp, d2)
	}
	exp = 2.5
	if !Equals(exp, ip[0]) {
		t.Errorf("Expected %f got %f", exp, ip[0])
	}
	if !Equals(exp, ip[1]) {
		t.Errorf("Expected %f got %f", exp, ip[1])
	}
	exp = 0.5
	if !Equals(exp, it) {
		t.Errorf("Expected %f got %f", exp, it)
	}
}
