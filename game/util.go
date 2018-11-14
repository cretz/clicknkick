package game

import (
	"math"

	"github.com/cretz/clicknkick/game/resources"

	"github.com/hajimehoshi/ebiten"
)

func newOpTrans(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}

const oneRadianDegree = math.Pi / 180

func degToRad(deg float64) float64 { return deg * oneRadianDegree }
func radToDeg(rad float64) float64 { return rad / oneRadianDegree }

func cursorPos() (x, y float64) {
	xInt, yInt := ebiten.CursorPosition()
	return float64(xInt), float64(yInt)
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

var arrowLineImage *ebiten.Image
var circleImage *ebiten.Image

func init() {
	var err error
	if arrowLineImage, err = resources.LoadEbitenImage("arrow-line.png"); err != nil {
		panic(err)
	}
	if circleImage, err = resources.LoadEbitenImage("circle.png"); err != nil {
		panic(err)
	}
}

func drawLineDot(dst *ebiten.Image, x1, y1, x2, y2 float64, op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	w, h := arrowLineImage.Size()
	op.GeoM.Translate(-float64(w)/2, -float64(h))
	op.GeoM.Scale(1, distance(x1, y1, x2, y2)/float64(h))
	op.GeoM.Rotate(math.Atan2(x2-x1, y1-y2))
	op.GeoM.Translate(x1, y1)
	dst.DrawImage(arrowLineImage, op)

	op.GeoM.Reset()
	w, h = circleImage.Size()
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(x2-float64(w)/4, y2-float64(h)/4)
	dst.DrawImage(circleImage, op)
}
