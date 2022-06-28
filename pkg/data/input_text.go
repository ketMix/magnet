package data

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var receivingKeyboardInput map[string]bool = make(map[string]bool)

// Are there any active inputs
func CurrentlyReceivingInput() bool {
	for _, v := range receivingKeyboardInput {
		if v {
			return v
		}
	}
	return false
}

func DrawStaticText(txt string, font font.Face, x, y int, color color.Color, screen *ebiten.Image, shouldCenter bool) {
	var offsetX int
	var offsetY int
	bounds := text.BoundString(font, txt)
	if shouldCenter {
		offsetX = bounds.Dx() / 2
		offsetY = bounds.Dy() / 2
	}
	text.Draw(
		screen,
		txt,
		font,
		x-offsetX,
		y+offsetY,
		color,
	)
}

// Implements UIComponent interface
type TextInput struct {
	Clickable
	label       string
	isActive    bool
	runes       []rune
	data        string
	counter     int
	borderWidth int
	maxLength   int
}

func NewTextInput(label string, x, y int) *TextInput {
	maxLength := 22
	borderWidth := 2
	textSize := text.BoundString(NormalFace, "t")
	outerImage := ebiten.NewImage((textSize.Dx()*maxLength)+borderWidth*7, textSize.Dy()+borderWidth*7)
	outerBounds := outerImage.Bounds()
	innerRectangle := image.Rectangle{
		image.Point{
			X: outerBounds.Min.X + borderWidth,
			Y: outerBounds.Min.Y + borderWidth,
		},
		image.Point{
			X: outerBounds.Max.X - borderWidth,
			Y: outerBounds.Max.Y - borderWidth,
		}}
	outerImage.Fill(color.White)
	innerImage := outerImage.SubImage(innerRectangle).(*ebiten.Image)
	innerImage.Fill(color.Black)
	return &TextInput{
		label:       label,
		borderWidth: borderWidth,
		maxLength:   maxLength,
		Clickable: Clickable{
			x:     x,
			y:     y,
			image: outerImage,
		},
	}
}

func (ti *TextInput) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		println("true")
		if ti.IsClicked() {
			println("activating")
			ti.isActive = true
			receivingKeyboardInput[ti.label] = true
		} else {
			println("deactivating")
			ti.isActive = false
			receivingKeyboardInput[ti.label] = false
			ti.counter = 0
		}
	}
	if ti.isActive {
		ti.counter++
		ti.runes = ebiten.AppendInputChars(ti.runes[:0])
		if len(ti.data) < ti.maxLength {
			ti.data += string(ti.runes)
		}
		println(ti.data)
	}
}

func (ti *TextInput) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		labelBounds := text.BoundString(NormalFace, ti.label)
		cursorX, cursorY := ebiten.CursorPosition()
		// Puike
		minX, maxX := ti.x-ti.image.Bounds().Dx()/3, float64(ti.x)+float64(ti.image.Bounds().Dx())/1.5
		minY, maxY := ti.y+labelBounds.Dy(), float64(ti.y)+float64(ti.image.Bounds().Dy())*1.5
		if int(minX) < cursorX && cursorX < int(maxX) {
			if int(minY) < cursorY && cursorY < int(maxY) {
				return true
			}
		}
	}

	return false
}

func (ti *TextInput) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	labelBounds := text.BoundString(NormalFace, ti.label)
	op := ebiten.DrawImageOptions{}
	// barf
	op.GeoM.Translate(float64(ti.x-ti.image.Bounds().Dx()/3), float64(ti.y+labelBounds.Dy()))
	screen.DrawImage(ti.image, &op)
	DrawStaticText(
		ti.label,
		NormalFace,
		ti.x,
		ti.y,
		color.White,
		screen,
		true,
	)
	t := ti.data
	if ti.counter%60 > 30 {
		t += "|"
	}
	// Need to position this
	DrawStaticText(
		ti.data,
		NormalFace,
		ti.x-ti.image.Bounds().Dx()/3+ti.borderWidth*2,
		ti.y+labelBounds.Dy()+ti.borderWidth,
		color.White,
		screen,
		false,
	)
}

func (ti *TextInput) GetInput() string {
	return ti.data
}
