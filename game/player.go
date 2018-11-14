package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type player struct {
	// Represents center
	x, y float64
	w, h float64
	body *ebiten.Image
	arm  *ebiten.Image
	leg  *ebiten.Image
	// Neg if not applicable
	lastX, lastY float64
	nextX, nextY []float64
}

type playerColor string

const (
	playerColorRed   playerColor = "Red"
	playerColorGreen             = "Green"
	playerColorBlue              = "Blue"
	playerColorWhite             = "White"
)

var playerColors = []playerColor{playerColorRed, playerColorGreen, playerColorBlue, playerColorWhite}

var playerPendingLineOp ebiten.DrawImageOptions
var playerLastLineOp ebiten.DrawImageOptions
var playerNextLineOp ebiten.DrawImageOptions

func init() {
	// Pending line is yellowish
	playerPendingLineOp.ColorM.Scale(1.5, 1.5, 0.7, 0.6)
	// Last line is gray
	playerLastLineOp.ColorM.Scale(0.1, 0.1, 0.1, 0.3)
	// Next line is greener
	playerNextLineOp.ColorM.Scale(0.5, 1.5, 0.5, 0.8)
}

// Num is 1 through 5.
func newPlayer(equip sprites, color playerColor, num int, alt bool) (*player, error) {
	n := num
	if alt {
		n += 5
	}
	armLegOffset := 0
	if num == 2 {
		armLegOffset = 1
	}
	p := &player{
		body:  equip[fmt.Sprintf("character%v (%v).png", color, n)],
		arm:   equip[fmt.Sprintf("character%v (%v).png", color, 11+armLegOffset)],
		leg:   equip[fmt.Sprintf("character%v (%v).png", color, 13+armLegOffset)],
		lastX: -1,
		lastY: -1,
	}
	if p.body == nil || p.arm == nil || p.leg == nil {
		return nil, fmt.Errorf("Missing player piece for color %v, num %v, and off %v", color, num, alt)
	}
	w, h := p.body.Size()
	p.w, p.h = float64(w), float64(h)
	return p, nil
}

func (p *player) draw(screen *ebiten.Image, g *Game) {
	// Draw player chain
	if p.lastX >= 0 {
		drawLineDot(screen, p.lastX, p.lastY, p.x, p.y, &playerLastLineOp)
	}
	currX, currY := p.x, p.y
	prevX, prevY := p.x, p.y
	// Default dir is always to ball
	dir := radToDeg(math.Atan2(g.ball.x-currX, currY-g.ball.y))
	for i, nextX := range p.nextX {
		nextY := p.nextY[i]
		drawLineDot(screen, prevX, prevY, nextX, nextY, &playerNextLineOp)
		prevX, prevY = nextX, nextY
		// The first one overrides the direction
		if i == 0 {
			dir = radToDeg(math.Atan2(nextX-p.x, p.y-nextY))
			// Also, if animating, we need to calc the diff between start and this next one
			if g.runningTurnPercent > 0 {
				currX += (nextX - currX) / 100 * g.runningTurnPercent
				currY += (nextY - currY) / 100 * g.runningTurnPercent
			}
		}
	}
	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.w/2, -p.h/2)
	op.GeoM.Rotate(degToRad(dir - 90))
	op.GeoM.Translate(currX, currY)
	screen.DrawImage(p.body, op)
}

func (p *player) setLeftTop(left, top float64) {
	p.x, p.y = left+p.w/2, top+p.h/2
}

func (p *player) advanceTurn(g *Game) {
	// Move the chain (we do a copy to prevent leaks)
	if len(p.nextX) > 0 {
		p.lastX, p.lastY = p.x, p.y
		p.x, p.y = p.nextX[0], p.nextY[0]
		copy(p.nextX, p.nextX[1:])
		p.nextX = p.nextX[:len(p.nextX)-1]
		copy(p.nextY, p.nextY[1:])
		p.nextY = p.nextY[:len(p.nextY)-1]
	} else {
		p.lastX, p.lastY = -1, -1
	}
}
