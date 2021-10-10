package graphics2d

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	g2dimg "github.com/jphsd/graphics2d/image"
	. "github.com/jphsd/graphics2d/util"
	"golang.org/x/image/vector"
)

// RenderColoredShape renders the supplied shape with the fill color
// into the destination image.
func RenderColoredShape(dst draw.Image, shape *Shape, fill color.Color) {
	filler := image.NewUniform(fill)
	_ = RenderShapeExt(dst, shape, []float32{0, 0}, filler, image.Point{}, nil, draw.Over)
}

// RenderShape renders the supplied shape with the fill image into
// the destination image.
func RenderShape(dst draw.Image, shape *Shape, filler image.Image, foffs image.Point) {
	_ = RenderShapeExt(dst, shape, []float32{0, 0}, filler, foffs, nil, draw.Over)
}

// RenderFlatten is the standard curve flattening value used when rendering.
const RenderFlatten = 0.6

// RenderShapeExt renders the supplied shape at an offset with the fill and clip images into
// the destination image using op.
func RenderShapeExt(dst draw.Image, shape *Shape, at []float32, filler image.Image, foffs image.Point, clip *image.Alpha, op draw.Op) error {
	rect := dst.Bounds()
	size := rect.Size()

	if clip != nil && !size.Eq(clip.Bounds().Size()) {
		return fmt.Errorf("clip image must have same size as destination image")
	}

	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process paths
	ox, oy := at[0], at[1]
	for _, path := range shape.paths {
		fp := path.Flatten(RenderFlatten) // tolerance 0.6
		step := ToF32(fp.steps[0][0]...)
		rasterizer.MoveTo(ox+step[0], oy+step[1])
		for i, lp := 1, len(fp.steps); i < lp; i++ {
			step = ToF32(fp.steps[i][0]...)
			rasterizer.LineTo(ox+step[0], oy+step[1])
		}
		rasterizer.ClosePath()
	}

	if clip != nil {
		// Obtain rasterizer mask and intersect it against the clip mask
		alpha := image.NewAlpha(rect)
		rasterizer.Draw(alpha, rect, image.Opaque, image.Point{})
		alpha = g2dimg.AlphaAnd(alpha, clip)
		// alpha now has origin {0, 0}
		draw.DrawMask(dst, rect, filler, foffs, alpha, image.Point{}, op)
	} else {
		rasterizer.Draw(dst, rect, filler, foffs)
	}

	return nil
}

// RenderShapeAlpha adds the shape to the mask at the offset with op.
func RenderShapeAlpha(dst *image.Alpha, shape *Shape, at []float32, op draw.Op) {
	rect := dst.Bounds()
	size := rect.Size()
	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process paths
	ox, oy := at[0], at[1]
	for _, path := range shape.paths {
		fp := path.Flatten(RenderFlatten) // tolerance 0.6
		step := ToF32(fp.steps[0][0]...)
		rasterizer.MoveTo(ox+step[0], oy+step[1])
		for i, lp := 1, len(fp.steps); i < lp; i++ {
			step = ToF32(fp.steps[i][0]...)
			rasterizer.LineTo(ox+step[0], oy+step[1])
		}
		rasterizer.ClosePath()
	}

	rasterizer.Draw(dst, rect, image.Opaque, image.Point{})
}

// The following functions use a pre-rendered mask of the shape to draw to the destination image.
// Note that the offsets are Points and not []float32 as previous.

// DrawColoredMask utilizes the supplied shape's mask to draw into the destination image at an offset
// with the fill color.
func DrawColoredMask(dst draw.Image, shape *Shape, offs image.Point, fill color.Color) {
	filler := image.NewUniform(fill)
	_ = DrawMaskExt(dst, shape, offs, filler, image.Point{}, nil, draw.Over)
}

// DrawMask utilizes the supplied shape's mask to draw into the destination image at an offset with
// the filler image.
func DrawMask(dst draw.Image, shape *Shape, offs image.Point, filler image.Image, foffs image.Point) {
	_ = DrawMaskExt(dst, shape, offs, filler, foffs, nil, draw.Over)
}

// DrawMaskExt utilizes the supplied shape's mask to draw into the destination image at an offset with
// the filler, also offset, and clip images using op.
func DrawMaskExt(dst draw.Image, shape *Shape, offs image.Point, filler image.Image, foffs image.Point, clip *image.Alpha, op draw.Op) error {
	rect := dst.Bounds()
	size := rect.Size()

	if clip != nil && !size.Eq(clip.Bounds().Size()) {
		return fmt.Errorf("clip image must have same size as destination image")
	}

	mask := shape.Mask()
	srect := mask.Bounds()
	drect := image.Rectangle{offs, image.Point{offs.X + srect.Dx(), offs.Y + srect.Dy()}}

	if clip != nil {
		// Obtain clip mask subimage and intersect it against the shape mask
		sub := clip.SubImage(drect).(*image.Alpha)
		mask = g2dimg.AlphaAnd(mask, sub)
		// mask now has origin {0, 0}
		draw.DrawMask(dst, drect, filler, foffs, mask, image.Point{}, op)
	} else {
		draw.DrawMask(dst, drect, filler, foffs, mask, srect.Min, op)
	}

	return nil
}
