package game

import (
	"github.com/cretz/clicknkick/game/resources"
	"github.com/hajimehoshi/ebiten"
)

type selectReticle struct {
	*ebiten.Image
	leftOffset, topOffset float64
	tickCounter           float64
	rotDeg                float64
}

func newSelectReticle() (*selectReticle, error) {
	img, err := resources.LoadEbitenImage("crosshair131.png")
	if err != nil {
		return nil, err
	}
	w, h := img.Size()
	return &selectReticle{Image: img, leftOffset: float64(w) / 2, topOffset: float64(h) / 2}, nil
}

const selectReticleTicksPerDegreeChange = 10
const selectReticleDegreeChange = 10

// x and y are the center
func (s *selectReticle) draw(screen *ebiten.Image, g *Game, x, y float64, selected bool) {
	s.tickCounter++
	s.rotDeg = (s.tickCounter / selectReticleTicksPerDegreeChange) * selectReticleDegreeChange
	if s.rotDeg >= 360 {
		s.rotDeg = 0
		s.tickCounter = 0
	}
	op := newOpTrans(-s.leftOffset, -s.topOffset)
	op.GeoM.Rotate(degToRad(s.rotDeg))
	op.GeoM.Translate(x, y)
	// When selected, change the color
	if selected {
		op.ColorM.Scale(1.7, 1.7, 0.5, 1)
	}
	screen.DrawImage(s.Image, op)
}
