package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	equip      sprites
	field      *field
	fieldScale ebiten.DrawImageOptions
	ball       *ball
	title      *title
	width      int
	height     int

	team1 *team
	team2 *team

	status gameStatus
}

type gameStatus int

const (
	gameStatusTitle gameStatus = iota
	gameStatusPlay
	gameStatusPause
)

var errorQuit = errors.New("Quit")

func New(width int, height int) (g *Game, err error) {
	g = &Game{width: width, height: height, status: gameStatusPlay}
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
	if g.status == gameStatusTitle {
		switch g.title.update(g) {
		case titleOptionPlay:
			panic("TODO")
		case titleOptionPractice:
			panic("TODO")
		case titleOptionExit:
			return errorQuit
		}
	}
	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	// Draw game components
	g.field.draw(screen, g)
	g.ball.draw(screen, g)
	g.team1.draw(screen, g)
	g.team2.draw(screen, g)
	// Draw the title screen if applicable
	if g.status == gameStatusTitle {
		g.title.draw(screen, g)
	}
	return ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}
