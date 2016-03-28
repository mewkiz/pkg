// Package imgutil implements some image utility functions.
package imgutil

import (
	"image"
	"image/color"
	_ "image/gif" // support for decoding gif images.
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/bmp"
)

// ReadFile reads an image file (gif, jpeg or png) specified by imgPath and
// returns it as an image.Image.
func ReadFile(imgPath string) (img image.Image, err error) {
	fr, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	img, _, err = image.Decode(fr)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// WriteFile writes the image data to a PNG file specified by imgPath. WriteFile
// creates the named file using mode 0666 (before umask), truncating it if it
// already exists.
func WriteFile(imgPath string, img image.Image) (err error) {
	fw, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer fw.Close()
	if err := png.Encode(fw, img); err != nil {
		return err
	}
	return nil
}

// WriteJPEG writes the image data to a JPEG file specified by imgPath.
// writeJPEG creates the named file using mode 0666 (before umask), truncating
// it if it already exists. The quality of the output image is within the range
// [1, 100]; higher is better.
func WriteJPEG(imgPath string, img image.Image, quality int) (err error) {
	fw, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer fw.Close()
	options := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(fw, img, options); err != nil {
		return err
	}
	return nil
}

// WriteBMP writes the image data to a BMP file specified by imgPath.
// WriteBMP creates the named file using mode 0666 (before umask), truncating
// it if it already exists
func WriteBMP(imgPath string, img image.Image) (err error) {
	fw, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer fw.Close()
	if err := bmp.Encode(fw, img); err != nil {
		return err
	}
	return nil
}

// Equal returns true if the images img1 and img2 are equal, and false
// otherwise.
func Equal(img1, img2 image.Image) bool {
	rect1 := img1.Bounds()
	rect2 := img2.Bounds()

	// Compare bounds.
	if rect1 != rect2 {
		return false
	}

	// Compare pixel colors.
	for x := rect1.Min.X; x < rect1.Max.X; x++ {
		for y := rect1.Min.Y; y < rect1.Max.Y; y++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)
			if !ColorEq(c1, c2) {
				return false
			}
		}
	}

	return true
}

// ColorEq returns true if the colors c1 and c2 are equal, and false otherwise.
func ColorEq(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
