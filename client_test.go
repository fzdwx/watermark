package watermark

import (
	"os"
	"testing"
	"time"
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

func TestCost(t *testing.T) {
	c := NewClient()

	cost(func() {
		if err := genText(c, "你好"); err != nil {
			t.Fatal(err)
		}
	})

	cost(func() {
		if err := genText(c, "Hello"); err != nil {
			t.Fatal(err)
		}
	})

	cost(func() {
		if err := genText(c, "测试测试测试测试测试测试"); err != nil {
			t.Fatal(err)
		}
	})

}

func genText(c Client, text string) error {
	f, err := os.Open("./.github/img.png")
	if err != nil {
		return err
	}
	defer f.Close()

	option := DefaultTextMarkOption()
	mask, err := c.AddTextMask(text, f, ImageFormatJpg, option)
	if err != nil {
		return err
	}

	f, err = os.CreateTemp("", "*.jpg")
	if err != nil {
		return err
	}
	_, err = f.Write(mask)
	return nil

}

func cost(f func()) {
	now := time.Now()
	f()
	println(time.Since(now).String())
}
