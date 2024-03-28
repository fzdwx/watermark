package font

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
)

var defaultFont *truetype.Font

//go:embed default.ttc
var defaultFontTtc []byte

func init() {
	font, _ := truetype.Parse(defaultFontTtc)
	defaultFont = font
}

func Get() *truetype.Font {
	return defaultFont
}
