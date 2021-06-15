package surface

// NormalMap provides an At function to determine the unit normal at a location.
type NormalMap interface {
	At(x, y int) []float64
}

// DefaultNM is a map where all normals are {0, 0, 1}.
type DefaultNM struct{}

// At implements the At function of the NormalMap interface.
func (n *DefaultNM) At(x, y int) []float64 {
	return []float64{0, 0, 1}
}
