package arn

import (
	"bytes"
	"image"
	"time"

	"github.com/animenotifier/arn/imageoutput"
)

// CoverMaxWidth is the maximum size for covers.
const CoverMaxWidth = 1920

// CoverMaxHeight is the maximum height for covers.
const CoverMaxHeight = 450

// CoverSmallWidth is the width used for mobile phones.
const CoverSmallWidth = 640

// CoverSmallHeight is the height used for mobile phones.
const CoverSmallHeight = 640

// CoverWebPQuality is the WebP quality of cover images.
const CoverWebPQuality = AvatarWebPQuality

// CoverJPEGQuality is the JPEG quality of cover images.
const CoverJPEGQuality = CoverWebPQuality

// Define the cover image outputs
var coverImageOutputs = []imageoutput.Output{
	// JPEG - Large
	&imageoutput.JPEGFile{
		Directory: "images/covers/large/",
		Width:     CoverMaxWidth,
		Height:    CoverMaxHeight,
		Quality:   CoverJPEGQuality,
	},

	// JPEG - Small
	&imageoutput.JPEGFile{
		Directory: "images/covers/small/",
		Width:     CoverSmallWidth,
		Height:    CoverSmallHeight,
		Quality:   CoverJPEGQuality,
	},

	// WebP - Large
	&imageoutput.WebPFile{
		Directory: "images/covers/large/",
		Width:     CoverMaxWidth,
		Height:    CoverMaxHeight,
		Quality:   CoverWebPQuality,
	},

	// WebP - Small
	&imageoutput.WebPFile{
		Directory: "images/covers/small/",
		Width:     CoverSmallWidth,
		Height:    CoverSmallHeight,
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
