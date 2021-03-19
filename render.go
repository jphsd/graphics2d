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
func RenderColoredPath(dst draw.Image, path *Path, at []float32, fill color.Color) {
	filler := image.NewUniform(fill)
	RenderPathExt(dst, path, at, filler, nil, draw.Over)
}

// RenderPath renders the specified path (forced closed) at an offset with the fill image
// into the destination image.
func RenderPath(dst draw.Image, path *Path, at []float32, filler image.Image) {
	RenderPathExt(dst, path, at, filler, nil, draw.Over)
}

// RenderPathExt renders the specified path (forced closed) at an offset with the fill and clip images
// into the destination image using op.
func RenderPathExt(dst draw.Image, path *Path, at []float32, filler image.Image, clip *image.Alpha, op draw.Op) error {
	rect := dst.Bounds()
	size := rect.Size()

	if clip != nil {
		crect := clip.Bounds()
		if !crect.Empty() && !crect.Size().Eq(size) {
			return fmt.Errorf("clip image must have same size as destination image")
		}
	}

	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process path
	ox, oy := at[0], at[1]
	fp := path.Flatten(0.6) // tolerance 0.6
	step := ToF32(fp.steps[0][0]...)
	rasterizer.MoveTo(ox+step[0], oy+step[1])
	for i, lp := 1, len(fp.steps); i < lp; i++ {
		step = ToF32(fp.steps[i][0]...)
		rasterizer.LineTo(ox+step[0], oy+step[1])
	}
	rasterizer.ClosePath()

	if clip != nil {
		// Obtain rasterizer mask and intersect it against the clip mask
		alpha := image.NewAlpha(image.Rect(0, 0, size.X, size.Y))
		rasterizer.Draw(alpha, rect, image.Opaque, image.Point{})
		alpha = g2dimg.AlphaAnd(alpha, clip)
		draw.DrawMask(dst, rect, filler, image.Point{}, alpha, image.Point{}, op)
	} else {
		rasterizer.Draw(dst, rect, filler, image.Point{})
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
	fp := path.Flatten(0.6) // tolerance 0.6
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
func RenderColoredShape(dst draw.Image, shape *Shape, at []float32, fill color.Color) {
	filler := image.NewUniform(fill)
	RenderShapeExt(dst, shape, at, filler, nil, draw.Over)
}

// RenderShape renders the supplied shape at an offset with the fill image into
// the destination image.
func RenderShape(dst draw.Image, shape *Shape, at []float32, filler image.Image) {
	RenderShapeExt(dst, shape, at, filler, nil, draw.Over)
}

// RenderShapeExt renders the supplied shape at an offset with the fill and clip images into
// the destination image using op.
func RenderShapeExt(dst draw.Image, shape *Shape, at []float32, filler image.Image, clip *image.Alpha, op draw.Op) error {
	rect := dst.Bounds()
	size := rect.Size()

	if clip != nil {
		crect := clip.Bounds()
		if !crect.Empty() && !crect.Size().Eq(size) {
			return fmt.Errorf("clip image must have same size as destination image")
		}
	}

	rasterizer := vector.NewRasterizer(size.X, size.Y)
	rasterizer.DrawOp = op

	// Process paths
	ox, oy := at[0], at[1]
	for _, path := range shape.paths {
		fp := path.Flatten(0.6) // tolerance 0.6
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
		alpha := image.NewAlpha(image.Rect(0, 0, size.X, size.Y))
		rasterizer.Draw(alpha, rect, image.Opaque, image.Point{})
		alpha = g2dimg.AlphaAnd(alpha, clip)
		draw.DrawMask(dst, rect, filler, image.Point{}, alpha, image.Point{}, op)
	} else {
		rasterizer.Draw(dst, rect, filler, image.Point{})
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
		fp := path.Flatten(0.6) // tolerance 0.6
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
