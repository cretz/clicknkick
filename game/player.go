package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type player struct {
	// Represents left and top
	x, y float64
	body *ebiten.Image
	arm  *ebiten.Image
	leg  *ebiten.Image
	// In degrees, 0 to 360
	dir float64
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
		body: equip[fmt.Sprintf("character%v (%v).png", color, n)],
		arm:  equip[fmt.Sprintf("character%v (%v).png", color, 11+armLegOffset)],
		leg:  equip[fmt.Sprintf("character%v (%v).png", color, 13+armLegOffset)],
	}
	if p.body == nil || p.arm == nil || p.leg == nil {
		return nil, fmt.Errorf("Missing player piece for color %v, num %v, and off %v", color, num, alt)
	}
	return p, nil
}

func (p *player) draw(screen *ebiten.Image, g *Game) {
	screen.DrawImage(p.body, newOpTrans(p.x, p.y))
}

func (p *player) bounds() (width, height float64) {
	// TODO: this changes based on rotation?
	w, h := p.body.Size()
	return float64(w), float64(h)
}
