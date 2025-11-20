# Yet Another 2D Graphics Package For Go
[![Go Reference](https://pkg.go.dev/badge/github.com/jphsd/graphics2d.svg)](https://pkg.go.dev/github.com/jphsd/graphics2d)
[![Go Report Card](https://goreportcard.com/badge/github.com/jphsd/graphics2d)](https://goreportcard.com/report/github.com/jphsd/graphics2d)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/jphsd/graphics2d)

[![Splash image created with graphics2d](./doc/splash.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Splash)

## 1. Introduction

[Graphics2D](https://pkg.go.dev/github.com/jphsd/graphics2d) is a vector based drawing package that
leverages [golang.org/x/image/vector](https://pkg.go.dev/golang.org/x/image/vector) to render shapes into an image.

The vector package extends [image/draw](https://pkg.go.dev/image/draw) to create a mask that a source image
is rendered through into the destination image.

Graphics2D follows this convention.

### Paths

The primary type in the package is the [Path](https://pkg.go.dev/github.com/jphsd/graphics2d#Path).
A path represents a single movement of a pen, from pen down to pen up. Paths are composed of steps with
some number of points in them.
The number of points determines the order of the Bezier curve generated.
The path methods [LineTo](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.LineTo)
and [CurveTo](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.CurveTo)
are just synonyms for [AddStep](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.AddStep)
Once created, a path can be left as is (open), or closed [Close](https://pkg.go.dev/github.com/jphsd/graphics2d#Path).
A closed path can no longer be extended and a line is created from the first point in the path to its last.

### Shapes

Shapes allow multiple paths to be combined to produce more complex drawings.
For example, the figure 8 is composed of three paths; the outline, and the two holes in it.

## 2. Basic Shapes
[![Fig1 image created with graphics2d](./doc/fig1.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig01)
The shapes above were created using [Line](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.Line),
[RegularPolygon](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.RegularPolygon),
[Circle](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.Circle),
and [Oval](https://pkg.go.dev/github.com/jphsd/graphics2d#Path.Oval).
These are just some of the constructors available for the Path type.

## 3. Bezier Curves
[![Fig2 image created with graphics2d](./doc/fig2.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig02)
[Bezier curves](https://en.wikipedia.org/wiki/B%C3%A9zier_curve)
are polynomial curves.
Most vector packages support first, second and third order curves, lines, quadratic and cubic curves respectively.
The path AddStep method has no upper limit on the number of control points that can be specified.
The last example above on the right is a quartic curve.

## 4. Arcs And ArcStyles
[![Fig3 image created with graphics3d](./doc/fig3.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig03)
Various arc path constructors are available and typically take an offset start angle and a sweep angle.
The arcs are approximated from cubic bezier curves. 
Arcs must have a [style](https://pkg.go.dev/github.com/jphsd/graphics2d#ArcStyle)
associated with them which is one of ArcOpen, ArcPie or ArcChord as shown above.

## 5. Reentrant Shapes
[![Fig4 image created with graphics2d](./doc/fig4.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig04)

## 6. Using Path Processors
[![Fig5 image created with graphics2d](./doc/fig5.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig05)
[![Fig6 image created with graphics2d](./doc/fig6.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig06)

## 7. Using Fonts
[![Fig7 image created with graphics2d](./doc/fig7.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig07)

## 8. Dashing With Path Processors
[![Fig8 image created with graphics2d](./doc/fig8.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig08)

## 9. Tracing With Path Processors
### Join
[![Fig9 image created with graphics2d](./doc/fig9.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig09)
[![Fig10 image created with graphics2d](./doc/fig10.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig10)

## 10. Outlining With Stroke Path Processor
### Cap
[![Fig11 image created with graphics2d](./doc/fig11.png)](https://pkg.go.dev/github.com/jphsd/graphics2d#example-package-Fig11)
