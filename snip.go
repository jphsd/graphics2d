package graphics2d

import (
	"math"

	. "github.com/jphsd/graphics2d/util"
)

// Snip contains the snip pattern and offset. The snip pattern represents lengths of state0, state1,
// ... stateN-1, and is in the same coordinate system as the path. The offset provides the ability to
// start from anywhere in the pattern.
type Snip struct {
	N       int
	Pattern []float64
	D       float64
	State   int
	length  float64
	patind  int
	delta   float64
}

// NewSnip creates a new snip with the supplied pattern and offset. If the pattern is not N in length
// then it is replicated to create a mod N length pattern.
func NewSnip(n int, pattern []float64, offs float64) *Snip {
	pat := pattern[:]
	for len(pat)%n != 0 {
		pat = append(pat, pattern...)
	}
	s := sum(pat)

	// Default flattening value is 1
	res := &Snip{n, pat, 1, 0, s, 0, 0}
	res.SetOffset(offs)

	/*
		neg := offs < 0
		if neg {
			offs = -offs
		}
		f := offs / s
		if f > 1 {
			f = math.Floor(f)
			offs -= f * s
		}
		if neg {
			offs = s - offs
		}

		// Figure out initial state, pattern index and delta based on offset.
		state := 0
		patind := 0          // which part of the pattern we're on
		delta := pat[patind] // distance to the next state change

		// Walk to offset in pattern
		for true {
			if offs > delta {
				offs -= delta
				patind++
				state++
				if state == n {
					state = 0
				}
				delta = pat[patind]
				continue
			}
			// offs is > 0 and < delta
			delta -= offs
			break
		}

		if Equals(delta, 0) {
			state++
			if state == n {
				state = 0
			}
			patind++
			if patind == len(pat) {
				patind = 0
			}
			delta = pat[patind]
		}

		// Default flattening value is 1
		return &Snip{n, pat, 1, state, s, patind, delta}
	*/
	return res
}

func sum(l []float64) float64 {
	var s float64 = 0
	for _, v := range l {
		s += v
	}
	return s
}

func (s *Snip) SetOffset(offs float64) {
	neg := offs < 0
	if neg {
		offs = -offs
	}
	f := offs / s.length
	if f > 1 {
		f = math.Floor(f)
		offs -= f * s.length
	}
	if neg {
		offs = s.length - offs
	}

	// Figure out initial state, pattern index and delta based on offset.
	state := 0
	patind := 0                // which part of the pattern we're on
	delta := s.Pattern[patind] // distance to the next state change

	// Walk to offset in pattern
	for true {
		if offs > delta {
			offs -= delta
			patind++
			state++
			if state == s.N {
				state = 0
			}
			delta = s.Pattern[patind]
			continue
		}
		// offs is > 0 and < delta
		delta -= offs
		break
	}

	if Equals(delta, 0) {
		state++
		if state == s.N {
			state = 0
		}
		patind++
		if patind == len(s.Pattern) {
			patind = 0
		}
		delta = s.Pattern[patind]
	}

	s.State = state
	s.patind = patind
	s.delta = delta
}

// Process implements the PathProcessor interface.
func (s *Snip) Process(p *Path) []*Path {
	// Flatten the path parts to D
	parts := p.Parts()
	np := len(parts)
	if np == 0 {
		return []*Path{p}
	}
	fparts := make([][][][]float64, np) // part:subparts:points:xy
	for i, part := range parts {
		fparts[i] = FlattenPart(s.D, part)
	}
	lparts := getLengths(fparts)

	// Use flattened parts to build list of parts and their t values where there's a state change
	patind := s.patind
	delta := s.delta
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
				if patind == len(s.Pattern) {
					patind = 0
				}
				delta = s.Pattern[patind]
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
		if Equals(t, 0) {
			// t == 0
			if len(pparts) == 0 {
				continue
			}
			lp, _ := PartsToPath(pparts)
			res = append(res, lp)
			// rem already set
		} else {
			// state change is in this part, split it at t
			pieces := SplitCurve(rem, t)
			pparts = append(pparts, pieces[0])
			lp, _ := PartsToPath(pparts)
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
	lp, _ := PartsToPath(pparts)
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
