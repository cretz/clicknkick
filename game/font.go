package game

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var simpleFont *truetype.Font

func init() {
	var err error
	if simpleFont, err = truetype.Parse(goregular.TTF); err != nil {
		panic(err)
	}
}
