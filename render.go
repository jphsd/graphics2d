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

// RenderColoredPath renders the specified path (forced closed) at an offset with the fill color
// into the destination image.
func RenderColoredPath(dst draw.Image, path *Path, fill color.Color) {
	filler := image.NewUniform(fill)
	_ = RenderPathExt(dst, path, []float32{0, 0}, filler, image.Point{}, nil, draw.Over)
}

// RenderPath renders the specified path (forced closed) at an offset with the fill image
// into the destination image.
func RenderPath(dst draw.Image, path *Path, filler image.Image, foffs image.Point) {
	_ = RenderPathExt(dst, path, []float32{0, 0}, filler, foffs, nil, draw.Over)
}

// RenderFlatten is the standard curve flattening value used when rendering.
const RenderFlatten = 0.6

// RenderPathExt renders the specified path (forced closed) at an offset with the fill and clip images
// into the destination image using op.
func RenderPathExt(dst draw.Image, path *Path, at []float32, filler image.Image, foffs image.Point, clip *image.Alpha, op draw.Op) error {
	rect := dst.Bounds()
	size := rect.Size()

	if clip != nil && !size.Eq(clip.Bounds().Size()) {
		return fmt.Errorf("clip image must have same size as destination image")
	}

	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process path
	ox, oy := at[0], at[1]
	fp := path.Flatten(RenderFlatten) // tolerance 0.6
	step := ToF32(fp.steps[0][0]...)
	rasterizer.MoveTo(ox+step[0], oy+step[1])
	for i, lp := 1, len(fp.steps); i < lp; i++ {
		step = ToF32(fp.steps[i][0]...)
		rasterizer.LineTo(ox+step[0], oy+step[1])
	}
	rasterizer.ClosePath()

	if clip != nil {
		// Obtain rasterizer mask and intersect it against the clip mask
		alpha := image.NewAlpha(rect)
		rasterizer.Draw(alpha, rect, image.Opaque, image.Point{})
		alpha = g2dimg.AlphaAnd(alpha, clip)
		// alpha now has {0, 0} origin
		draw.DrawMask(dst, rect, filler, foffs, alpha, image.Point{}, op)
	} else {
		rasterizer.Draw(dst, rect, filler, foffs)
	}

	return nil
}

// RenderPathAlpha adds the path (forced closed) to the mask at the supplied offset with op.
func RenderPathAlpha(dst *image.Alpha, path *Path, at []float32, op draw.Op) {
	rect := dst.Bounds()
	size := rect.Size()
	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process path
	ox, oy := at[0], at[1]
	fp := path.Flatten(RenderFlatten) // tolerance 0.6
	step := ToF32(fp.steps[0][0]...)
	rasterizer.MoveTo(ox+step[0], oy+step[1])
	for i, lp := 1, len(fp.steps); i < lp; i++ {
		step = ToF32(fp.steps[i][0]...)
		rasterizer.LineTo(ox+step[0], oy+step[1])
	}
	rasterizer.ClosePath()

	rasterizer.Draw(dst, rect, image.Opaque, image.Point{})
}

// RenderColoredShape renders the supplied shape at an offset with the fill color
// into the destination image.
func RenderColoredShape(dst draw.Image, shape *Shape, fill color.Color) {
	filler := image.NewUniform(fill)
	_ = RenderShapeExt(dst, shape, []float32{0, 0}, filler, image.Point{}, nil, draw.Over)
}

// RenderShape renders the supplied shape at an offset with the fill image into
// the destination image.
func RenderShape(dst draw.Image, shape *Shape, filler image.Image, foffs image.Point) {
	_ = RenderShapeExt(dst, shape, []float32{0, 0}, filler, foffs, nil, draw.Over)
}

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

// DrawColoredShape utilizes the supplied shape's mask to draw into the destination image at an offset
// with the fill color.
func DrawColoredShape(dst draw.Image, shape *Shape, offs image.Point, fill color.Color) {
	filler := image.NewUniform(fill)
	_ = DrawShapeExt(dst, shape, offs, filler, image.Point{}, nil, draw.Over)
}

// DrawShape utilizes the supplied shape's mask to draw into the destination image at an offset with
// the filler image.
func DrawShape(dst draw.Image, shape *Shape, offs image.Point, filler image.Image, foffs image.Point) {
	_ = DrawShapeExt(dst, shape, offs, filler, foffs, nil, draw.Over)
}

// DrawShapeExt utilizes the supplied shape's mask to draw into the destination image at an offset with
// the filler, also offset, and clip images using op.
func DrawShapeExt(dst draw.Image, shape *Shape, offs image.Point, filler image.Image, foffs image.Point, clip *image.Alpha, op draw.Op) error {
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
