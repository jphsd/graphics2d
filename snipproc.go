package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// SnipProc contains the snip pattern and offset. The snip pattern represents lengths of state0, state1,
// ... stateN-1, and is in the same coordinate system as the path. The offset provides the ability to
// start from anywhere in the pattern.
type SnipProc struct {
	N       int
	Pattern []float64
	Flatten float64
	State   int
	length  float64
	patind  int
	delta   float64
}

// NewSnipProc creates a new snip path processor with the supplied pattern and offset. If the pattern is
// not N in length then it is replicated to create a mod N length pattern.
func NewSnipProc(n int, pattern []float64, offs float64) *SnipProc {
	pat := pattern[:]
	for len(pat)%n != 0 {
		pat = append(pat, pattern...)
	}
	s := sum(pat)

	// Default flattening value is 1
	res := &SnipProc{n, pat, 1, 0, s, 0, 0}
	res.Offset(offs)
	return res
}

func sum(l []float64) float64 {
	var s float64 = 0
	for _, v := range l {
		s += v
	}
	return s
}

// Offset determines where in the pattern the path processor will start.
func (sp *SnipProc) Offset(offs float64) {
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
func (sp *SnipProc) Process(p *Path) []*Path {
	// Flatten the path parts
	parts := p.Parts()
	np := len(parts)
	if np == 0 {
		return []*Path{p}
	}
	fparts := make([][][][]float64, np) // part:subparts:points:xy
	for i, part := range parts {
		fparts[i] = FlattenPart(sp.Flatten, part)
	}
	lparts := getLengths(fparts)

	// Use flattened parts to build list of parts and their t values where there's a state change
	patind := sp.patind
	delta := sp.delta
	chind := []int{}
	cht := []float64{}

	for i, lpart := range lparts {
		// last lpart contains the length of this part
		nsegs := len(lpart) - 1
		partlen := lpart[nsegs][0] // total length of this part
		if partlen < delta {
			// skip this part
			delta -= partlen
			continue
		}
		// walk individual line segs until we find the one with the state change in it
		for j := 0; j < nsegs; j++ {
			length := lpart[j][1]
			if length < delta {
				// skip this seg
				delta -= length
				continue
			}
			lsum := 0.0
			for length > delta {
				rem := length - delta
				vlen := lpart[j][0] + lsum + delta
				lsum += delta

				tlen := vlen / partlen
				chind = append(chind, i)
				cht = append(cht, tlen)

				patind++
				if patind == len(sp.Pattern) {
					patind = 0
				}
				delta = sp.Pattern[patind]
				length = rem
			}
			delta -= length
		}
	}

	npp := len(chind)
	if npp == 0 {
		// All of path is in the snip
		return []*Path{p.Open()}
	}

	// Build snipped paths based on part indices and t values as split points.
	// This way we preserve the original curves.
	cht = convTVals(chind, cht)
	res := make([]*Path, 0, npp+1)
	rem := parts[0]           // part remaining after last split
	pind := 0                 // current index into parts
	pparts := [][][]float64{} // parts collected towards next path

	for i := 0; i < npp; i++ {
		p, t := chind[i], cht[i]
		for p > pind {
			pparts = append(pparts, rem)
			pind++
			rem = parts[pind]
		}
		// rem is the correct part
		if util.Equals(t, 0) {
			// t == 0
			if len(pparts) == 0 {
				continue
			}
			lp := PartsToPath(pparts...)
			res = append(res, lp)
			// rem already set
		} else {
			// state change is in this part, split it at t
			pieces := util.SplitCurve(rem, t)
			pparts = append(pparts, pieces[0])
			lp := PartsToPath(pparts...)
			res = append(res, lp)
			rem = pieces[1]
		}
		pparts = [][][]float64{}
	}

	// Handle remaining path
	pparts = append(pparts, rem)
	pind++
	for pind < np {
		pparts = append(pparts, parts[pind])
		pind++
	}
	lp := PartsToPath(pparts...)
	return append(res, lp)
}

// part:lineseg:cumpartlen/len
func getLengths(parts [][][][]float64) [][][]float64 {
	res := make([][][]float64, len(parts))
	for i, part := range parts {
		sumPart := 0.0
		n := len(part) + 1
		res[i] = make([][]float64, n)
		for j, lineseg := range part {
			dx, dy := lineseg[1][0]-lineseg[0][0], lineseg[1][1]-lineseg[0][1]
			len := math.Sqrt(dx*dx + dy*dy)
			res[i][j] = []float64{sumPart, len}
			sumPart += len
		}
		// Last chunk in part is the part total
		res[i][n-1] = []float64{sumPart}
	}
	return res
}

// t vals need to be relative to the remaining part after splitting.
// e.g. part split at 0.25, 0.5 and 0.666 => 0.25, 0.333 and .333
func convTVals(chind []int, cht []float64) []float64 {
	n := len(chind)
	res := make([]float64, n)
	p := -1
	lt := 0.0
	ll := 1.0
	for i := 0; i < n; i++ {
		if chind[i] != p {
			// Reset
			lt = cht[i]
			ll = 1 - lt
			res[i] = lt
			p = chind[i]
			continue
		}
		t := cht[i]
		d := t - lt
		res[i] = d / ll
		lt = t
		ll -= d
	}
	return res
}

// DashProc contains the dash pattern and offset. The dash pattern represents lengths of pen down, pen up,
// ... and is in the same coordinate system as the path. The offset provides the ability to start from
// anywhere in the pattern.
type DashProc struct {
	Snip *SnipProc
}

// NewDashProc creates a new dash path processor with the supplied pattern and offset. If the pattern is
// odd in length then it is replicated to create an even length pattern.
func NewDashProc(pattern []float64, offs float64) *DashProc {
	return &DashProc{NewSnipProc(2, pattern, offs)}
}

// Process implements the PathProcessor interface.
func (d *DashProc) Process(p *Path) []*Path {
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

// MunchProc contains the length.
type MunchProc struct {
	Munch *CompoundProc
}

// NewMunchProc creates a munching path processor. It calculates points along a path spaced l apart
// and creates new paths that join the points with lines.
func NewMunchProc(l float64) *MunchProc {
	if l < 0 {
		l = -l
	}

	return &MunchProc{NewCompoundProc(NewSnipProc(2, []float64{l, l}, 0), &LineProc{})}
}

// Process implements the PathProcessor interface.
func (mp *MunchProc) Process(p *Path) []*Path {
	return p.Process(mp.Munch)
}
