package arn

import (
	"bytes"
	"image"
	"time"

	"github.com/animenotifier/arn/imageoutput"
)

// UserAvatar ...
type UserAvatar struct {
	Extension    string `json:"extension"`
	Source       string `json:"source"`
	LastModified int64  `json:"lastModified"`
}

// RefreshAvatar ...
func (user *User) RefreshAvatar() {
	// TODO: ...
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

// SetAvatar ...
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
