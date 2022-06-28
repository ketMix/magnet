package data

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Implements UIComponent interface
type Clickable struct {
	image   *ebiten.Image
	x, y    int
	onClick func()
}

func (c *Clickable) SetPos(x, y int) {
	c.x = x
	c.y = y
}

func (c *Clickable) Image() *ebiten.Image {
	return c.image
}

func (c *Clickable) Update() {
	if c.IsClicked() {
		c.onClick()
	}
}

func (c *Clickable) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if c.image == nil {
		return
	}
	screenOp.GeoM.Translate(
		float64(c.x-c.image.Bounds().Dx()/2),
		float64(c.y-c.image.Bounds().Dy()/2),
	)
	screen.DrawImage(c.image, screenOp)
}

func (c *Clickable) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		minX, maxX := c.x-c.image.Bounds().Dx()/2, c.x+c.image.Bounds().Dx()/2
		minY, maxY := c.y-c.image.Bounds().Dy()/2, c.y+c.image.Bounds().Dy()/2
		if int(minX) < cursorX && cursorX < int(maxX) {
			if int(minY) < cursorY && cursorY < int(maxY) {
				return true
			}
		}
	}

	return false
}

type BGMIcon struct {
	Clickable
}

func NewBGMIcon() *BGMIcon {
	image, err := GetImage("ui/bgm.png")
	if err != nil {
		return nil
	}
	return &BGMIcon{
		Clickable: Clickable{
			image: image,
			onClick: func() {
				BGM.ToggleMute()
			},
		},
	}
}

func (bgm *BGMIcon) Update() {
	if bgm.IsClicked() {
		bgm.onClick()
	}
}

func (bgm *BGMIcon) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if BGM.Muted {
		screenOp.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
	}
	bgm.Clickable.Draw(screen, screenOp)
}

type SFXIcon struct {
	Clickable
}

func NewSFXIcon() *SFXIcon {
	image, err := GetImage("ui/sfx.png")
	if err != nil {
		return nil
	}
	return &SFXIcon{
		Clickable: Clickable{
			image: image,
			onClick: func() {
				SFX.ToggleMute()
			},
		},
	}
}

func (sfx *SFXIcon) Update() {
	if sfx.IsClicked() {
		sfx.onClick()
	}
}

func (sfx *SFXIcon) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if SFX.Muted {
		screenOp.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
	}
	sfx.Clickable.Draw(screen, screenOp)
}

const borderWidth = 5

type Button struct {
	Clickable
	text string
}

func NewButton(x, y int, txt string, onClick func()) *Button {
	bounds := text.BoundString(NormalFace, txt)
	bgImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())

	return &Button{
		Clickable: Clickable{
			x:       x,
			y:       y,
			image:   bgImage,
			onClick: onClick,
		},
		text: txt,
	}
}

func (b *Button) Update() {
	if b.IsClicked() {
		b.onClick()
	}
}
func (b *Button) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		bounds := text.BoundString(NormalFace, b.text)
		cursorX, cursorY := ebiten.CursorPosition()
		minX, maxX := b.x-bounds.Dx()/2, b.x+bounds.Dx()/2
		minY, maxY := b.y-bounds.Dy()/2, b.y+bounds.Dy()/2
		if int(minX) < cursorX && cursorX < int(maxX) {
			if int(minY) < cursorY && cursorY < int(maxY) {
				return true
			}
		}
	}

	return false
}

func (b *Button) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	text.Draw(
		screen,
		b.text,
		NormalFace,
		int(b.x)-b.image.Bounds().Dx()/2,
		int(b.y)+b.image.Bounds().Dy()/2,
		color.White,
	)
}
