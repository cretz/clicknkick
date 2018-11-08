package game

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

func newOpTrans(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}

const oneRadianDegree = math.Pi / 180

func degToRad(deg float64) float64 { return deg * oneRadianDegree }
