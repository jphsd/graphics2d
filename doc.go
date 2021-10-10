/*
Package graphics2d contains types and functions for 2D graphics rendering.

If all you want to do is render a point, line, arc or path of a particular color into an image, then
use the draw functions which hide all of the mechanics described below from you:
  DrawPoint
  DrawLine
  DrawArc
  DrawPath

However, to take full advantage of the package, here's a few more things you'll need to know.

The lowest type is a Path. Paths start from a location and add steps using a number of control points. The
number of control points determines the polynomial order of the step: 0 - line; 1 - quadratic; 2 - cubic;
3 - quartic; ... A path can be closed which forces a line from the start location to the last. Once a path
is closed, no more steps can be added. Unlike other implementations, a path here represents a single
stroke (pen down). There is no move (pen up) step.

The Shape type is a container for paths and represents something that can be filled and rendered.
When paths are rendered, they must be fillable (i.e. closed), so they are forced closed. This can lead to
unexpected results... If a shape is rendered with a pen (see DrawShape) and the pen has a stroke
associated with it, then open paths are not an issue since strokes return closed paths.

A pen is a combination of a color (or image), a stroke and a transformation from shape space to image
space. If you're defining operations in the same space as the image, then the transformation is
simply the identity transformation. If, however, say your operations are in a space [0, 1]^2, then
you'd specify a transformation that maps [0, 1] => [0, width-1] etc. (i.e. scale by image width and
height). Note that this transforamtion is applied *after* the stroke, so for a 1 pixel wide stroke
the stroke width would be 1 / width. You don't have to use pens, you can use the RenderShape, PathProcessor,
image filler and transformation functions directly. Pens just provide a convenient abstraction.

Note that if the image is written to every graphics operation (as it is with the Draw*() functions), this
will kill performance as the entire image is written every time. It's better to collect all the operations
associated with a color in a shape and then render that shape once.

The PathProcessor interface is where the magic happens. Given a path, a function implementing this
interface returns a collection of paths derived from it. This allows for stroking, dashing and a variety
of other possibilities:
  CapsProc - adds shapes at the start and end of a path
  CompoundProc - allows multiple path processors to be run in sequence
  CurvesToLinesProc - replaces curved steps with lines between the points
  DashProc - wraps SnipProc to produce a dashed path
  FlattenProc - wraps Path.Flatten
  JitterProc - randomly move segment endpoints by some percentage of the segment's length
  LineProc - wraps Path.Line
  MunchProc - converts a path into length sized line paths
  OpenProc - wraps Path.Open
  PointsProc - adds shapes at the start of each step and the end of the last step
  ReverseProc - wraps Path.Reverse
  SimplifyProc - wraps Path.Simplify
  ShapesProc - distributes shapes along a path separated by some distance
  SnipProc - cuts up a path into smaller pieces according to a pattern
  SplitProc - splits each step into its own path
  StrokeProc - creates a fixed width outline of a path with options for the cap and
  join styles
  TraceProc - creates a new path by tracing the normals of the path at a fixed distance
  TransformProc - wraps Path.Transform

Shapes are rendered with the render functions. Paths are forced closed when rendered (see shapes
above). Convenience methods are provided for rendering with a single color or an image (see also pens,
above). The full render function allows a clip mask and offset to be supplied, and the draw.Op to be specified.

The Aff3 type provides the ability to specify affine transforms on Paths and Shapes.

Utility functions are provided to generate common forms of paths:
  Point
  Line
  PolyLine
  Curve
  PolyCurve
  Arc
  ArcFromPoint
  PolyArcFromPoint
  Circle
  Ellipse
  EllipticalArc
  EllipticalArcFromPoint
  RegularPolygon

A shape function is provided to capture glyphs:
  GlyphToShape
*/
package graphics2d
