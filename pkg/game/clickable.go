package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebijam22/pkg/data"
)

type ClickableUI interface {
	SetPos(x, y float64)
	Image() *ebiten.Image
	Update()
	Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions)
	IsClicked() bool
	OnClick()
}

type Clickable struct {
	image *ebiten.Image
	x     float64
	y     float64
}

func (c *Clickable) SetPos(x, y float64) {
	c.x = x
	c.y = y
}

func (c *Clickable) Image() *ebiten.Image {
	return c.image
}

func (c *Clickable) Update() {
	if c.IsClicked() {
		c.OnClick()
	}
}

func (c *Clickable) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if c.image == nil {
		return
	}
	screenOp.GeoM.Translate(
		c.x-float64(c.image.Bounds().Dx())/2,
		c.y-float64(c.image.Bounds().Dy())/2,
	)
	screen.DrawImage(c.image, screenOp)
}

func (c *Clickable) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		minX, maxX := c.x-float64(c.image.Bounds().Dx())/2, c.x+float64(c.image.Bounds().Dx())/2
		minY, maxY := c.y-float64(c.image.Bounds().Dx())/2, c.y+float64(c.image.Bounds().Dx())/2
		if int(minX) < cursorX && cursorX < int(maxX) {
			if int(minY) < cursorY && cursorY < int(maxY) {
				return true
			}
		}
	}

	return false
}

func (c *Clickable) OnClick() {
	return
}

type BGMIcon struct {
	Clickable
}

func NewBGMIcon() *BGMIcon {
	image, err := data.GetImage("ui/bgm.png")
	if err != nil {
		return nil
	}
	return &BGMIcon{
		Clickable: Clickable{
			image: image,
		},
	}
}

func (bgm *BGMIcon) Update() {
	if bgm.IsClicked() {
		bgm.OnClick()
	}
}

func (bgm *BGMIcon) OnClick() {
	data.BGM.ToggleMute()
}

func (bgm *BGMIcon) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if data.BGM.Muted {
		screenOp.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
	}
	bgm.Clickable.Draw(screen, screenOp)
}

type SFXIcon struct {
	Clickable
}

func NewSFXIcon() *SFXIcon {
	image, err := data.GetImage("ui/sfx.png")
	if err != nil {
		return nil
	}
	return &SFXIcon{
		Clickable: Clickable{
			image: image,
		},
	}
}

func (sfx *SFXIcon) Update() {
	if sfx.IsClicked() {
		sfx.OnClick()
	}
}

func (sfx *SFXIcon) OnClick() {
	data.SFX.ToggleMute()
}

func (sfx *SFXIcon) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if data.SFX.Muted {
		screenOp.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
	}
	sfx.Clickable.Draw(screen, screenOp)
}
