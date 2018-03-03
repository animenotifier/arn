package arn

import (
	"image"
	"os"
	"path"

	"github.com/animenotifier/arn/imageoutput"
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

	// AvatarWebPQuality is the WebP quality of avatars.
	AvatarWebPQuality = 80
)

// Define the avatar outputs
var avatarOutputs = []imageoutput.Output{
	// Original - Large
	&imageoutput.OriginalFile{
		Directory: "images/avatars/large/",
		Size:      AvatarMaxSize,
	},

	// Original - Small
	&imageoutput.OriginalFile{
		Directory: "images/avatars/small/",
		Size:      AvatarSmallSize,
	},

	// WebP - Large
	&imageoutput.WebPFile{
		Directory: "images/avatars/large/",
		Size:      AvatarMaxSize,
		Quality:   AvatarWebPQuality,
	},

	// WebP - Small
	&imageoutput.WebPFile{
		Directory: "images/avatars/small/",
		Size:      AvatarSmallSize,
		Quality:   AvatarWebPQuality,
	},
}

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

// FindFileWithExtension tries to test different file extensions.
func FindFileWithExtension(baseName string, dir string, extensions []string) string {
	for _, ext := range extensions {
		if _, err := os.Stat(path.Join(dir, baseName+ext)); !os.IsNotExist(err) {
			return dir + baseName + ext
		}
	}

	return ""
}
