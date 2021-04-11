# Yet Another 2D Graphics Package For Go

Still very much under development although the image functions are about as complete as I'm going to make them.

The top level Path and Shape types are complete, and the following PathProcessors implemented:
- Stroke - fixed width strokes with a variety of cap and join types.
- Snip - chops up a path according to a pattern
- Dash - wrapper around Snip for creating a dashed path
- CompoundProcessor - allows concatenation of PathProcessors
> dashedstroke := NewCompoundProcessor(NewDash(pattern, offs), NewStroke(1))

Package documentation [here](https://pkg.go.dev/github.com/jphsd/graphics2d)

Wiki entries [here](https://github.com/jphsd/graphics2d/wiki)
