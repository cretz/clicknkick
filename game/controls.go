package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

type controls struct {
}

type controlOption int

const (
	controlOptionNone controlOption = iota - 1
	controlOptionTurnComplete
)

func newControls() *controls {
	return &controls{}
}

func (*controls) update(g *Game) controlOption {
	if g.runningTurnPercent == 0 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := cursorPos()
		if x >= 0 && x <= 250 && y >= float64(g.height)-50 && y <= float64(g.height) {
			return controlOptionTurnComplete
		}
	}
	return controlOptionNone
}

func (*controls) draw(screen *ebiten.Image, g *Game) {
	if g.runningTurnPercent == 0 {
		text.Draw(screen, "Turn Complete", controlFontFace, 10, g.height-20, color.Black)
	}
}
