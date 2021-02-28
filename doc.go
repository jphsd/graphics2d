/*
Package graphics2d contains types and functions for 2D graphics rendering.

The lowest type is a Path. Paths start from a location and add steps using a number of control points. The
number of control points determines the polynomial order of the step: 0 - line; 1 - quadratic; c - cubic;
3 - quartic; ... A path can be closed which forces a line from the start location to the last. Once a path
is closed, no more steps can be added. Unlike other implementations, a path here represents a single
stroke (pen down). There is no move (pen up) step.

The Shape type is a container for closed paths and represents something that can be filled and rendered.
As Paths are added, if they're not already closed, they are forced closed.

The PathProcessor interface is where the magic happens. Given a path, a function implementing this interface returns
a collection of paths derived from it. This allows for stoking, dashing and a variety of other possibilities.

Shapes and Paths can be rendered with the render functions. Paths are forced closed when rendered.
Convenience methods are provided for rendering with a single color and an image. The full render function
allows a clip mask to be supplied and the draw.Op to be specified.

The Aff3 type provides the ability to specify affine transforms on Paths and Shapes.
*/
package graphics2d
