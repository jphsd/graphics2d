# SVG Utilities

## 1. Introduction

This package provides ways to render paths and shapes as [SVG](https://www.w3.org/TR/SVG11/Overview.html)
documents.

It makes use of the [encodings/xml](https://pkg.go.dev/encodings/xml)
package within the standard library.

This package is not intended as a full SVG implementation.
It simply provides a way to output the paths and shapes created
in [graphics2d](https://pkg.go.dev/github.com/jphsd/graphics2d)
as SVG.
The separate [xml/svg](https://pkg.go.dev/github.com/jphsd/xml/svg)
package provides ways of converting an SVG document into
a [Renderable](https://pkg.go.dev/github.com/jphsd/graphics2d#Renderable).

## 2. MarshallXML

[Path](https://pkg.go.dev/github.com/jphsd/graphics2d#Path)
and [Shape](https://pkg.go.dev/github.com/jphsd/graphics2d#Shape)
both implement the [xml.Marshaler](https://pkg.go.dev/xml#Marshaler)
interface allowing calls to [xml.Marshal](https://pkg.go.dev/xml#Marshal)
to be made with either type as an argument.
The result will be a snippet of XML.

### Path

This produces an SVG [path element](https://www.w3.org/TR/SVG11/paths.html#PathElement)
with the [d attribute](https://www.w3.org/TR/SVG11/paths.html#DAttribute)
containing the path steps.
Since SVG doesn't handle high order steps, the path is flattened prior to marshaling.

### Shape

SVG doesn't have an explicit shape element so the [group element](https://www.w3.org/TR/SVG11/struct.html#GElement)
is used instead.
All the paths in a shape will be marshaled under a single group.

## 3. Rendering

The rendering functions follow the same pattern as the image rendering ones.
Where an image was previously provided, now an XML encoder is used.

The two rendering functions supported are:
- [DrawShape](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#DrawShape)
- [RenderColoredShape](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#RenderColoredShape)

### Caveats

- Clipped shapes are not supported in the SVG encoding.
  An error will be returned if clipping has been used.
- Only monotone fillers are supported.
  The color used to fill a shape is taken from the color value at (0, 0) in the filler image.

## 4. Renderable

The [Renderable](https://pkg.go.dev/github.com/jphsd/graphics2d#Renderable)
type is rendered in SVG using [RenderRenderable](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#RenderRenderable).
The same caveats mentioned in section 3 apply.

## 5. SVG Wrapper

A convenience function, [NewEncoder](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#NewEncoder),
is supplied that wraps an [io.Writer](https://pkg.go.dev/io#Writer)
in an [xml.Encoder](https://pkg.go.dev/xml#Encoder)
with the standard [SVG element](https://www.w3.org/TR/SVG11/struct.html#SVGElement)
already applied.

Another function, [Complete](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#Complete),
writes the trailing SVG token at the end.

## 6. Image Conversion

[Image](https://pkg.go.dev/github.com/jphsd/graphics2d/svg#Image)
provides a way to encode an Image as an SVG [image element](https://www.w3.org/TR/SVG11/struct.html#ImageElement).
The image is first converted to PNG format and then stored base64 encoded in the image
[xlink:href attribute](https://www.w3.org/TR/SVG11/struct.html#ImageElementHrefAttribute).
