# Watermark

Add text mark to image.


## Install

```shell
go get github.com/fzdwx/watermark
```

## Usage

```go
package main

import (
	"github.com/fzdwx/watermark"
	"os"
	"testing"
)

func TestText(t *testing.T) {
	c := watermark.NewClient()

	f, err := os.Open("demo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	option := watermark.DefaultTextMarkOption()
	mask, err := c.AddTextMask("Hello", f, watermark.ImageFormatJpg, option)
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Create("demo_watermark.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(mask)
}

```

### Before

![before](./.github/demo.jpg)

### After

![after](./.github/demo_watermark.jpg)

