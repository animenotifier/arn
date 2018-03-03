package imageoutput

import (
	"image"
	"os"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

// WebPFile ...
type WebPFile struct {
	Directory string
	Size      int
	Quality   float32
}

// Save writes the image in WebP format to the file system.
func (output *WebPFile) Save(avatar *MetaImage, baseName string) error {
	img := avatar.Image

	// Resize if needed
	if img.Bounds().Dx() > output.Size {
		img = resize.Resize(uint(output.Size), 0, img, resize.Lanczos3)
	}

	// Write to file
	fileName := output.Directory + baseName + ".webp"
	return saveWebP(img, fileName, output.Quality)
}

// saveWebP saves an image as a file in WebP format.
func saveWebP(img image.Image, out string, quality float32) error {
	file, writeErr := os.Create(out)

	if writeErr != nil {
		return writeErr
	}

	defer file.Close()

	encodeErr := webp.Encode(file, img, &webp.Options{
		Quality: quality,
	})

	return encodeErr
}
