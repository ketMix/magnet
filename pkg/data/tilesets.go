package data

import (
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var tilesets = make(map[string]TileSet)

type TileSet struct {
	OpenImage    *ebiten.Image
	OpenImage2   *ebiten.Image
	BlockedImage *ebiten.Image
}

func LoadTileSet(n string) (TileSet, error) {
	if t, ok := tilesets[n]; ok {
		return t, nil
	}
	t := TileSet{}
	if img, err := ReadImage(path.Join(n, "open.png")); err == nil {
		t.OpenImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}
	if img, err := ReadImage(path.Join(n, "open2.png")); err == nil {
		t.OpenImage2 = ebiten.NewImageFromImage(img)
	} else {
		t.OpenImage2 = t.OpenImage
	}
	if img, err := ReadImage(path.Join(n, "blocked.png")); err == nil {
		t.BlockedImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}

	tilesets[n] = t

	return t, nil
}
