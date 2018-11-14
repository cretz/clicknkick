package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type ball struct {
	*ebiten.Image
	op   ebiten.DrawImageOptions
	x, y float64

	frameCount int
	moving     bool
}

func newBall(equip sprites) (*ball, error) {
	b := &ball{Image: equip["ball_soccer2.png"], moving: true}
	if b.Image == nil {
		return nil, fmt.Errorf("Missing soccer ball png")
	}
	return b, nil
}

const ballRotateFrameCount = 10

func (b *ball) draw(screen *ebiten.Image, g *Game) {
	w, h := b.Image.Size()
	b.op.GeoM.Reset()
	b.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

	b.frameCount++
	if b.frameCount < ballRotateFrameCount && b.moving {
		b.op.GeoM.Rotate(degToRad(90))
	} else if b.frameCount > ballRotateFrameCount*2 {
		b.frameCount = 0
	}

	b.op.GeoM.Concat(g.fieldScale.GeoM)
	b.op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.Image, &b.op)
}
