package data

import (
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var tilesets = make(map[string]TileSet)

type TileSet struct {
	OpenPositiveImage *ebiten.Image
	OpenNeutralImage  *ebiten.Image
	OpenNegativeImage *ebiten.Image
	BlockedImage      *ebiten.Image
}

func LoadTileSet(n string) (TileSet, error) {
	if t, ok := tilesets[n]; ok {
		return t, nil
	}
	t := TileSet{}
	if img, err := ReadImage(path.Join(n, "open-neutral.png")); err == nil {
		t.OpenNeutralImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}
	if img, err := ReadImage(path.Join(n, "open-positive.png")); err == nil {
		t.OpenPositiveImage = ebiten.NewImageFromImage(img)
	} else {
		t.OpenPositiveImage = t.OpenNeutralImage
	}
	if img, err := ReadImage(path.Join(n, "open-negative.png")); err == nil {
		t.OpenNegativeImage = ebiten.NewImageFromImage(img)
	} else {
		t.OpenNegativeImage = t.OpenNeutralImage
	}
	if img, err := ReadImage(path.Join(n, "blocked.png")); err == nil {
		t.BlockedImage = ebiten.NewImageFromImage(img)
	} else {
		return t, err
	}

	tilesets[n] = t

	return t, nil
}
