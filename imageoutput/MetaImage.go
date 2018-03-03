package imageoutput

import (
	"fmt"
	"image"
)

// MetaImage represents a single image and the name of the format.
type MetaImage struct {
	Image  image.Image
	Data   []byte
	Format string
}

// Extension ...
func (avatar *MetaImage) Extension() string {
	switch avatar.Format {
	case "jpg", "jpeg":
		return ".jpg"
	case "png":
		return ".png"
	case "gif":
		return ".gif"
	default:
		return ""
	}
}

// String returns a text representation of the format, width and height.
func (avatar *MetaImage) String() string {
	return fmt.Sprint(avatar.Format, " | ", avatar.Image.Bounds().Dx(), "x", avatar.Image.Bounds().Dy())
}
