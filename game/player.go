package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type player struct {
	body *ebiten.Image
	arm  *ebiten.Image
	leg  *ebiten.Image
	// In degrees, 0 to 360
	dir float64
}

// Color is "Red", "Green", "Blue", or "White". Num is 1 through 5.
func newPlayer(equip sprites, color string, num int, alt bool) (*player, error) {
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
