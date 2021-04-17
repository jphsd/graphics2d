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

The Shape type is a container for closed paths and represents something that can be filled and rendered.
As Paths are added, if they're not already closed, they are forced closed.

The PathProcessor interface is where the magic happens. Given a path, a function implementing this
interface returns a collection of paths derived from it. This allows for stroking, dashing and a variety
of other possibilities:
  CapsProc - adds shapes at the start and end of a path
  CompoundProc - allows path processors to be run in sequence
  CurvesToLinesProc - replaces curved steps with lines between the points
  DashProc - wraps SnipProc to produce a dashed path
  FlattenProc - wraps Path.Flatten
  LineProc - reduces a path to a single line from start to end
  MunchProc - converts a path into length sized line paths
  PointsProc - adds shapes at the start of each step and the end of the last step
  SimplifyProc - wraps Path.Simplify
  ShapesProc - distributes shapes along a path separated by some distance
  SnipProc - cuts up a path into smaller pieces according to a pattern
  SplitProc - splits each step into its own path
  StrokeProc - creates a fixed width outline of a path with options for the cap and
  join styles

Shapes and Paths are be rendered with the render functions. Paths are forced closed when rendered.
Convenience methods are provided for rendering with a single color or an image. The full render function
allows a clip mask and offset to be supplied and the draw.Op to be specified.

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
