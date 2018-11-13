package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type team struct {
	players []*player
	color   playerColor
	left    bool

	selectedPlayer         int
	selectedPlayerNewChain bool
	hoveredPlayer          int
}

func newTeam(g *Game, color playerColor, left bool) (t *team, err error) {
	t = &team{
		players:        make([]*player, 11),
		color:          color,
		left:           left,
		selectedPlayer: -1,
		hoveredPlayer:  -1,
	}
	for i := 0; i < 11; i++ {
		// First player is goalie
		if t.players[i], err = newPlayer(g.equip, color, rand.Intn(5)+1, i == 0); err != nil {
			break
		}
		// Face center
		if left {
			t.players[i].dir = 90
		} else {
			t.players[i].dir = 270
		}
	}
	// Now let's place them all, 4-4-2 by default
	goalX, _, _ := g.field.goalLine(g, left)
	centerX, centerY := g.field.center(g)
	// Move the goalX a bit to not center the player on it
	if w := t.players[0].w; left {
		goalX += w / 2
	} else {
		goalX -= w / 2
	}
	xStep := (centerX - goalX) / 3.5
	t.players[0].x, t.players[0].y = goalX, centerY
	const vertPad = 80
	for i := 0; i < 4; i++ {
		t.players[1+i].x, t.players[1+i].y = goalX+xStep, centerY-(vertPad*1.5)+(vertPad*float64(i))
		t.players[5+i].x, t.players[5+i].y = goalX+xStep*2, centerY-(vertPad*1.5)+(vertPad*float64(i))
	}
	t.players[9].x, t.players[9].y = goalX+xStep*3, centerY-vertPad/2
	t.players[10].x, t.players[10].y = goalX+xStep*3, centerY+vertPad/2
	return
}

func (t *team) update(g *Game, myTeam bool) {
	if !myTeam {
		return
	}
	// If there is a selected player, right click unselects, otherwise can click on hovered player
	if t.selectedPlayer >= 0 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		t.selectedPlayer = -1
	}
	if t.selectedPlayer == -1 {
		t.hoveredPlayer, _ = t.playerAtCursor()
		if t.hoveredPlayer >= 0 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			t.selectedPlayer = t.hoveredPlayer
			t.selectedPlayerNewChain = true
			t.hoveredPlayer = -1
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// When mouse is clicked in screen, it's a new point
		if p, _, _, destX, destY := t.selectedPlayerPendingPoint(g); destX >= 0 {
			if t.selectedPlayerNewChain {
				t.selectedPlayerNewChain = false
				p.nextX = nil
				p.nextY = nil
			}
			p.nextX = append(p.nextX, destX)
			p.nextY = append(p.nextY, destY)
		}
	}
}

func (t *team) draw(screen *ebiten.Image, g *Game, myTeam bool) {
	for _, player := range t.players {
		player.draw(screen, g)
	}
	if myTeam {
		if t.selectedPlayer >= 0 {
			// If the player is selected, mouse cursor is where they can move (limited by max)
			p, sourceX, sourceY, destX, destY := t.selectedPlayerPendingPoint(g)
			g.selectReticle.draw(screen, g, p.x, p.y, true)
			if sourceX >= 0 {
				ebitenutil.DrawLine(screen, sourceX, sourceY, destX, destY, color.Black)
			}
		} else if t.hoveredPlayer >= 0 {
			p := t.players[t.hoveredPlayer]
			g.selectReticle.draw(screen, g, p.x, p.y, false)
		}
	}
}

func (t *team) selectedPlayerPendingPoint(g *Game) (p *player, sourceX, sourceY, destX, destY float64) {
	p = t.players[t.selectedPlayer]
	// Off-screen means not there
	if destX, destY = cursorPos(); destX < 0 || destX > float64(g.width) || destY < 0 || destY > float64(g.height) {
		return p, -1, -1, -1, -1
	}
	sourceX, sourceY = p.x, p.y
	if !t.selectedPlayerNewChain {
		sourceX, sourceY = p.nextX[len(p.nextX)-1], p.nextY[len(p.nextY)-1]
	}
	// Max out the distance
	if dist := math.Abs(distance(sourceX, sourceY, destX, destY)); dist > g.maxPlayerMove {
		off := g.maxPlayerMove / dist
		destX = sourceX + ((destX - sourceX) * off)
		destY = sourceY + ((destY - sourceY) * off)
	}
	return
}

func (t *team) playerAtCursor() (playerIndex int, fromCenter float64) {
	x, y := cursorPos()
	// Closest player by combined from-center amounts
	playerIndex, fromCenter = -1, math.MaxFloat64
	for i, p := range t.players {
		fromCentX, fromCentY := math.Abs(p.x-x), math.Abs(p.y-y)
		if fromCentX <= p.w/2 && fromCentY <= p.h/2 && fromCentX+fromCentY < fromCenter {
			playerIndex = i
			fromCenter = fromCentX + fromCentY
		}
	}
	return
}
