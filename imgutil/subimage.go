package imgutil

import (
	"image"
	"image/draw"
	"log"
)

// SubImager is an interface that extends the basic image.Image interface with
// the SubImage method.
type SubImager interface {
	image.Image
	// SubImage returns an image representing the portion of the image visible
	// through r. The returned value shares pixels with the original image.
	SubImage(r image.Rectangle) image.Image
}

// SubFallback returns the provided image.Image as a SubImager. It provides a
// fallback for images missing the SubImage method.
func SubFallback(img image.Image) SubImager {
	sub, ok := img.(SubImager)
	if !ok {
		// The SubImage fallback method will have a severe impact on performance.
		log.Println("imgutil.SubFallback: no SubImage method defined; using fallback.")
		return &subFallback{Image: img}
	}
	return sub
}

// subFallback provides a fallback SubImage method for images.
type subFallback struct {
	image.Image
}

// SubImage returns an image representing the portion of the image src visible
// through r. The returned value doesn't shares pixels with the original image.
func (src *subFallback) SubImage(r image.Rectangle) image.Image {
	dstRect := image.Rect(0, 0, r.Dx(), r.Dy())
	dst := image.NewRGBA(dstRect)
	draw.Draw(dst, dstRect, src, r.Min, draw.Over)
	return dst
}
