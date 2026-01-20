package svg

import (
	"encoding/xml"
	"io"
)

const (
	SVGHeader = "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\">"
	SVGFooter = "</svg>"
)

// NewEncoder returns an xml.Encoder wrapped around w that's already had the SVG header written to it.
func NewEncoder(w io.Writer) *xml.Encoder {
	w.Write([]byte(SVGHeader))
	return xml.NewEncoder(w)
}

// Complete adds the SVG footer to w.
func Complete(w io.Writer) {
	w.Write([]byte(SVGFooter))
}
