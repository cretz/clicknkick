package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type ball struct {
	*ebiten.Image
	x, y float64

	frameCount int
	moving     bool
}

func newBall(equip sprites) (*ball, error) {
	b := &ball{Image: equip["ball_soccer2.png"], x: 100, y: 100, moving: true}
	if b.Image == nil {
		return nil, fmt.Errorf("Missing soccer ball png")
	}
	return b, nil
}

const ballRotateFrameCount = 20

func (b *ball) draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}

	b.frameCount++
	if b.frameCount < ballRotateFrameCount && b.moving {
		op.GeoM.Rotate(degToRad(90))
		w, _ := b.Image.Size()
		op.GeoM.Translate(float64(w), 0)
	} else if b.frameCount > ballRotateFrameCount*2 {
		b.frameCount = 0
	}

	op.GeoM.Concat(g.fieldScale.GeoM)
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.Image, op)
}