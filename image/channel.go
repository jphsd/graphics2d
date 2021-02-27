package image

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
)

// Lookups for alpha scaling
var invPremult = [][]uint8{}
var premult = [][]uint8{}
var pmlock = sync.Mutex{}

// ExtractChannel returns an image with just the selected channel.
// The returned image is scaled by 1/alpha, if not the alpha channel.
func ExtractChannel(img *image.RGBA, ch int) *image.Gray {
	if ch < 0 || ch > 3 {
		panic(fmt.Errorf("Requested channel not in range"))
	}

	pmlock.Lock()
	if len(invPremult) == 0 {
		invPremult = initInvPremult()
	}
	pmlock.Unlock()

	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewGray(image.Rect(0, 0, w, h))
	sp, dp := img.Pix, res.Pix
	if img.Stride == 4*w {
		// All at once
		if ch != 3 {
			j := 0
			for i := 0; i < len(res.Pix); i++ {
				dp[i] = invPremult[sp[j+ch]][sp[j+3]]
				j += 4
			}
		} else { // Alpha
			j := 0
			for i := 0; i < len(res.Pix); i++ {
				dp[i] = sp[j+3]
				j += 4
			}
		}
	} else {
		// Scan line at a time
		if ch != 3 {
			for i := 0; i < h; i++ {
				so, do := img.PixOffset(0, i), res.PixOffset(0, i)
				for j := 0; j < w; j++ {
					dp[do] = invPremult[sp[so+ch]][sp[so+3]]
					so += 4
					do++
				}
			}
		} else { // Alpha
			for i := 0; i < h; i++ {
				so, do := img.PixOffset(0, i)+3, res.PixOffset(0, i)
				for j := 0; j < w; j++ {
					dp[do] = sp[so]
					so += 4
					do++
				}
			}
		}
	}
	return res
}

// ReplaceChannel in an image with the one supplied. The supplied image is scaled by alpha,
// if one of R, G or B. If the alpha channel is replaced then RGB are first scaled by 1/alpha and
// then by alpha'
func ReplaceChannel(img *image.RGBA, ch int, rep *image.Gray) *image.RGBA {
	if ch < 0 || ch > 3 {
		panic(fmt.Errorf("Requested channel not in range"))
	}

	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	if rep.Rect.Dx() != w || rep.Rect.Dy() != h {
		panic(fmt.Errorf("images are different sizes"))
	}

	pmlock.Lock()
	if len(invPremult) == 0 {
		invPremult = initInvPremult()
	}
	if len(premult) == 0 {
		premult = initPremult()
	}
	pmlock.Unlock()

	res := image.NewRGBA(image.Rect(0, 0, w, h))
	sp, dp, rp := img.Pix, res.Pix, rep.Pix
	if img.Stride == 4*w && rep.Stride == w {
		// All at once
		for i := 0; i < len(res.Pix); i++ {
			j := i * 4
			a := sp[j+3]
			switch ch {
			case 0:
				dp[j] = premult[rp[i]][a]
				dp[j+1] = sp[j+1]
				dp[j+2] = sp[j+2]
				dp[j+3] = a
			case 1:
				dp[j] = sp[j]
				dp[j+1] = premult[rp[i]][a]
				dp[j+2] = sp[j+2]
				dp[j+3] = a
			case 2:
				dp[j] = sp[j]
				dp[j+1] = sp[j+1]
				dp[j+2] = premult[rp[i]][a]
				dp[j+3] = a
			case 3:
				r := invPremult[sp[j]][a]
				g := invPremult[sp[j+1]][a]
				b := invPremult[sp[j+2]][a]
				a = rp[i]
				dp[j] = premult[r][a]
				dp[j+1] = premult[g][a]
				dp[j+2] = premult[b][a]
				dp[j+3] = a
			}
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so, do, ro := img.PixOffset(0, i), res.PixOffset(0, i), rep.PixOffset(0, i)
			for k := 0; k < w; k++ {
				j := k * 4
				sj, dj := so+j, do+j
				a := sp[sj+3]
				switch ch {
				case 0:
					dp[dj] = premult[rp[ro+k]][a]
					dp[dj+1] = sp[sj+1]
					dp[dj+2] = sp[sj+2]
					dp[dj+3] = a
				case 1:
					dp[dj] = sp[sj]
					dp[dj+1] = premult[rp[ro+k]][a]
					dp[dj+2] = sp[sj+2]
					dp[dj+3] = a
				case 2:
					dp[dj] = sp[sj]
					dp[dj+1] = sp[sj+1]
					dp[dj+2] = premult[rp[ro+k]][a]
					dp[dj+3] = a
				case 3:
					r := invPremult[sp[sj]][a]
					g := invPremult[sp[sj+1]][a]
					b := invPremult[sp[sj+2]][a]
					a = rp[ro+k]
					dp[dj] = premult[r][a]
					dp[dj+1] = premult[g][a]
					dp[dj+2] = premult[b][a]
					dp[dj+3] = a
				}
			}
		}
	}
	return res
}

