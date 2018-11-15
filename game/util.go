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

var arrowLineImage = resources.MustLoadEbitenImage("arrow-line.png")
var circleImage = resources.MustLoadEbitenImage("circle.png")

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

var selectReticleImage = resources.MustLoadEbitenImage("crosshair131.png")
var selectReticleTickCounter = 0.0

const selectReticleTicksPerDegreeChange = 10
const selectReticleDegreeChange = 10

func drawSelectReticle(dst *ebiten.Image, x, y, scale float64, op *ebiten.DrawImageOptions) {
	selectReticleTickCounter++
	rotDeg := (selectReticleTickCounter / selectReticleTicksPerDegreeChange) * selectReticleDegreeChange
	if rotDeg >= 360 {
		rotDeg = 0
		selectReticleTickCounter = 0
	}
	op.GeoM.Reset()
	w, h := selectReticleImage.Size()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(degToRad(rotDeg))
	op.GeoM.Translate(x, y)
	dst.DrawImage(selectReticleImage, op)
}
