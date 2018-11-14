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
	prevX, prevY := p.x, p.y
	for i, nextX := range p.nextX {
		nextY := p.nextY[i]
		drawLineDot(screen, prevX, prevY, nextX, nextY, &playerNextLineOp)
		prevX, prevY = nextX, nextY
	}
	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.w/2, -p.h/2)
	op.GeoM.Rotate(degToRad(p.currDir(g) - 90))
	op.GeoM.Translate(p.currPos(g))
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

func (p *player) speedFactor() float64 {
	// Just the distance between current and next distance for now
	if len(p.nextX) == 0 {
		return 0
	}
	return math.Abs(distance(p.x, p.y, p.nextX[0], p.nextY[0]))
}

func (p *player) currPos(g *Game) (x, y float64) {
	if len(p.nextX) == 0 || g.runningTurnPercent == 0 {
		return p.x, p.y
	}
	nextX, nextY := p.nextX[0], p.nextY[0]
	return p.x + ((nextX - p.x) / 100 * g.runningTurnPercent), p.y + ((nextY - p.y) / 100 * g.runningTurnPercent)
}

// In degrees
func (p *player) currDir(g *Game) float64 {
	// Face the ball or the next dir
	if len(p.nextX) == 0 {
		return radToDeg(math.Atan2(g.ball.x-p.x, p.y-g.ball.y))
	}
	return radToDeg(math.Atan2(p.nextX[0]-p.x, p.y-p.nextY[0]))
}

func (p *player) putBallAtFeet(g *Game) {
	dirRad := degToRad(p.currDir(g))
	currX, currY := p.currPos(g)
	g.ball.x = currX + p.w*math.Cos(dirRad)
	g.ball.y = currY + p.w*math.Sin(dirRad)
}
