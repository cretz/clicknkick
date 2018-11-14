package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	width         int
	height        int
	maxPlayerMove float64

	equip         sprites
	field         *field
	fieldScale    ebiten.DrawImageOptions
	ball          *ball
	title         *title
	controls      *controls
	team1         *team
	team2         *team
	selectReticle *selectReticle

	iAmTeam1           bool
	practice           bool
	status             gameStatus
	runningTurnPercent float64

	ballPossessionPlayerIndex int
	ballPossessionTeam1       bool
}

type gameStatus int

const (
	gameStatusTitle gameStatus = iota
	gameStatusPlay
	gameStatusPause
)

var errorQuit = errors.New("Quit")

const debug = true

func New(width int, height int) (g *Game, err error) {
	g = &Game{width: width, height: height, status: gameStatusTitle, ballPossessionPlayerIndex: -1}
	if debug {
		g.status = gameStatusPlay
		g.iAmTeam1 = true
		g.practice = true
	}
	// Load sprites
	if g.equip, err = newSprites("sheet_charactersEquipment"); err != nil {
		return nil, err
	}
	// Create field
	if g.field, err = newField(21, 10); err != nil {
		return nil, err
	}
	// Scale it
	fieldWidth, fieldHeight := g.field.Size()
	g.fieldScale.GeoM.Scale(float64(width)/float64(fieldWidth), float64(height)/float64(fieldHeight))
	// Create ball
	if g.ball, err = newBall(g.equip); err != nil {
		return nil, err
	}
	g.ball.x, g.ball.y = g.field.center(g)
	// Can only move player so many tiles
	g.maxPlayerMove, _ = g.fieldScale.GeoM.Apply(5*fieldTileSize, 0)
	// Create teams (making sure they don't share a color)
	if g.team1, err = newTeam(g, playerColors[rand.Intn(len(playerColors))], true); err != nil {
		return nil, err
	}
	var team2Color playerColor
	for {
		if team2Color = playerColors[rand.Intn(len(playerColors))]; g.team1.color != team2Color {
			break
		}
	}
	if g.team2, err = newTeam(g, team2Color, false); err != nil {
		return nil, err
	}
	// Create the reticle for selecting
	if g.selectReticle, err = newSelectReticle(); err != nil {
		return nil, err
	}
	// Create title and controls
	g.title = newTitle(width, height)
	g.controls = newControls()
	return
}

func (g *Game) Run() error {
	if err := ebiten.Run(g.tick, g.width, g.height, 1, "Click n' Kick"); err != errorQuit {
		return err
	}
	return nil
}

func (g *Game) tick(screen *ebiten.Image) (err error) {
	err = g.update(screen)
	if err == nil && !ebiten.IsDrawingSkipped() {
		err = g.draw(screen)
	}
	return
}

const secondsOfAnimation = 2

// Have to be within this distance to obtain possession
const maxGainPossessionDistance = 20

func (g *Game) update(screen *ebiten.Image) error {
	switch g.status {
	case gameStatusTitle:
		switch g.title.update(g) {
		case titleOptionPlay:
			panic("TODO")
		case titleOptionPractice:
			g.status = gameStatusPlay
			g.practice = true
			g.iAmTeam1 = true
		case titleOptionExit:
			return errorQuit
		}
	case gameStatusPlay:
		switch g.controls.update(g) {
		case controlOptionTurnComplete:
			// Start the animation
			g.myTeam().selectedPlayer = -1
			g.myTeam().hoveredPlayer = -1
			g.runningTurnPercent = 0.01
		}
		if g.runningTurnPercent >= 100 {
			g.runningTurnPercent = 0
			g.team1.advanceTurn(g)
			g.team2.advanceTurn(g)
			// One last put-at-feet when anim stops if there is a next move
			if p := g.ballPossessionPlayer(); p != nil && len(p.nextX) > 0 {
				p.putBallAtFeet(g)
			}
		} else if g.runningTurnPercent > 0 {
			totalTicksInAnimation := secondsOfAnimation * ebiten.CurrentTPS()
			g.runningTurnPercent += 100 / totalTicksInAnimation
		}
		g.team1.update(g, g.iAmTeam1)
		g.team2.update(g, !g.iAmTeam1)
		// If we're moving, then possession can be updated
		if g.runningTurnPercent > 0 {
			g.updatePossession()
			g.ball.moving = true
			// If there is a possessing player, we need to put the ball at his feet
			if p := g.ballPossessionPlayer(); p != nil && len(p.nextX) > 0 {
				p.putBallAtFeet(g)
			}
		} else {
			g.ball.moving = false
		}
	}
	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	// Draw game components
	g.field.draw(screen, g)
	g.team1.draw(screen, g, g.iAmTeam1)
	if !g.practice {
		g.team2.draw(screen, g, !g.iAmTeam1)
	}
	g.ball.draw(screen, g)
	// Draw the controls if applicable
	switch g.status {
	case gameStatusTitle:
		g.title.draw(screen, g)
	case gameStatusPlay:
		g.controls.draw(screen, g)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nRun Pct: %0.2f", g.runningTurnPercent))
	return ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) myTeam() *team {
	if g.iAmTeam1 {
		return g.team1
	}
	return g.team2
}

func (g *Game) updatePossession() {
	// Find any player within ball range that is moving the slowest
	newSlowest, speedFactor := g.team1.slowestPlayerWithinBallRange(g)
	team1 := true
	if !g.practice {
		team2Player, team2SpeedFactor := g.team2.slowestPlayerWithinBallRange(g)
		if newSlowest == -1 || (team2Player != -1 && team2SpeedFactor < speedFactor) {
			newSlowest, speedFactor = team2Player, team2SpeedFactor
			team1 = false
		}
	}
	// No player? No prob
	if newSlowest == -1 {
		return
	}
	// Someone already possessing?
	if currPlayer := g.ballPossessionPlayer(); currPlayer != nil {
		// Can't steal from my own team and can't steal from slower player
		if team1 == g.ballPossessionTeam1 || currPlayer.speedFactor() < speedFactor {
			return
		}
	}
	fmt.Printf("Changed ball possessor to %v (team1? %v)\n", newSlowest, team1)
	// New slowest takes over
	g.ballPossessionPlayerIndex, g.ballPossessionTeam1 = newSlowest, team1
}

func (g *Game) ballPossessionPlayer() *player {
	if g.ballPossessionPlayerIndex == -1 {
		return nil
	} else if g.ballPossessionTeam1 {
		return g.team1.players[g.ballPossessionPlayerIndex]
	} else {
		return g.team2.players[g.ballPossessionPlayerIndex]
	}
}
