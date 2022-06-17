package game

import (
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var tilesets = make(map[string]TileSet)

type TileSet struct {
	openImage    *ebiten.Image
	openImage2   *ebiten.Image
	blockedImage *ebiten.Image
}

func loadTileSet(n string) (TileSet, error) {
	if t, ok := tilesets[n]; ok {
		return t, nil
	}
	t := TileSet{}
	if img, err := readImage(path.Join(n, "open.png")); err == nil {
		t.openImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}
	if img, err := readImage(path.Join(n, "open2.png")); err == nil {
		t.openImage2 = ebiten.NewImageFromImage(img)
	} else {
		t.openImage2 = t.openImage
	}
	if img, err := readImage(path.Join(n, "blocked.png")); err == nil {
		t.blockedImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}

	tilesets[n] = t

	return t, nil
}
