package arn

import (
	"image"
	"os"
	"path"

	"github.com/chai2010/webp"
)

// OriginalImageExtensions includes all the formats that an avatar source could have sent to us.
var OriginalImageExtensions = []string{
	".jpg",
	".png",
	".gif",
}

const (
	// AvatarSmallSize is the minimum size in pixels of an avatar.
	AvatarSmallSize = 100

	// AvatarMaxSize is the maximum size in pixels of an avatar.
	AvatarMaxSize = 560
)

// LoadImage loads an image from the given path.
func LoadImage(path string) (img image.Image, format string, err error) {
	f, openErr := os.Open(path)

	if openErr != nil {
		return nil, "", openErr
	}

	img, format, decodeErr := image.Decode(f)

	if decodeErr != nil {
		return nil, "", decodeErr
	}

	return img, format, nil
}

// SaveWebP saves an image as a file in WebP format.
func SaveWebP(img image.Image, out string, quality float32) error {
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

// FindFileWithExtension tries to test different file extensions.
func FindFileWithExtension(baseName string, dir string, extensions []string) string {
	for _, ext := range extensions {
		if _, err := os.Stat(path.Join(dir, baseName+ext)); !os.IsNotExist(err) {
			return dir + baseName + ext
		}
	}

	return ""
}
