package game

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var simpleFont *truetype.Font
var titleFontFace font.Face
var controlFontFace font.Face

func init() {
	var err error
	if simpleFont, err = truetype.Parse(goregular.TTF); err != nil {
		panic(err)
	}
	titleFontFace = truetype.NewFace(simpleFont, &truetype.Options{Size: 70})
	controlFontFace = truetype.NewFace(simpleFont, &truetype.Options{Size: 30})
}
