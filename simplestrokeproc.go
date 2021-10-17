package graphics2d

// SimpleStrokeProc implements the simplest stroke path processor - flatten
// everything to RenderFlatten tolerance and turn each flattened line step into
// a rectangle Width wide. No joins or caps, just lots of little rectangles.
type SimpleStrokeProc struct {
	Width float64
}

// Process implements the PathProcessor interface.
func (cp *SimpleStrokeProc) Process(p *Path) []*Path {
	p = p.Flatten(RenderFlatten)
	d := cp.Width / 2

	// Output a block for every piece of the flattened path
	res := []*Path{}
	parts := p.Parts()
	for _, part := range parts {
		dx, dy := part[1][0]-part[0][0], part[1][1]-part[0][1]
		dx, dy = -dy, dx
		dx, dy = unit(dx, dy)
		dx *= d
		dy *= d
		pp := NewPath([]float64{part[0][0] + dx, part[0][1] + dy})
		pp.AddStep([]float64{part[1][0] + dx, part[1][1] + dy})
		pp.AddStep([]float64{part[1][0] - dx, part[1][1] - dy})
		pp.AddStep([]float64{part[0][0] - dx, part[0][1] - dy})
		pp.Close()
		res = append(res, pp)
	}

	return res
}
