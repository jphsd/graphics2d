package graphics2d

// Transform interface provides the Apply function to transform a set of points.
type Transform interface {
	Apply(points ...[]float64) [][]float64
}
