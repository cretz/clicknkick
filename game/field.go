package game

import (
	"image"

	"github.com/cretz/clicknkick/game/resources"
	"github.com/hajimehoshi/ebiten"
)

const fieldTileSize = 64

type field struct {
	*ebiten.Image
}

func newField(xTiles, yTiles int) (*field, error) {
	f := &field{}
	f.Image, _ = ebiten.NewImage(xTiles*fieldTileSize, yTiles*fieldTileSize, ebiten.FilterDefault)
	maxX, maxY := float64(xTiles), float64(yTiles)
	centerX, centerY := maxX/2-0.5, maxY/2-0.5
	// Load tiles (can allow resource to change)
	goGroundImg, err := resources.LoadImage("groundGrass_mownWide.png")
	if err != nil {
		return nil, err
	}
	tiles, _ := ebiten.NewImageFromImage(goGroundImg, ebiten.FilterDefault)
	// Draw grass
	tile := f.tileImg(tiles, 0, 0)
	for x := 0.0; x < maxX; x++ {
		for y := 0.0; y < maxY; y++ {
			f.drawTileImg(tile, x, y)
		}
	}
	// Corners
	f.drawTile(tiles, 5, 0, 0, 0)
	f.drawTile(tiles, 6, 0, maxX-1, 0)
	f.drawTile(tiles, 5, 1, 0, maxY-1)
	f.drawTile(tiles, 6, 1, maxX-1, maxY-1)
	// Connect corners
	for x := 1.0; x < maxX-1; x++ {
		f.drawTileImg(f.tileImg(tiles, 2, 0), x, 0)
		f.drawTileImg(f.tileImg(tiles, 2, 3), x, maxY-1)
	}
	for y := 1.0; y < maxY-1; y++ {
		f.drawTileImg(f.tileImg(tiles, 1, 1), 0, y)
		f.drawTileImg(f.tileImg(tiles, 3, 1), centerX, y)
		f.drawTileImg(f.tileImg(tiles, 4, 1), maxX-1, y)
	}
	// Mid connectors
	f.drawTile(tiles, 3, 0, centerX, 0)
	f.drawTile(tiles, 3, 3, centerX, maxY-1)
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
	// Goal box left
	f.drawTile(tiles, 1, 2, 0, centerY-2)
	f.drawTile(tiles, 11, 2, 1, centerY-2)
	f.drawTile(tiles, 4, 1, 1, centerY-1)
	f.drawTile(tiles, 4, 1, 1, centerY)
	f.drawTile(tiles, 4, 1, 1, centerY+1)
	f.drawTile(tiles, 11, 3, 1, centerY+2)
	f.drawTile(tiles, 1, 2, 0, centerY+2)
	// Goal box right
	f.drawTile(tiles, 4, 2, maxX-1, centerY-2)
	f.drawTile(tiles, 10, 2, maxX-2, centerY-2)
	f.drawTile(tiles, 1, 1, maxX-2, centerY-1)
	f.drawTile(tiles, 1, 1, maxX-2, centerY)
	f.drawTile(tiles, 1, 1, maxX-2, centerY+1)
	f.drawTile(tiles, 10, 3, maxX-2, centerY+2)
	f.drawTile(tiles, 4, 2, maxX-1, centerY+2)
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
