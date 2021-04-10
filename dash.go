package graphics2d

// Dash contains the dash pattern and offset. The dash pattern represents lengths of pen down, pen up,
// ... and is in the same coordinate system as the path. The offset provides the ability to start from
// anywhere in the pattern.
type Dash struct {
	Snip *Snip
}

// NewDash creates a new dash with the supplied pattern and offset. If the pattern is odd in length
// then it is replicated to create an even length pattern.
func NewDash(pattern []float64, offs float64) *Dash {
	return &Dash{NewSnip(2, pattern, offs)}
}

// Process implements the PathProcessor interface.
func (d *Dash) Process(p *Path) []*Path {
	paths := d.Snip.Process(p)
	np := len(paths)
	dp := np / 2
	res1 := make([]*Path, 0, dp+1)
	res2 := make([]*Path, 0, dp+1)
	for i := 0; i < np; i++ {
		if i%2 == 0 {
			res1 = append(res1, paths[i])
		} else {
			res2 = append(res2, paths[i])
		}
	}
	if d.Snip.State == 0 {
		return res1
	}
	return res2
}
