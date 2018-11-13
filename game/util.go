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

func cursorPos() (x, y float64) {
	xInt, yInt := ebiten.CursorPosition()
	return float64(xInt), float64(yInt)
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
