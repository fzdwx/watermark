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

func (c *client) addTextMarkToGif(source io.Reader, mark draw.Image, option *TextMarkOption) ([]byte, error) {
	img, err := gif.DecodeAll(source)
	if err != nil {
		return nil, fmt.Errorf("encode gif: %w", err)
	}

	newGIF := &gif.GIF{}
	for i, frame := range img.Image {
		newImg := image.NewPaletted(frame.Bounds(), frame.Palette)
		draw.Draw(newImg, newImg.Bounds(), frame, frame.Bounds().Min, draw.Src)

		c.draw(newImg, mark, option.StepX, option.StepY)

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
		Skew:      10,
	}
}

func (c *client) addTextMarkToImage(source io.Reader, format ImageFormat, mark draw.Image, option *TextMarkOption) ([]byte, error) {
	img, _, err := image.Decode(source)
	if err != nil {
		return nil, fmt.Errorf("decode images: %w", err)
	}

	targetImg := image.NewNRGBA(img.Bounds())
	draw.Draw(targetImg, targetImg.Bounds(), img, img.Bounds().Min, draw.Src)
	c.draw(targetImg, mark, option.StepX, option.StepY)

	var buf bytes.Buffer
	switch format {
	case ImageFormatJpeg, ImageFormatJpg:
		err = jpeg.Encode(&buf, targetImg, &jpeg.Options{Quality: 80})
	case ImageFormatPng:
		err = png.Encode(&buf, targetImg)
	case ImageFormatWebp:
		err = webp.Encode(&buf, targetImg, &webp.Options{Quality: 80})
	default:
		return nil, fmt.Errorf("unsuport image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("encode image: %w", err)
	}
	return buf.Bytes(), nil
}

func (c *client) draw(
	target draw.Image,
	mark draw.Image,
	stepY, stepX int,
) {
	for y := -stepY; y <= target.Bounds().Max.Y+stepY; y += stepY {
		for x := -stepX; x <= target.Bounds().Max.X+stepX; x += stepX {
			offsetX := 0
			if (y/stepY)%2 == 1 {
				offsetX = stepX / 2
			}

			draw.Draw(target, mark.Bounds().Add(image.Pt(x+offsetX, y)), mark, image.Pt(0, 0), draw.Over)
		}
	}
}

func (c *client) newTextImg(text string, opt *TextMarkOption) draw.Image {
	f := freetype.NewContext()
	f.SetDPI(opt.Dpi)
	f.SetFont(opt.Font)
	f.SetFontSize(opt.FontSize)
	f.SetSrc(image.NewUniform(opt.TextColor))

	textImg := image.NewRGBA(image.Rect(0, 0, len(text)*int(f.PointToFixed(12)>>6), int(f.PointToFixed(12*1.5)>>6)))
	f.SetClip(textImg.Bounds())
	f.SetDst(textImg)

	pt := freetype.Pt(0, int(f.PointToFixed(12)>>6))
	_, _ = f.DrawString(text, pt)

	if opt.Skew > 0 {
		return imaging.Rotate(textImg, opt.Skew, image.Transparent)
	}
	return textImg
}
