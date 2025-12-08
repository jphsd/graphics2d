package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// FSnipProc contains the snip pattern and offset. The snip pattern represents lengths of state0, state1,
// ... stateN-1, and is in the same coordinate system as the path. The offset provides the ability to
// start from anywhere in the pattern.
type FSnipProc struct {
	N       int
	Pattern []float64
	Flatten float64
	State   int
	length  float64
	patind  int
	delta   float64
}

// SnipProc preserves the curve order of the parts. FSnipProc flattens the path so all parts are linear.

// NewFSnipProc creates a new snip path processor with the supplied pattern and offset. If the pattern is
// not N in length then it is replicated to create a mod N length pattern.
func NewFSnipProc(n int, pattern []float64, offs float64) *FSnipProc {
	pat := pattern[:]
	for len(pat)%n != 0 {
		pat = append(pat, pattern...)
	}
	s := sum(pat)

	res := &FSnipProc{n, pat, RenderFlatten, 0, s, 0, 0}
	res.Offset(offs)
	return res
}

// Offset determines where in the pattern the path processor will start.
func (sp *FSnipProc) Offset(offs float64) {
	neg := offs < 0
	if neg {
		offs = -offs
	}
	f := offs / sp.length
	if f > 1 {
		f = math.Floor(f)
		offs -= f * sp.length
	}
	if neg {
		offs = sp.length - offs
	}

	// Figure out initial state, pattern index and delta based on offset.
	state := 0
	patind := 0                 // which part of the pattern we're on
	delta := sp.Pattern[patind] // distance to the next state change

	// Walk to offset in pattern
	for true {
		if offs > delta {
			offs -= delta
			patind++
			state++
			if state == sp.N {
				state = 0
			}
			delta = sp.Pattern[patind]
			continue
		}
		// offs is > 0 and < delta
		delta -= offs
		break
	}

	if util.Equals(delta, 0) {
		state++
		if state == sp.N {
			state = 0
		}
		patind++
		if patind == len(sp.Pattern) {
			patind = 0
		}
		delta = sp.Pattern[patind]
	}

	sp.State = state
	sp.patind = patind
	sp.delta = delta
}

// Process implements the PathProcessor interface.
func (sp *FSnipProc) Process(p *Path) []*Path {
	// Check path isn't a point
	if len(p.Steps()) == 0 {
		return []*Path{p}
	}

	// Flatten the path and get the parts as lines
	path := p.Flatten(sp.Flatten).Process(&StepsToLinesProc{false})[0]
	parts := path.Parts()

	// Pattern state variables
	patind := sp.patind
	delta := sp.delta

	// Part state variables
	pi := 0
	cur, end := parts[pi][0], parts[pi][1]
	avail := math.Hypot(end[0]-cur[0], end[1]-cur[1])
	sparts := []Part{}

	// Turn parts into snip paths
	res := []*Path{}
	for {
		// Collect parts until we have met the delta
		// unless we run out of parts, in which case we're done
		// Once delta met, add path to res, bump State, patind and delta
		// and continue
		if avail < delta {
			sparts = append(sparts, Part{cur, end})
			delta -= avail
			pi++
			if pi > len(parts)-1 {
				return append(res, PartsToPath(sparts...))
			}
			cur = end
			end = parts[pi][1]
			avail = math.Hypot(end[0]-cur[0], end[1]-cur[1])
		} else {
			// avail >= delta
			t := delta / avail
			omt := 1 - t
			tmp := []float64{cur[0]*omt + end[0]*t, cur[1]*omt + end[1]*t}
			sparts = append(sparts, Part{cur, tmp})
			avail -= delta
			res = append(res, PartsToPath(sparts...))
			sparts = nil
			patind++
			if patind == len(sp.Pattern) {
				patind = 0
			}
			delta = sp.Pattern[patind]
			if util.Equals(avail, 0) {
				pi++
				if pi > len(parts)-1 {
					return res
				}
				cur = end
				end = parts[pi][1]
				avail = math.Hypot(end[0]-cur[0], end[1]-cur[1])
			} else {
				cur = tmp
			}
		}
	}

	return res
}

// DashProc contains the dash pattern and offset. The dash pattern represents lengths of pen down, pen up,
// ... and is in the same coordinate system as the path. The offset provides the ability to start from
// anywhere in the pattern.
type DashProc struct {
	Snip *FSnipProc
}

// NewDashProc creates a new dash path processor with the supplied pattern and offset. If the pattern is
// odd in length then it is replicated to create an even length pattern.
func NewDashProc(pattern []float64, offs float64) DashProc {
	return DashProc{NewFSnipProc(2, pattern, offs)}
}

// Process implements the PathProcessor interface.
func (d DashProc) Process(p *Path) []*Path {
	paths := d.Snip.Process(p)
	np := len(paths)
	dp := np / 2
	res1 := make([]*Path, 0, dp+1)
	res2 := make([]*Path, 0, dp+1)
	for i := range np {
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

// MunchProc contains the munching compound path processor.
type MunchProc struct {
	Comp CompoundProc
}

// NewMunchProc creates a munching path processor. It calculates points along a path spaced l apart
// and creates new paths that join the points with lines.
func NewMunchProc(l float64) MunchProc {
	if l < 0 {
		l = -l
	}

	return MunchProc{NewCompoundProc(NewFSnipProc(2, []float64{l, l}, 0), PathToLineProc{})}
}

// Process implements the PathProcessor interface.
func (mp MunchProc) Process(p *Path) []*Path {
	return p.Process(mp.Comp)
}
