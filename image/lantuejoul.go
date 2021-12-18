package image

import (
	"image"
	"image/draw"
)

// LJSkeleton calculates the skeletons of the image using Lanuejoul's formula over n iterations with support b.
// Each skeleton is formed from the tophat of the image eroded by the support b dilated by itself n times.
// The last skeleton is the union of the previous skeletons.
// See https://en.wikipedia.org/wiki/Morphological_skeleton
func LJSkeleton(img image.Image, b [][]bool, n int) []*image.Gray {
	r := img.Bounds()

	// Convert to grayscale if necessary
	gray, ok := img.(*image.Gray)
	if !ok {
		gray := image.NewGray(r)
		draw.Draw(gray, r, img, r.Min, draw.Src)
	}

	sw, sh := len(b[0]), len(b)
	dw, dh := sw/2, sh/2
	simg := SupportToGray(Z0)
	res := make([]*image.Gray, n+1)

	// First round (n = 0) is just regular tophat
	res[0] = TopHat(gray, b)
	res[n] = Copy(res[0])
	for i := 1; i < n; i++ {
		// Dilate current support by b
		simg = growImg(simg, dw, dh)
		simg = Dilate(simg, b)
		suppt := GrayToSupport(simg)

		tmp := Erode(gray, suppt)
		res[i] = TopHat(tmp, b)
		res[n] = Or(res[n], res[i], r.Min)
	}
	return res
}

// LJReconstitute turns a set of skeletons back into the opened version of the original image using the support.
func LJReconstitute(skels []*image.Gray, b [][]bool) *image.Gray {
	sw, sh := len(b[0]), len(b)
	dw, dh := sw/2, sh/2
	simg := SupportToGray(Z0)

	// First round (n = 0) is just the skeleton
	res := Copy(skels[0])

	for i := 1; i < len(skels); i++ {
		// Dilate current support by b
		simg = growImg(simg, dw, dh)
		simg = Dilate(simg, b)
		suppt := GrayToSupport(simg)

		r := skels[i].Bounds()
		res = Or(res, Dilate(skels[i], suppt), r.Min)
	}
	return res
}

// Helper to grow the support image. Grows the img by dw and dh on each edge.
func growImg(img *image.Gray, dw, dh int) *image.Gray {
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewGray(image.Rect(0, 0, w+dw+dw, h+dh+dh))
	draw.Draw(res, image.Rect(dw, dh, dw+w, dh+h), img, imgR.Min, draw.Src)
	return res
}
