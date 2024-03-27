package watermark

import (
	"bytes"
	"fmt"
	webp "github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func (c *client) addTextMaskToGif(text string, source io.Reader, option *TextMarkOption) ([]byte, error) {
	img, err := gif.DecodeAll(source)
	if err != nil {
		return nil, fmt.Errorf("encode gif: %w", err)
	}

	newGIF := &gif.GIF{}
	for i, frame := range img.Image {
		newImg := image.NewPaletted(frame.Bounds(), frame.Palette)
		draw.Draw(newImg, newImg.Bounds(), frame, frame.Bounds().Min, draw.Src)

		c.draw(newImg, text, option)

		newGIF.Image = append(newGIF.Image, newImg)
		newGIF.Delay = append(newGIF.Delay, img.Delay[i])
		newGIF.Disposal = append(newGIF.Disposal, img.Disposal[i])
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, newGIF)
	if err != nil {
		return nil, fmt.Errorf("encode to gif: %w", err)
	}
	return buf.Bytes(), nil

}

type TextMarkOption struct {
	Font         *truetype.Font
	FontSize     float64
	TextColor    color.Color
	Dpi          float64
	StepX, StepY int
	Skew         float64
}

func DefaultTextMarkOption() *TextMarkOption {
	return &TextMarkOption{
		Font:      defaultFont,
		FontSize:  12,
		TextColor: color.RGBA{R: 128, G: 128, B: 128, A: 80},
		Dpi:       150,
		StepX:     100,
		StepY:     100,
		Skew:      22.5,
	}
}

func (c *client) addTextMaskToImage(text string, source io.Reader, format ImageFormat, option *TextMarkOption) ([]byte, error) {
	img, _, err := image.Decode(source)
	if err != nil {
		return nil, fmt.Errorf("decode images: %w", err)
	}

	newImg := image.NewNRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Src)
	c.draw(newImg, text, option)

	var buf bytes.Buffer
	switch format {
	case ImageFormatJpeg, ImageFormatJpg:
		err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 80})
	case ImageFormatPng:
		err = png.Encode(&buf, newImg)
	case ImageFormatWebp:
		err = webp.Encode(&buf, newImg, &webp.Options{Quality: 80})
	default:
		return nil, fmt.Errorf("unsuport image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("encode image: %w", err)
	}
	return buf.Bytes(), nil
}

// text 水印文字
// font 字体相关
// watermarkColor 水印颜色
// stepX and stepY 是 x轴和y轴的步长
// dpi dpi
// skew 倾斜角度
func (c *client) draw(
	newImg draw.Image,
	text string,
	option *TextMarkOption,
) {
	textImg := c.newTextImg(text, option.Font, option.FontSize, image.NewUniform(option.TextColor), option.Dpi)
	for y := -option.StepY; y <= newImg.Bounds().Max.Y+option.StepY; y += option.StepY {
		for x := -option.StepX; x <= newImg.Bounds().Max.X+option.StepX; x += option.StepX {
			offsetX := 0
			if (y/option.StepY)%2 == 1 {
				offsetX = option.StepX / 2
			}

			rotated := imaging.Rotate(textImg, option.Skew, image.Transparent)
			draw.Draw(newImg, rotated.Bounds().Add(image.Pt(x+offsetX, y)), rotated, image.Pt(0, 0), draw.Over)
		}
	}
}

func (c *client) newTextImg(text string, font *truetype.Font, fontSize float64, watermarkColor *image.Uniform, dpi float64) image.Image {
	f := freetype.NewContext()
	f.SetDPI(dpi)
	f.SetFont(font)
	f.SetFontSize(fontSize)
	f.SetSrc(watermarkColor)

	textImg := image.NewRGBA(image.Rect(0, 0, len(text)*int(f.PointToFixed(12)>>6), int(f.PointToFixed(12*1.5)>>6)))
	f.SetClip(textImg.Bounds())
	f.SetDst(textImg)

	pt := freetype.Pt(0, int(f.PointToFixed(12)>>6))
	_, _ = f.DrawString(text, pt)
	return textImg
}
