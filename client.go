package watermark

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
	"io"
)

type Client interface {

	// AddTextMask 添加文字水印
	// text 水印文字
	// source 图片源
	// format 图片格式
	// option 水印参数 see DefaultTextMarkOption
	AddTextMask(text string, source io.Reader, format ImageFormat, option *TextMarkOption) ([]byte, error)
}

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

func (c *client) AddTextMask(text string, source io.Reader, format ImageFormat, option *TextMarkOption) ([]byte, error) {
	if option == nil {
		option = DefaultTextMarkOption()
	}
	switch format {
	case ImageFormatGif:
		return c.addTextMaskToGif(text, source, option)
	default:
		return c.addTextMaskToImage(text, source, format, option)
	}
}
