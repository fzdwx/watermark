package watermark

import (
	"io"
)

type Client interface {

	// AddTextMark add text mark
	// text water text
	// source image source
	// format image format
	// option see DefaultTextMarkOption
	AddTextMark(text string, source io.Reader, format ImageFormat, option *TextMarkOption) ([]byte, error)
}

type Option struct {
	// Output image quality
	// for jpg/jpeg/webp
	// 0~100
	Quality float64
	// Watermark Image gap
	StepX, StepY int
	// Watermark Image rotate angle
	Skew float64
}

// NewClient new watermark client
func NewClient() Client {
	return &client{}
}

type client struct {
}

func (c *client) AddTextMark(text string, source io.Reader, format ImageFormat, option *TextMarkOption) ([]byte, error) {
	if option == nil {
		option = DefaultTextMarkOption()
	}

	img := c.newTextImg(text, option)
	switch format {
	case ImageFormatGif:
		return c.addTextMarkToGif(source, img, option.Option)
	default:
		return c.addTextMarkToImage(source, format, img, option.Option)
	}
}
