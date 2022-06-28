package data

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UIComponent interface {
	SetPos(x, y int)
	Image() *ebiten.Image
	Update()
	Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions)
	IsClicked() bool
}
