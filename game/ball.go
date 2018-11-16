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

	selected bool
	hovered  bool

	lastX, lastY float64
	nextX, nextY float64

	off            bool
	offKickerTeam1 bool
}

func newBall(equip sprites) (*ball, error) {
	b := &ball{Image: equip["ball_soccer2.png"], moving: true, lastX: -1, lastY: -1, nextX: -1, nextY: -1}
	if b.Image == nil {
		return nil, fmt.Errorf("Missing soccer ball png")
	}
	return b, nil
}

const ballRotateFrameCount = 10

var ballPendingLineOp ebiten.DrawImageOptions
var ballLastLineOp ebiten.DrawImageOptions
var ballNextLineOp ebiten.DrawImageOptions

func init() {
	// Pending line is redish
	ballPendingLineOp.ColorM.Scale(1.5, 0.5, 0.5, 0.6)
	// Last line is gray/red
	ballLastLineOp.ColorM.Scale(1.2, 0.1, 0.1, 0.3)
	// Next line is redish
	ballNextLineOp.ColorM.Scale(1.5, 0.5, 0.5, 0.8)
}

func (b *ball) draw(screen *ebiten.Image, g *Game) {
	// Draw last and next lines
	if b.lastX >= 0 {
		drawLineDot(screen, b.lastX, b.lastY, b.x, b.y, &ballLastLineOp)
	}
	if b.nextX >= 0 {
		drawLineDot(screen, b.x, b.y, b.nextX, b.nextY, &ballNextLineOp)
	}
	// Draw reticle if selected or hovered
	if b.selected || b.hovered {
		// Certain color for selected
		g.selectReticleOp.ColorM.Reset()
		if b.selected {
			g.selectReticleOp.ColorM.Scale(1.7, 1.7, 0.5, 1)
		}
		drawSelectReticle(screen, b.x, b.y, 0.5, &g.selectReticleOp)
	}
	// Draw ball
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
	b.op.GeoM.Translate(b.currPos(g))
	screen.DrawImage(b.Image, &b.op)
}

func (b *ball) cursorOver() bool {
	x, y := cursorPos()
	w, h := b.Image.Size()
	return x >= b.x-float64(w)/2 && x <= b.x+float64(w)/2 && y >= b.y-float64(h)/2 && y <= b.y+float64(h)/2
}

func (b *ball) advanceTurn(g *Game) {
	if b.nextX >= 0 {
		b.lastX, b.lastY = b.x, b.y
		b.x, b.y = b.nextX, b.nextY
		b.nextX, b.nextY = -1, -1
	} else {
		b.lastX, b.lastY = -1, -1
	}
}

func (b *ball) ballPendingPoint(g *Game) (sourceX, sourceY, destX, destY float64) {
	// Off-screen means not there
	if destX, destY = cursorPos(); destX < 0 || destX > float64(g.width) || destY < 0 || destY > float64(g.height) {
		return -1, -1, -1, -1
	}
	// Otherwise, no max distance
	return b.x, b.y, destX, destY
}

func (b *ball) currPos(g *Game) (x, y float64) {
	if b.nextX < 0 || g.runningTurnPercent == 0 {
		return b.x, b.y
	}
	return b.x + ((b.nextX - b.x) / 100 * g.runningTurnPercent), b.y + ((b.nextY - b.y) / 100 * g.runningTurnPercent)
}

func (b *ball) offField(g *Game) (off, left, corner, goal bool) {
	x, y := b.currPos(g)
	// Entire ball must be over the side
	w, h := b.Image.Size()
	x1, y1, x2, y2 := g.field.fieldBounds(g)
	left = x+float64(w)/2 < x1
	if corner = left || x-float64(w)/2 > x2; corner {
		_, goalTop, goalBottom := g.field.goalLine(g, false)
		goal = y-float64(h)/2 > goalTop && y+float64(h)/2 < goalBottom
	}
	off = corner || y+float64(h)/2 < y1 || y-float64(h)/2 > y2
	return
}
