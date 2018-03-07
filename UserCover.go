package arn

import (
	"bytes"
	"image"
	"time"

	"github.com/animenotifier/arn/imageoutput"
)

// CoverMaxSize is the maximum size for covers.
const CoverMaxSize = 1920

// CoverSmallSize is the size used for mobile phones.
const CoverSmallSize = 640

// CoverWebPQuality is the WebP quality of cover images.
const CoverWebPQuality = AvatarWebPQuality

// CoverJPEGQuality is the JPEG quality of cover images.
const CoverJPEGQuality = CoverWebPQuality

// Define the cover image outputs
var coverImageOutputs = []imageoutput.Output{
	// JPEG - Large
	&imageoutput.JPEGFile{
		Directory: "images/covers/large/",
		Size:      CoverMaxSize,
		Quality:   CoverJPEGQuality,
	},

	// JPEG - Small
	&imageoutput.JPEGFile{
		Directory: "images/covers/small/",
		Size:      CoverSmallSize,
		Quality:   CoverJPEGQuality,
	},

	// WebP - Large
	&imageoutput.WebPFile{
		Directory: "images/covers/large/",
		Size:      CoverMaxSize,
		Quality:   CoverWebPQuality,
	},

	// WebP - Small
	&imageoutput.WebPFile{
		Directory: "images/covers/small/",
		Size:      CoverSmallSize,
		Quality:   CoverWebPQuality,
	},
}

// UserCover ...
type UserCover struct {
	Extension    string `json:"extension"`
	LastModified int64  `json:"lastModified"`
}

// SetCoverBytes accepts a byte buffer that represents an image file and updates the cover image.
func (user *User) SetCoverBytes(data []byte) error {
	// Decode
	img, format, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return err
	}

	return user.SetCover(&imageoutput.MetaImage{
		Image:  img,
		Format: format,
		Data:   data,
	})
}

// SetCover sets the cover image to the given MetaImage.
func (user *User) SetCover(cover *imageoutput.MetaImage) error {
	var lastError error

	// Save the different image formats and sizes
	for _, output := range coverImageOutputs {
		err := output.Save(cover, user.ID)

		if err != nil {
			lastError = err
		}
	}

	user.Cover.Extension = cover.Extension()
	user.Cover.LastModified = time.Now().Unix()
	return lastError
}
