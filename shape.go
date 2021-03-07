package graphics2d

import "image"

/*
 * A Shape is a fillable collection of paths. For a path to be fillable,
 * it must be closed, so paths added to the shape are forced closed.
 */

type Shape struct {
	paths  []*Path
	bounds image.Rectangle
}

func (s *Shape) Bounds() image.Rectangle {
	if s.bounds.Empty() && s.paths != nil {
		rect := s.paths[0].Bounds()
		for i := 1; i < len(s.paths); i++ {
			rect = rect.Union(s.paths[i].Bounds())
		}
		s.bounds = rect
	}
	return s.bounds
}

func (s *Shape) AddPath(p *Path) {
	lp := p.Copy()
	lp.Close()
	if s.paths == nil {
		s.paths = make([]*Path, 1)
		s.paths[0] = lp
	} else {
		s.paths = append(s.paths, lp)
	}
	s.bounds = image.Rectangle{}
}

func (s *Shape) AddPaths(paths []*Path) {
	for _, p := range paths {
		s.AddPath(p)
	}
}

func (s *Shape) AddShape(shape *Shape) {
	s.AddPaths(shape.Paths())
}

func (s *Shape) Paths() []*Path {
	return s.paths[:]
}

func (s *Shape) Copy() *Shape {
	np := make([]*Path, len(s.paths))
	copy(np, s.paths)
	return &Shape{np, s.bounds}
}

// Apply an affine transform to all the paths in a shape
// and return a new shape
func (s *Shape) Transform(xfm *Aff3) *Shape {
	np := make([]*Path, len(s.paths))
	for i, path := range s.paths {
		np[i] = path.Transform(xfm)
	}
	return &Shape{np, image.Rectangle{}}
}

// Apply a processor to a shape and force close on
// the processed paths
func (s *Shape) Process(proc PathProcessor) *Shape {
	np := make([]*Path, 0)
	for _, p := range s.paths {
		npaths := p.Process(proc)
		for _, pp := range npaths {
			pp.Close()
			np = append(np, pp)
		}
	}

	return &Shape{np, image.Rectangle{}}
}

// Applies a collection of PathProcessors to a shape
func (s *Shape) CompoundProcess(procs []PathProcessor) *Shape {
	np := make([]*Path, 0)
	for _, p := range s.paths {
		npaths := p.CompoundProcess(procs)
		for _, pp := range npaths {
			pp.Close()
			np = append(np, pp)
		}
	}

	return &Shape{np, image.Rectangle{}}
}
