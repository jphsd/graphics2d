# Yet Another 2D Graphics Package For Go
[![Go Reference](https://pkg.go.dev/badge/github.com/jphsd/graphics2d.svg)](https://pkg.go.dev/github.com/jphsd/graphics2d)
[![godocs.io](http://godocs.io/github.com/jphsd/graphics2d?status.svg)](http://godocs.io/github.com/jphsd/graphics2d)
[![Go Report Card](https://goreportcard.com/badge/github.com/jphsd/graphics2d)](https://goreportcard.com/report/github.com/jphsd/graphics2d)
[![Build Status](https://travis-ci.com/jphsd/graphics2d.svg?branch=master)](https://travis-ci.com/github/jphsd/graphics2d)

Still very much under development although the image functions are about as complete as I'm going to make them.
Note - the go:embed feature from v1.16 is required to build and go vet will fail for go versions < 1.17 due
to the new way of handling build tags.

The top level Path and Shape types are complete, and the majority of PathProcessors implemented, including:
- StrokeProc - fixed width strokes with a variety of cap and join types.
- SnipProc - chops up a path according to a pattern
- DashProc - wrapper around SnipProc for creating a dashed path
- CompoundProc - allows concatenation of PathProcessors
> dashedstroke := NewCompoundProc(NewDashProc(pattern, offs), NewStrokeProc(1))

Wiki entries [here](https://github.com/jphsd/graphics2d/wiki)

## Code Structure
![Visualization of the codebase](./diagram.svg)
