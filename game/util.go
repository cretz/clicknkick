package game

import (
	"github.com/hajimehoshi/ebiten"
)

func newOpTrans(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}
