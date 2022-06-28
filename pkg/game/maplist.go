package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/world"
)

type MapList struct {
	buttons     []*data.Button
	selectedMap string
}

func (m *MapList) Init() error {
	mapFiles, err := data.GetPathFiles("levels")
	if err != nil {
		return err
	}
	x := 48
	y := 0
	for _, f := range mapFiles {
		if f == "README" {
			continue
		}
		if m.selectedMap == "" {
			m.selectedMap = f
		}
		(func(f string) {
			b := data.NewButton(x, y, f, func() {
				m.selectedMap = f
			})
			bounds := text.BoundString(data.NormalFace, f)
			m.buttons = append(m.buttons, b)
			x += bounds.Dx() * 2
			if x >= world.ScreenWidth-100 {
				y += bounds.Dy() + 8
				x = 8
			}
		})(f)
	}
	return nil
}

func (m *MapList) Update() error {
	for _, b := range m.buttons {
		b.Update()
		if m.selectedMap == b.Text() {
			b.Active = true
		} else {
			b.Active = false
		}
	}
	return nil
}

func (m *MapList) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	data.DrawStaticText("Map: ", data.BoldFace, int(op.GeoM.Element(0, 2)), int(op.GeoM.Element(1, 2))+4, color.White, screen, false)
	for _, b := range m.buttons {
		b.Draw(screen, op)
	}
}
