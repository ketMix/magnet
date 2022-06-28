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
	innerImage  *ebiten.Image
	OnChange    func(s string)
}

func NewTextInput(label, placeholder string, maxLength, x, y int) *TextInput {
	borderWidth := 2
	textSize := text.BoundString(NormalFace, "m")
	outerImage := ebiten.NewImage((textSize.Dx()*(maxLength))+borderWidth*7, textSize.Dy()+borderWidth*7)
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
		data:        placeholder,
		innerImage:  innerImage,
		Clickable: Clickable{
			x:     x - outerBounds.Dx()/2,
			y:     y - outerBounds.Dy()/2,
			image: outerImage,
		},
	}
}

func (ti *TextInput) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if ti.IsClicked() {
			println("Clicked ", ti.label)
			ti.isActive = true
			receivingKeyboardInput[ti.label] = true
		} else {
			ti.isActive = false
			receivingKeyboardInput[ti.label] = false
			ti.counter = 0
		}
	}
	if ti.isActive {
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(ti.data) > 0 {
				ti.data = ti.data[0 : len(ti.data)-1]
				if ti.OnChange != nil {
					ti.OnChange(ti.data)
				}
			}
		}
		ti.counter++
		ti.runes = ebiten.AppendInputChars(ti.runes[:0])
		if len(ti.data) < ti.maxLength {
			ti.data += string(ti.runes)
			if len(ti.runes) > 0 && ti.OnChange != nil {
				ti.OnChange(ti.data)
			}
		}
	}
}

func (ti *TextInput) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// labelBounds := text.BoundString(NormalFace, ti.label)
		cursorX, cursorY := ebiten.CursorPosition()
		width := ti.innerImage.Bounds().Dx()
		height := ti.innerImage.Bounds().Dy() / 2
		labelBounds := text.BoundString(NormalFace, ti.label)
		minX, maxX := ti.x, ti.x+width
		minY, maxY := ti.y+labelBounds.Dy()-height, ti.y+labelBounds.Dy()+height
		if int(minX) < cursorX && cursorX < int(maxX) {
			if int(minY) < cursorY && cursorY < int(maxY) {
				return true
			}
		}
	}

	return false
}

func (ti *TextInput) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(ti.x), float64(ti.y))

	// Draw the input box
	screen.DrawImage(ti.image, &op)
	DrawStaticText(
		ti.label,
		NormalFace,
		ti.x,
		ti.y-ti.borderWidth,
		color.White,
		screen,
		false,
	)

	// Add a lil _ if we can type
	t := ti.data
	if ti.counter%60 > 30 {
		t += "_"
	}

	// Draw the input
	DrawStaticText(
		t,
		NormalFace,
		ti.x+ti.borderWidth*2,
		ti.y+ti.borderWidth*7,
		color.White,
		screen,
		false,
	)
}

func (ti *TextInput) GetInput() string {
	return ti.data
}
