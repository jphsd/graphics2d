package image

import "image"

// ExtractChannel returns an image with just the selected channel.
// The returned image is scaled by 1/alpha, if not the alpha channel.
func ColorConvert(img *image.NRGBA, xfm *Aff5) *image.NRGBA {
	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewNRGBA(image.Rect(0, 0, w, h))
	sp, dp := img.Pix, res.Pix
	if img.Stride == 4*w {
		// All at once
		j := 0
		for i := 0; i < len(res.Pix); i++ {
			nv := xfm.Transform(sp[j : j+4])
			copy(nv, dp[j:j+4])
			j += 4
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so := img.PixOffset(0+imgR.Min.X, i+imgR.Min.Y)
			do := res.PixOffset(0, i)
			for j := 0; j < w; j++ {
				nv := xfm.Transform(sp[so : so+4])
				copy(nv, dp[do:do+4])
				so += 4
				do += 4
			}
		}
	}
	return res
}
