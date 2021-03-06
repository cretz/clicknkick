package resources

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/hajimehoshi/ebiten"
)

func Load(name string) ([]byte, error) {
	// TODO: embed ala go-bindata
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("Unable to get current dir")
	}
	return ioutil.ReadFile(path.Join(path.Dir(filename), name))
}

func LoadImage(name string) (image.Image, error) {
	byts, err := Load(name)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(byts))
	return img, err
}

func LoadEbitenImage(name string) (*ebiten.Image, error) {
	img, err := LoadImage(name)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func MustLoadEbitenImage(name string) *ebiten.Image {
	img, err := LoadEbitenImage(name)
	if err != nil {
		panic(err)
	}
	return img
}
