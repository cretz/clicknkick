package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	field       *field
	fieldDrawOp ebiten.DrawImageOptions
	width       int
	height      int
}

func New(width int, height int) (g *Game, err error) {
	g = &Game{width: width, height: height}
	// Create field
	if g.field, err = newField(21, 10); err != nil {
		return nil, err
	}
	// Scale it
	fieldWidth, fieldHeight := g.field.Size()
	g.fieldDrawOp.GeoM.Scale(float64(width)/float64(fieldWidth), float64(height)/float64(fieldHeight))
	return
}

func (g *Game) Run() error {
	return ebiten.Run(g.tick, g.width, g.height, 1, "Click n' Kick")
}

func (g *Game) tick(screen *ebiten.Image) (err error) {
	err = g.update(screen)
	if err == nil && !ebiten.IsDrawingSkipped() {
		err = g.draw(screen)
	}
	return
}

func (g *Game) update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	// Draw the field as scaled
	screen.DrawImage(g.field.Image, &g.fieldDrawOp)
	return ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}