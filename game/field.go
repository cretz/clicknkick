package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/cretz/clicknkick/game/resources"
	"github.com/hajimehoshi/ebiten"
)

const fieldTileSize = 64

type field struct {
	*ebiten.Image
	xTiles, yTiles int
}

func newField(xTiles, yTiles int) (*field, error) {
	f := &field{xTiles: xTiles, yTiles: yTiles}
	f.Image, _ = ebiten.NewImage(xTiles*fieldTileSize, yTiles*fieldTileSize, ebiten.FilterDefault)
	maxX, maxY := float64(xTiles), float64(yTiles)
	centerX, centerY := f.centerTile()
	// Load tiles (can allow resource to change)
	tiles, err := resources.LoadEbitenImage("groundGrass_mownWide.png")
	if err != nil {
		return nil, err
	}
	// We're going to draw the field with 1 tile padding all around
	// Draw grass
	tile := f.tileImg(tiles, 0, 0)
	for x := 0.0; x < maxX; x++ {
		for y := 0.0; y < maxY; y++ {
			f.drawTileImg(tile, x, y)
		}
	}
	// Corners
	f.drawTile(tiles, 5, 0, 1, 1)
	f.drawTile(tiles, 6, 0, maxX-2, 1)
	f.drawTile(tiles, 5, 1, 1, maxY-2)
	f.drawTile(tiles, 6, 1, maxX-2, maxY-2)
	// Connect corners
	for x := 2.0; x < maxX-2; x++ {
		f.drawTileImg(f.tileImg(tiles, 2, 0), x, 1)
		f.drawTileImg(f.tileImg(tiles, 2, 3), x, maxY-2)
	}
	for y := 2.0; y < maxY-2; y++ {
		f.drawTileImg(f.tileImg(tiles, 1, 1), 1, y)
		f.drawTileImg(f.tileImg(tiles, 3, 1), centerX, y)
		f.drawTileImg(f.tileImg(tiles, 4, 1), maxX-2, y)
	}
	// Mid connectors
	f.drawTile(tiles, 3, 0, centerX, 1)
	f.drawTile(tiles, 3, 3, centerX, maxY-2)
	// Middle circle clockwise from top center
	f.drawTile(tiles, 7, 3, centerX, centerY-1)
	f.drawTile(tiles, 9, 0, centerX+1, centerY-1)
	f.drawTile(tiles, 9, 1, centerX+1, centerY)
	f.drawTile(tiles, 9, 2, centerX+1, centerY+1)
	f.drawTile(tiles, 9, 3, centerX, centerY+1)
	f.drawTile(tiles, 7, 2, centerX-1, centerY+1)
	f.drawTile(tiles, 7, 1, centerX-1, centerY)
	f.drawTile(tiles, 7, 0, centerX-1, centerY-1)
	// Middle dot
	f.drawTile(tiles, 2, 1, centerX, centerY)
	// Penalty box left
	f.drawTile(tiles, 1, 2, 1, centerY-2)
	f.drawTile(tiles, 11, 2, 2, centerY-2)
	f.drawTile(tiles, 4, 1, 2, centerY-1)
	f.drawTile(tiles, 4, 1, 2, centerY)
	f.drawTile(tiles, 4, 1, 2, centerY+1)
	f.drawTile(tiles, 11, 3, 2, centerY+2)
	f.drawTile(tiles, 1, 2, 1, centerY+2)
	// Penalty box right
	f.drawTile(tiles, 4, 2, maxX-2, centerY-2)
	f.drawTile(tiles, 10, 2, maxX-3, centerY-2)
	f.drawTile(tiles, 1, 1, maxX-3, centerY-1)
	f.drawTile(tiles, 1, 1, maxX-3, centerY)
	f.drawTile(tiles, 1, 1, maxX-3, centerY+1)
	f.drawTile(tiles, 10, 3, maxX-3, centerY+2)
	f.drawTile(tiles, 4, 2, maxX-2, centerY+2)
	// Elements for goals
	elems, err := resources.LoadEbitenImage("elements.png")
	if err != nil {
		return nil, err
	}
	// Goal left
	f.drawTile(elems, 4, 3, 0.12, centerY-1)
	f.drawTile(elems, 4, 4, 0.12, centerY)
	f.drawTile(elems, 4, 5, 0.12, centerY+1)
	// Goal right
	f.drawTile(elems, 8, 3, maxX-1.12, centerY-1)
	f.drawTile(elems, 8, 4, maxX-1.12, centerY)
	f.drawTile(elems, 8, 5, maxX-1.12, centerY+1)
	// TODO:
	// * penalty arches
	// * penalty spots
	// * goal boxes
	return f, nil
}

func (f *field) tileImg(tiles *ebiten.Image, sx, sy int) *ebiten.Image {
	return tiles.SubImage(image.Rect(sx*fieldTileSize, sy*fieldTileSize,
		(sx+1)*fieldTileSize, (sy+1)*fieldTileSize)).(*ebiten.Image)
}

func (f *field) drawTile(tiles *ebiten.Image, sx, sy int, dx, dy float64) {
	f.drawTileImg(f.tileImg(tiles, sx, sy), dx, dy)
}

func (f *field) drawTileImg(img *ebiten.Image, dx, dy float64) {
	f.DrawImage(img, newOpTrans(dx*fieldTileSize, dy*fieldTileSize))
}

func (f *field) draw(screen *ebiten.Image, g *Game) {
	screen.DrawImage(f.Image, &g.fieldScale)
	centerX, centerY := f.center(g)
	ebitenutil.DrawLine(screen, centerX, centerY, centerX+5, centerY+5, color.Black)
}

func (f *field) fieldPos(tileX, tileY float64, g *Game) (x, y float64) {
	return g.fieldScale.GeoM.Apply(tileX*fieldTileSize, tileY*fieldTileSize)
}

func (f *field) centerTile() (x, y float64) {
	return float64(f.xTiles)/2 - 0.5, float64(f.yTiles)/2 - 0.5
}

func (f *field) center(g *Game) (x, y float64) {
	cX, cY := f.centerTile()
	return f.fieldPos(cX+.5, cY+.5, g)
}

func (f *field) goalLine(g *Game, left bool) (x, top, bottom float64) {
	_, centerY := f.centerTile()
	if left {
		x, top = f.fieldPos(1.12, centerY-1, g)
	} else {
		x, top = f.fieldPos(float64(f.xTiles)-1.12, centerY-1, g)
	}
	_, bottom = f.fieldPos(0, centerY+1, g)
	return
}

func (f *field) size(g *Game) (x, y float64) {
	return f.fieldPos(float64(f.xTiles), float64(f.yTiles), g)
}
