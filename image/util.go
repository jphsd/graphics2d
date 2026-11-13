package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	//_ "github.com/jdeng/goheif"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	_ "image/png"
)

// Fill fills the image with the color
func Fill(img draw.Image, col color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, Point{}, draw.Src)
}

// SaveImage is a utility function to save an image as a .png.
func SaveImage(img Image, name string) error {
	fDst, err := os.Create(fmt.Sprintf("%s.png", name))
	if err != nil {
		return err
	}
	defer fDst.Close()
	return png.Encode(fDst, img)
}

// ReadImage is a utility function to read an image from a file.
// The following formats are supported - bmp, heic, jpeg, png, tiff, webp.
func ReadImage(name string) (Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
