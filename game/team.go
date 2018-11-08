package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

type team struct {
	players []*player
	color   playerColor
	left    bool
}

func newTeam(g *Game, color playerColor, left bool) (t *team, err error) {
	t = &team{players: make([]*player, 11), color: color, left: left}
	for i := 0; i < 11; i++ {
		// First player is goalie
		if t.players[i], err = newPlayer(g.equip, color, rand.Intn(5)+1, i == 0); err != nil {
			break
		}
	}
	// Now let's place them all, 4-4-2 by default
	goalX, _, _ := g.field.goalLine(g, left)
	centerX, centerY := g.field.center(g)
	// Offset the x and y based on left/right and player
	playerW, playerH := t.players[0].bounds()
	centerY -= playerH / 2
	if !left {
		goalX -= playerW
		centerX -= playerW
	}

	xStep := (centerX - goalX) / 3.5
	t.players[0].x, t.players[0].y = goalX, centerY
	const vertPad = 40
	for i := 0; i < 4; i++ {
		t.players[1+i].x, t.players[1+i].y = goalX+xStep, centerY-(vertPad*1.5)+(vertPad*float64(i))
		t.players[5+i].x, t.players[5+i].y = goalX+xStep*2, centerY-(vertPad*1.5)+(vertPad*float64(i))
	}
	t.players[9].x, t.players[9].y = goalX+xStep*3, centerY-vertPad/2
	t.players[10].x, t.players[10].y = goalX+xStep*3, centerY+vertPad/2
	return
}

func (t *team) draw(screen *ebiten.Image, g *Game) {
	for _, player := range t.players {
		player.draw(screen, g)
	}
}