// SwitchChannels in an image (not alpha)
func SwitchChannels(img *image.RGBA, ch1, ch2 int) *image.RGBA {
	if ch1 < 0 || ch1 > 2 || ch2 < 0 || ch2 > 2 {
		panic(fmt.Errorf("channel out of range"))
	}

	imgR := img.Bounds()
	w, h := imgR.Dx(), imgR.Dy()
	res := image.NewRGBA(image.Rect(0, 0, w, h))

	if ch1 == ch2 {
		// Copy the image as is
		draw.Draw(res, res.Bounds(), img, image.Point{}, draw.Src)
		return res
	}

	// Find the odd one out which isn't changing
	ooo := 0
	if ch1 == 0 || ch2 == 0 {
		ooo = 1
	}
	if ooo == 1 && (ch1 == 1 || ch2 == 1) {
		ooo = 2
	}

	sp, dp := img.Pix, res.Pix
	if img.Stride == 4*w {
		// All at once
		for i := 0; i < len(res.Pix); i += 4 {
			switch ooo {
			case 0:
				dp[i], dp[i+1], dp[i+2], dp[i+3] = sp[i], sp[i+2], sp[i+1], sp[i+3]
			case 1:
				dp[i], dp[i+1], dp[i+2], dp[i+3] = sp[i+2], sp[i+1], sp[i], sp[i+3]
			case 2:
				dp[i], dp[i+1], dp[i+2], dp[i+3] = sp[i+1], sp[i], sp[i+2], sp[i+3]
			}
		}
	} else {
		// Scan line at a time
		for i := 0; i < h; i++ {
			so, do := img.PixOffset(0, i), res.PixOffset(0, i)
			for j := 0; j < w; j++ {
				switch ooo {
				case 0:
					dp[do], dp[do+1], dp[do+2], dp[do+3] = sp[so], sp[so+2], sp[so+1], sp[so+3]
				case 1:
					dp[do], dp[do+1], dp[do+2], dp[do+3] = sp[so+2], sp[so+1], sp[so], sp[so+3]
				case 2:
					dp[do], dp[do+1], dp[do+2], dp[do+3] = sp[so+1], sp[so], sp[so+2], sp[so+3]
				}
				so += 4
				do += 4
			}
		}
	}
	return res
}

// CombineChannels combines mono-channel images into a single image. The R, G, B channels
// will be scaled by alpha if scale is true.
func CombineChannels(chR, chG, chB, chA *image.Gray, scale bool) *image.RGBA {
	chRR := chR.Bounds()
	w, h := chRR.Dx(), chRR.Dy()
	chGR := chG.Bounds()
	chBR := chB.Bounds()
	chAR := chA.Bounds()

	if w != chGR.Dx() || h != chGR.Dy() || w != chBR.Dx() || h != chBR.Dy() ||
		w != chAR.Dx() || h != chAR.Dy() {
		panic(fmt.Errorf("images are different sizes"))
	}

	res := image.NewRGBA(image.Rect(0, 0, w, h))
	if scale {
		pmlock.Lock()
		if len(premult) == 0 {
			premult = initPremult()
		}
		pmlock.Unlock()
		if chR.Stride == w && chG.Stride == w && chB.Stride == w && chA.Stride == w {
			// All at once
			j := 0
			for i := 0; i < len(chR.Pix); i++ {
				a := chA.Pix[i]
				res.Pix[j] = premult[chR.Pix[i]][a]
				res.Pix[j+1] = premult[chG.Pix[i]][a]
				res.Pix[j+2] = premult[chB.Pix[i]][a]
				res.Pix[j+3] = a
				j += 4
			}
		} else {
			// Scan line at a time
			for i := 0; i < h; i++ {
				ro, ggo, bo, ao, do := chR.PixOffset(0, i), chG.PixOffset(0, i), chB.PixOffset(0, i), chA.PixOffset(0, i), res.PixOffset(0, i)
				for j := 0; j < w; j++ {
					k := do + j*4
					a := chA.Pix[ao+j]
					res.Pix[k] = premult[chR.Pix[ro+j]][a]
					res.Pix[k+1] = premult[chG.Pix[ggo+j]][a]
					res.Pix[k+2] = premult[chB.Pix[bo+j]][a]
					res.Pix[k+3] = a
				}
			}
		}
	} else {
		if chR.Stride == w && chG.Stride == w && chB.Stride == w && chA.Stride == w {
			// All at once
			j := 0
			for i := 0; i < len(chR.Pix); i++ {
				res.Pix[j] = chR.Pix[i]
				res.Pix[j+1] = chG.Pix[i]
				res.Pix[j+2] = chB.Pix[i]
				res.Pix[j+3] = chA.Pix[i]
				j += 4
			}
		} else {
			// Scan line at a time
			for i := 0; i < h; i++ {
				ro, ggo, bo, ao, do := chR.PixOffset(0, i), chG.PixOffset(0, i), chB.PixOffset(0, i), chA.PixOffset(0, i), res.PixOffset(0, i)
				for j := 0; j < w; j++ {
					k := do + j*4
					res.Pix[k] = chR.Pix[ro+j]
					res.Pix[k+1] = chG.Pix[ggo+j]
					res.Pix[k+2] = chB.Pix[bo+j]
					res.Pix[k+3] = chA.Pix[ao+j]
				}
			}
		}
	}

	return res
}

// Note - this is lossy
func initInvPremult() [][]uint8 {
	res := make([][]uint8, 256)
	for i := 0; i < 256; i++ {
		nres := make([]uint8, 256)
		v := uint32(i)
		v |= v << 8
		nres[0] = uint8(0)
		for a := 1; a < 256; a++ {
			nv := (v * 0xffff) / uint32(a*0x101)
			nres[a] = uint8(nv >> 8)
		}
		res[i] = nres
	}
	return res
}

func initPremult() [][]uint8 {
	res := make([][]uint8, 256)
	for i := 0; i < 256; i++ {
		nres := make([]uint8, 256)
		v := uint32(i)
		v |= v << 8
		for a := 0; a < 256; a++ {
			nv := uint32(a) * v
			nv /= 0xff
			nres[a] = uint8(nv >> 8)
		}
		res[i] = nres
	}
	return res
}
