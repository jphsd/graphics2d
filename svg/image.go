package svg

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"image"
	"image/png"
)

type ximage struct {
	Id     string `xml:"id,attr,omitempty"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Data   string `xml:"xlink:href,attr"`
}

// Image writes img to the encoder as a base64 encoded png using the <image> element.
func Image(enc *xml.Encoder, img image.Image, id string) error {
	// Encode image as .png bytes
	b := &bytes.Buffer{}
	png.Encode(b, img)

	// Convert to b64
	bb := b.Bytes()
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(bb)))
	base64.StdEncoding.Encode(b64, bb)

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	data := "data:image/png;base64," + string(b64)
	return enc.EncodeElement(ximage{id, width, height, data}, xml.StartElement{Name: xml.Name{"", "image"}})
}
