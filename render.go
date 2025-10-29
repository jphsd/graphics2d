package graphics2d

import (
	//"fmt"
	"image"
	"image/draw"

	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/util"
	"golang.org/x/image/vector"
)

// RenderColoredShape renders the supplied shape with the fill color into the destination image.
func RenderColoredShape(dst draw.Image, shape *Shape, fill color.Color) {
	RenderShape(dst, shape, image.NewUniform(fill))
}

// RenderShape renders the supplied shape with the fill image into the destination image.
func RenderShape(dst draw.Image, shape *Shape, filler image.Image) {
	r := dst.Bounds()
	RenderShapeExt(dst, r, shape, filler, r.Min, nil, image.Point{}, draw.Over)
}

// RenderClippedShape renders the supplied shape with the fill image into the destination image
// as masked by the clip shape.
func RenderClippedShape(dst draw.Image, shape, clip *Shape, filler image.Image) {
	r := dst.Bounds()
	RenderShapeExt(dst, r, shape, filler, r.Min, clip.Mask(), r.Min, draw.Over)
}

// DefaultRenderFlatten is the standard curve flattening value.
const DefaultRenderFlatten = 0.6

// RenderFlatten is the curve flattening value used when rendering.
var RenderFlatten = DefaultRenderFlatten

// RenderShapeExt renders the supplied shape with the fill and clip images into
// the destination image region using op.
func RenderShapeExt(dst draw.Image, drect image.Rectangle, shape *Shape, filler image.Image, fp image.Point, mask *image.Alpha, mp image.Point, op draw.Op) {
	orig := drect.Min

	// To avoid unnecessary work, reduce the rasterizer size to the shape width and height
	// clipped by the destination image bounds, the filler image and the clip image
	srect := shape.Bounds()
	drect = drect.Intersect(srect)
	// the filler bounds
	drect = drect.Intersect(filler.Bounds().Add(orig.Sub(fp)))
	// and the clip bounds (if present)
	if mask != nil {
		drect = drect.Intersect(mask.Bounds().Add(orig.Sub(mp)))
	}
	if drect.Empty() {
		return
	}

	size := drect.Size()
	dx, dy := drect.Min.X-orig.X, drect.Min.Y-orig.Y

	// Make rasterizer, note rasterizer has implicit r.Min of {0, 0}
	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process paths translated by -drect.Min since mp is {0, 0} in the vectorizer
	minx, miny := float32(drect.Min.X), float32(drect.Min.Y)

	for _, path := range shape.paths {
		prect := path.Bounds() // shape.Bounds() will have caused these to be generated already
		prect = drect.Intersect(prect)
		if prect.Empty() {
			continue
		}
		fp := path.Flatten(RenderFlatten) // default tolerance 0.6
		step := util.ToF32(fp.steps[0][0]...)
		rasterizer.MoveTo(step[0]-minx, step[1]-miny)
		for i, lp := 1, len(fp.steps); i < lp; i++ {
			step = util.ToF32(fp.steps[i][0]...)
			rasterizer.LineTo(step[0]-minx, step[1]-miny)
		}
		rasterizer.ClosePath()
	}

	fp.X += dx
	fp.Y += dy

	if mask == nil {
		rasterizer.Draw(dst, drect, filler, fp)
		return
	}

	// Process clip mask - obtain rasterizer mask and intersect it against the clip mask
	nmask := image.NewAlpha(drect)
	mp.X += dx
	mp.Y += dy
	rasterizer.Draw(nmask, drect, mask, mp)
	draw.DrawMask(dst, drect, filler, fp, nmask, drect.Min, op)
}
