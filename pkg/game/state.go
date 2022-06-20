package game

import "github.com/hajimehoshi/ebiten/v2"

type State interface {
	Init() error
	Dispose() error
	Update() error
	Draw(screen *ebiten.Image)
}
