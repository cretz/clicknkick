package game

import (
	"image"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var optionFace font.Face

func init() {
	optionFace = truetype.NewFace(simpleFont, &truetype.Options{Size: 40})
}

type option struct {
	rect  image.Rectangle
	op    *ebiten.DrawImageOptions
	reg   *ebiten.Image
	hover *ebiten.Image
}

func newOptionSet(screenW, screenH, width, height, startY, margin int, paddingTop int, texts ...string) []*option {
	opts := make([]*option, len(texts))
	x := screenW/2 - width/2
	y := startY
	for i, txt := range texts {
		opt := &option{rect: image.Rect(x, y, x+width, y+height), op: newOpTrans(float64(x), float64(y))}
		opts[i] = opt
		textWidth := font.MeasureString(optionFace, txt)
		textX, textY := width/2-textWidth.Round()/2, paddingTop
		// First, create the simple image
		opt.reg, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
		text.Draw(opt.reg, txt, optionFace, textX, textY, color.RGBA{255, 255, 255, 255})
		// Now the hovered
		opt.hover, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
		opt.hover.Fill(color.RGBA{200, 200, 200, 100})
		text.Draw(opt.hover, txt, optionFace, textX, textY, color.RGBA{0, 0, 0, 200})
		y += height + margin
	}
	return opts
}

func (o *option) draw(screen *ebiten.Image, g *Game) {
	img := o.reg
	if o.contains(ebiten.CursorPosition()) {
		img = o.hover
	}
	screen.DrawImage(img, o.op)
}

func (o *option) contains(x, y int) bool {
	return o.rect.Min.X <= x && x < o.rect.Max.X &&
		o.rect.Min.Y <= y && y < o.rect.Max.Y
}
