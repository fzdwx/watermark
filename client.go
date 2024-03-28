package watermark

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
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

// NewClient new watermark client
func NewClient() Client {
	return &client{}
}

var defaultFont *truetype.Font

//go:embed default.ttc
var defaultFontTtc []byte

func init() {
	font, _ := truetype.Parse(defaultFontTtc)
	defaultFont = font
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
		return c.addTextMarkToGif(source, img, option)
	default:
		return c.addTextMarkToImage(source, format, img, option)
	}
}
