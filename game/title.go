package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/inpututil"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/text"

	"github.com/hajimehoshi/ebiten"
)

type title struct {
	*ebiten.Image
	options []*option
}

type titleOption int

const (
	titleOptionNone titleOption = iota - 1
	titleOptionPlay
	titleOptionPractice
	titleOptionExit
)

func newTitle(width, height int) *title {
	t := &title{options: newOptionSet(width, height, 300, 60, 300, 10, 40, "Play", "Practice", "Exit")}
	// Create a new semi-transparent overlay
	t.Image, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	t.Image.Fill(color.RGBA{0, 0, 0, 200})
	// Draw title, centered
	textWidth := font.MeasureString(titleFontFace, "Click n' Kick").Round()
	text.Draw(t.Image, "Click n' Kick", titleFontFace, width/2-textWidth/2, 150, color.White)
	return t
}

// Returns the clicked title option if any
func (t *title) update(g *Game) titleOption {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i, option := range t.options {
			if option.contains(x, y) {
				return titleOption(i)
			}
		}
	}
	return titleOptionNone
}

func (t *title) draw(screen *ebiten.Image, g *Game) {
	screen.DrawImage(t.Image, nil)
	for _, option := range t.options {
		option.draw(screen, g)
	}
}
