package graphics2d

import (
	//"fmt"
	"image"
	"image/color"
	"image/draw"

	g2dimg "github.com/jphsd/graphics2d/image"
	"github.com/jphsd/graphics2d/util"
	"golang.org/x/image/vector"
)

// RenderColoredShape renders the supplied shape with the fill color
// into the destination image.
func RenderColoredShape(dst draw.Image, shape *Shape, fill color.Color) {
	filler := image.NewUniform(fill)
	RenderShapeExt(dst, shape, filler, image.Point{}, nil, image.Point{}, draw.Over)
}

// RenderShape renders the supplied shape with the offset fill image into
// the destination image.
func RenderShape(dst draw.Image, shape *Shape, filler image.Image, fx, fy int) {
	RenderShapeExt(dst, shape, filler, image.Point{fx, fy}, nil, image.Point{}, draw.Over)
}

// DefaultRenderFlatten is the standard curve flattening value.
const DefaultRenderFlatten = 0.6

// RenderFlatten is the curve flattening value used when rendering.
var RenderFlatten = DefaultRenderFlatten

// RenderShapeExt renders the supplied shape with the fill and clip images into
// the destination image using op.
func RenderShapeExt(dst draw.Image, shape *Shape, filler image.Image, foffs image.Point, clip *image.Alpha, coffs image.Point, op draw.Op) {
	rect := dst.Bounds()

	// To avoid unnecessary work, reduce the rasterizer size to the shape width and height
	// clipped by the destination image bounds
	srect := shape.Bounds()
	srect = rect.Intersect(srect)
	size := srect.Size()
	rasterizer := vector.NewRasterizer(size.X, size.Y) // Rasterizer has implicit r.Min of {0, 0}
	rasterizer.DrawOp = op

	// Process paths translated by -srect.Min and add srect.Min fo filler offest
	foffs = image.Point{foffs.X + srect.Min.X, foffs.Y + srect.Min.Y}
	minx, miny := float32(srect.Min.X), float32(srect.Min.Y)

	for _, path := range shape.paths {
		fp := path.Flatten(RenderFlatten) // tolerance 0.6
		step := util.ToF32(fp.steps[0][0]...)
		rasterizer.MoveTo(step[0]-minx, step[1]-miny)
		for i, lp := 1, len(fp.steps); i < lp; i++ {
			step = util.ToF32(fp.steps[i][0]...)
			rasterizer.LineTo(step[0]-minx, step[1]-miny)
		}
		rasterizer.ClosePath()
	}

	if clip != nil {
		// Obtain rasterizer mask and intersect it against the clip mask
		mask := image.NewAlpha(srect)
		rasterizer.Draw(mask, srect, image.Opaque, image.Point{})
		mask = g2dimg.AlphaAnd(mask, clip, coffs)
		draw.DrawMask(dst, srect, filler, foffs, mask, srect.Min, op)
		return
	}

	rasterizer.Draw(dst, srect, filler, foffs)
}

// RenderShapeAlpha creates and returns the shape's alpha mask. The mask size is determined by the shape's
// bounds and the mask is located at {0, 0}.
func RenderShapeAlpha(shape *Shape) *image.Alpha {
	srect := shape.Bounds()
	size := srect.Size()
	rect := image.Rectangle{image.Point{}, image.Point{size.X, size.Y}}

	rasterizer := vector.NewRasterizer(size.X, size.Y) // Rasterizer has implicit r.Min of {0, 0}
	rasterizer.DrawOp = draw.Src

	// Process paths translated by -srect.Min
	minx, miny := float32(srect.Min.X), float32(srect.Min.Y)
	for _, path := range shape.paths {
		fp := path.Flatten(RenderFlatten) // tolerance 0.6
		step := util.ToF32(fp.steps[0][0]...)
		rasterizer.MoveTo(step[0]-minx, step[1]-miny)
		for i, lp := 1, len(fp.steps); i < lp; i++ {
			step = util.ToF32(fp.steps[i][0]...)
			rasterizer.LineTo(step[0]-minx, step[1]-miny)
		}
		rasterizer.ClosePath()
	}

	mask := image.NewAlpha(rect)
	rasterizer.Draw(mask, rect, image.Opaque, image.Point{})
	return mask
}
