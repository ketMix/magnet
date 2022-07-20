package data

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
		return c.IsHit()
	}

	return false
}

func (c *Clickable) IsHit() bool {
	cursorX, cursorY := ebiten.CursorPosition()
	minX, maxX := c.x-c.image.Bounds().Dx()/2, c.x+c.image.Bounds().Dx()/2
	minY, maxY := c.y-c.image.Bounds().Dy()/2, c.y+c.image.Bounds().Dy()/2
	if int(minX) < cursorX && cursorX < int(maxX) {
		if int(minY) < cursorY && cursorY < int(maxY) {
			return true
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
	code             string
	text             *string
	OffsetX, OffsetY int // Forgive me.
	Active           bool
	Bold             bool
	Underline        bool
	Hover            bool
	isHovered        bool
}

func NewButton(x, y int, code string, onClick func()) *Button {
	txtString := GiveMeString(code)
	bounds := text.BoundString(NormalFace, txtString)
	bgImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())

	return &Button{
		Clickable: Clickable{
			x:       x,
			y:       y,
			image:   bgImage,
			onClick: onClick,
		},
		code: code,
		text: &txtString,
	}
}

func NewImageButton(x, y int, image *ebiten.Image, onClick func()) *Button {
	return &Button{
		Clickable: Clickable{
			x:       x,
			y:       y,
			image:   image,
			onClick: onClick,
		},
	}
}

func (b *Button) Update() {
	if b.IsClicked() {
		b.onClick()
		// This isn't the right thing to do, but it's easier to always turn the cursor back on click. (this prevents the pointer cursor being set when a button causes a travel)
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		return
	}

	if b.Hover {
		if b.IsHit() {
			if !b.isHovered {
				b.isHovered = true
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
			}
		} else {
			if b.isHovered {
				b.isHovered = false
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
		}
	}
}

func (b *Button) IsClicked() bool {
	if b.text == nil {
		return b.Clickable.IsClicked()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return b.IsHit()
	}

	return false
}

func (b *Button) IsHit() bool {
	if b.text == nil {
		return b.Clickable.IsHit()
	}

	bounds := text.BoundString(NormalFace, *b.text)
	cursorX, cursorY := ebiten.CursorPosition()
	minX, maxX := b.OffsetX+b.x-bounds.Dx()/2, b.OffsetX+b.x+bounds.Dx()/2
	minY, maxY := b.OffsetY+b.y-bounds.Dy()/2, b.OffsetY+b.y+bounds.Dy()/2
	if int(minX) < cursorX && cursorX < int(maxX) {
		if int(minY) < cursorY && cursorY < int(maxY) {
			return true
		}
	}
	return false
}

func (b *Button) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	if b.text == nil {
		b.Clickable.Draw(screen, screenOp)
		return
	}

	b.OffsetX = int(screenOp.GeoM.Element(0, 2))
	b.OffsetY = int(screenOp.GeoM.Element(1, 2))
	c := color.RGBA{255, 255, 255, 255}
	if b.Active {
		c = color.RGBA{255, 255, 0, 255}
	}

	face := NormalFace
	if b.Bold {
		face = BoldFace
	}

	bounds := DrawStaticTextByCode(
		b.code,
		face,
		b.OffsetX+(b.x),
		b.OffsetY+(b.y),
		c,
		screen,
		true,
	)

	if b.Underline {
		ebitenutil.DrawLine(
			screen, // Wanna see a magic (number) trick?
			float64(b.OffsetX+(b.x)-bounds.Dx()/2)+4,
			float64(b.OffsetY+(b.y)+bounds.Dy())-3,
			float64(b.OffsetX+(b.x)+bounds.Dx()/2)+4,
			float64(b.OffsetY+(b.y)+bounds.Dy())-3,
			c,
		)
	}
}

func (b *Button) Text() string {
	return *b.text
}
