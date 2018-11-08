package game

import (
	"encoding/xml"
	"image"

	"github.com/cretz/clicknkick/game/resources"
	"github.com/hajimehoshi/ebiten"
)

type sprites map[string]*ebiten.Image

func newSprites(sheetName string) (sprites, error) {
	// Load XML
	sheetXml := &struct {
		SubTexture []struct {
			Name   string `xml:"name,attr"`
			X      int    `xml:"x,attr"`
			Y      int    `xml:"y,attr"`
			Width  int    `xml:"width,attr"`
			Height int    `xml:"height,attr"`
		}
	}{}
	if byts, err := resources.Load(sheetName + ".xml"); err != nil {
		return nil, err
	} else if err = xml.Unmarshal(byts, sheetXml); err != nil {
		return nil, err
	}
	// Load PNG
	sheetImg, err := resources.LoadEbitenImage(sheetName + ".png")
	if err != nil {
		return nil, err
	}
	// Chop it up
	s := make(map[string]*ebiten.Image, len(sheetXml.SubTexture))
	for _, tex := range sheetXml.SubTexture {
		s[tex.Name] = sheetImg.SubImage(image.Rect(tex.X, tex.Y, tex.X+tex.Width, tex.Y+tex.Height)).(*ebiten.Image)
	}
	return s, nil
}
