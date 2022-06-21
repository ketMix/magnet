package data

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Images for drawing lines.
	EmptyImage    *ebiten.Image
	EmptySubImage *ebiten.Image

	// Images
	wallImage                                  *ebiten.Image
	turretNegativeImage                        *ebiten.Image
	turretPositiveImage                        *ebiten.Image
	spawnerImage, spawnerShadowImage           *ebiten.Image
	spawnerPositiveImage, spawnerNegativeImage *ebiten.Image
	toolSlotImage, toolSlotActiveImage         *ebiten.Image
	toolDestroyImage                           *ebiten.Image
	toolGunImage                               *ebiten.Image
	projecticlePositiveImage                   *ebiten.Image
	projecticleNegativeImage                   *ebiten.Image
	projecticleNeutralImage                    *ebiten.Image

	// SFX
	turretPlaceSound *Sound
)

// Images is the map of all loaded images.
var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

// GetImage returns the image matching the given file name. IT ALSO LOADS IT.
func GetImage(p string) (*ebiten.Image, error) {
	if v, ok := Images[p]; ok {
		return v, nil
	}
	if img, err := ReadImage(p); err != nil {
		return nil, err
	} else {
		eimg := ebiten.NewImageFromImage(img)
		Images[p] = eimg
		return eimg, nil
	}
}

// Sounds is the map of all loaded sounds.
var Sounds map[string]*Sound = make(map[string]*Sound)

// GetSound returns the sound matching the given file name. IT ALSO LOADS IT.
func GetSound(p string) (*Sound, error) {
	if v, ok := Sounds[p]; ok {
		return v, nil
	}
	if snd, err := ReadSound(p); err != nil {
		return nil, err
	} else {
		Sounds[p] = snd
		return snd, nil
	}
}

// LoadData loads some data.
func LoadData() error {
	//
	EmptyImage = ebiten.NewImage(3, 3)
	EmptySubImage = EmptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	EmptyImage.Fill(color.White)

	return nil
}
