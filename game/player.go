package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

type player struct {
	// Represents center
	x, y float64
	w, h float64
	body *ebiten.Image
	arm  *ebiten.Image
	leg  *ebiten.Image
	// In degrees, 0 to 360
	dir float64
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
		ebitenutil.DrawLine(screen, p.lastX, p.lastY, p.x, p.y, color.Black)
	}
	prevX, prevY := p.x, p.y
	for i, nextX := range p.nextX {
		ebitenutil.DrawLine(screen, prevX, prevY, nextX, p.nextY[i], color.Black)
		prevX, prevY = nextX, p.nextY[i]
	}
	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.w/2, -p.h/2)
	op.GeoM.Rotate(degToRad(p.dir - 90))
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.body, op)
}

func (p *player) setLeftTop(left, top float64) {
	p.x, p.y = left+p.w/2, top+p.h/2
}
