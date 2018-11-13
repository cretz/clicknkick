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
	team1         *team
	team2         *team
	selectReticle *selectReticle

	iAmTeam1 bool
	practice bool
	status   gameStatus
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
	g = &Game{width: width, height: height, status: gameStatusTitle}
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
	g.ball.setCenter(g.field.center(g))
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
	// Create title
	g.title = newTitle(width, height)
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
		g.team1.update(g, g.iAmTeam1)
		g.team2.update(g, !g.iAmTeam1)
	}
	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	// Draw game components
	g.field.draw(screen, g)
	g.ball.draw(screen, g)
	g.team1.draw(screen, g, g.iAmTeam1)
	if !g.practice {
		g.team2.draw(screen, g, !g.iAmTeam1)
	}
	// Draw the title screen if applicable
	switch g.status {
	case gameStatusTitle:
		g.title.draw(screen, g)
	case gameStatusPlay:
	}
	return ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) myTeam() *team {
	if g.iAmTeam1 {
		return g.team1
	}
	return g.team2
}
