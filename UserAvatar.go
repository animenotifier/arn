package arn

import (
	"bytes"
	"image"
	"path"
	"time"

	"github.com/blitzprog/imageoutput"
)

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
		Directory: path.Join(Root, "images/avatars/large/"),
		Width:     AvatarMaxSize,
		Height:    AvatarMaxSize,
	},

	// Original - Small
	&imageoutput.OriginalFile{
		Directory: path.Join(Root, "images/avatars/small/"),
		Width:     AvatarSmallSize,
		Height:    AvatarSmallSize,
	},

	// WebP - Large
	&imageoutput.WebPFile{
		Directory: path.Join(Root, "images/avatars/large/"),
		Width:     AvatarMaxSize,
		Height:    AvatarMaxSize,
		Quality:   AvatarWebPQuality,
	},

	// WebP - Small
	&imageoutput.WebPFile{
		Directory: path.Join(Root, "images/avatars/small/"),
		Width:     AvatarSmallSize,
		Height:    AvatarSmallSize,
		Quality:   AvatarWebPQuality,
	},
}

// UserAvatar ...
type UserAvatar struct {
	Extension    string `json:"extension"`
	Source       string `json:"source"`
	LastModified int64  `json:"lastModified"`
}

// SetAvatarBytes accepts a byte buffer that represents an image file and updates the avatar.
func (user *User) SetAvatarBytes(data []byte) error {
	// Decode
	img, format, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return err
	}

	return user.SetAvatar(&imageoutput.MetaImage{
		Image:  img,
		Format: format,
		Data:   data,
	})
}

// SetAvatar sets the avatar to the given MetaImage.
func (user *User) SetAvatar(avatar *imageoutput.MetaImage) error {
	var lastError error

	// Save the different image formats and sizes
	for _, output := range avatarOutputs {
		err := output.Save(avatar, user.ID)

		if err != nil {
			lastError = err
		}
	}

	user.Avatar.Extension = avatar.Extension()
	user.Avatar.LastModified = time.Now().Unix()
	return lastError
}
