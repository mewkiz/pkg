// Package pngutil implements some png utility functions.
package pngutil

import "image"
import "image/png"
import "os"

// ReadFile reads a PNG file specified by imgPath and returns it as an
// image.Image.
func ReadFile(imgPath string) (img image.Image, err error) {
	fr, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	img, err = png.Decode(fr)
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
	err = png.Encode(fw, img)
	if err != nil {
		return err
	}
	return nil
}
