package imgutil

import (
	"image"
	"image/color"
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

// subImage represents a subimage of the given source image.
type subImage struct {
	// Underlying image.
	image.Image
	// Bounds of subimage.
	bounds image.Rectangle
}

// NewSubImage returns a subimage of the given source image based on the
// specified bounds.
func NewSubImage(src image.Image, bounds image.Rectangle) image.Image {
	return &subImage{
		Image:  src,
		bounds: bounds,
	}
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (sub *subImage) Bounds() image.Rectangle {
	return sub.bounds
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (sub *subImage) At(x, y int) color.Color {
	if !(image.Pt(x, y).In(sub.bounds)) {
		return color.Transparent
	}
	return sub.Image.At(x, y)
}
