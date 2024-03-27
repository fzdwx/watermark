package watermark

import (
	"os"
	"testing"
)

func TestText(t *testing.T) {
	c := NewClient()

	f, err := os.Open("./.github/demo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	option := DefaultTextMarkOption()
	mask, err := c.AddTextMask("Hello", f, ImageFormatJpg, option)
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Create("./.github/demo_watermark.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(mask)
}
