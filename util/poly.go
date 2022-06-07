package util

type wrule int

const (
	OddEven wrule = iota
	NonZero
)

// WindingRule defines how edges in an edge crossing polygon are treated.
var WindingRule = OddEven

// PointInPoly returns true is a point is within the polygon defined by the list of vertices according to the
// setting of WindingRule (OddEven or NonZero).
func PointInPoly(pt []float64, poly ...[]float64) bool {
	n := len(poly)
	if n == 0 {
		return false
	}

	bb := BoundingBox(poly...)
	x, y := pt[0], pt[1]
	if x < bb[0][0] || x > bb[1][0] || y < bb[0][1] || y > bb[1][1] {
		return false
	}

	prev := poly[0]
	poly = append(poly, prev)
	n++
	lsum, rsum := 0, 0
	// Project rays out to left and right counting how many edges are crossed
	for i := 1; i < n; i++ {
		cur := poly[i]
		x0, x1 := prev[0], cur[0]
		y0, y1 := prev[1], cur[1]
		if Equals(y0, y1) {
			// Ignore horizontal edges
			prev = cur
			continue
		}
		up := false
		if y1 < y0 {
			up = true
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}
		// If ray transits through vertex, only want to count once.
		if y0 < y && y <= y1 {
			// Calc xi
			t := (y - y0) / (y1 - y0)
			xi := Lerp(t, x0, x1)
			if xi < x {
				if WindingRule == OddEven {
					lsum++
				} else if up {
					lsum++
				} else {
					lsum--
				}
			} else if x < xi {
				if WindingRule == OddEven {
					rsum++
				} else if up {
					rsum++
				} else {
					rsum--
				}
			}
		}
		prev = cur
	}

	if WindingRule == OddEven {
		return lsum%2 == 1 && rsum%2 == 1
	}
	return lsum != 0 && rsum != 0
}
