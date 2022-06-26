package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebijam22/pkg/data"
)

// Toolbelt is the interface for containing user actions for placing turrets and similar.
type Toolbelt struct {
	items      []*ToolbeltItem
	activeItem *ToolbeltItem
}

// Update updates the toolbelt. This seems a bit silly, but oh well.
func (t *Toolbelt) Update() (request Request) {
	// This is a stupid check.
	if t.activeItem == nil && len(t.items) > 0 {
		t.activeItem = t.items[0]
		t.activeItem.active = true
	}

	// Update our individual slots.
	for _, item := range t.items {
		r := item.Update()
		if r != nil {
			switch r.(type) {
			case SelectToolbeltItemRequest:
				// If we're selecting the same item twice, cycle the toolbelt item
				if t.activeItem == item {
					t.activeItem.Cycle()
				} else {
					if t.activeItem != nil {
						t.activeItem.active = false
					}
					t.activeItem = item
					t.activeItem.active = true
				}
			}
			request = r
			break
		}
	}
	return request
}

func (t *Toolbelt) CheckHit(x, y int) bool {
	return false
}

// Position positions the toolbelt and all its tools.
func (t *Toolbelt) Position() {
	toolSlotImage, _ := data.GetImage("toolslot.png")
	x, y := 8, ScreenHeight-8-toolSlotImage.Bounds().Dy()+toolSlotImage.Bounds().Dy()/2

	for _, ti := range t.items {
		ti.Position(&x, &y)
	}
}

func (t *Toolbelt) Draw(screen *ebiten.Image) {
	// Draw the belt slots.
	for _, ti := range t.items {
		ti.DrawSlot(screen)
	}
	// Then the slot items.
	for _, ti := range t.items {
		ti.Draw(screen)
	}
}

type ToolKind int

const (
	ToolNone ToolKind = iota
	ToolGun
	ToolTurret
	ToolWall
	ToolDestroy
)

// ToolbeltItem is a toolbelt entry.
type ToolbeltItem struct {
	tool     ToolKind
	kind     data.EntityConfig
	polarity data.Polarity
	x, y     int
	key      ebiten.Key // Key to check against for activation.
	active   bool
}

func (t *ToolbeltItem) Update() (request Request) {
	toolSlotImage, _ := data.GetImage("toolslot.png")
	// Does the cursor intersect us?
	if inpututil.IsKeyJustPressed(t.key) {
		return SelectToolbeltItemRequest{t.tool}
	} else {
		x, y := ebiten.CursorPosition()
		x1, x2 := t.x-toolSlotImage.Bounds().Dx()/2, t.x+toolSlotImage.Bounds().Dx()/2
		y1, y2 := t.y-toolSlotImage.Bounds().Dy()/2, t.y+toolSlotImage.Bounds().Dy()/2

		if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				return SelectToolbeltItemRequest{t.tool}
			}
			// Do a dummy return to prevent click through.
			return DummyRequest{}
		}
	}
	return nil
}

// Position assigns the center position for the toolbelt item.
func (t *ToolbeltItem) Position(sx, sy *int) {
	toolSlotImage, _ := data.GetImage("toolslot.png")
	toolDestroyImage, _ := data.GetImage("tool-destroy.png")
	t.x = *sx + toolDestroyImage.Bounds().Dx()/2
	t.y = *sy

	// Move forward our cursor.
	*sx += toolSlotImage.Bounds().Dx() + 1
}

func (t *ToolbeltItem) DrawSlot(screen *ebiten.Image) {
	toolSlotImage, _ := data.GetImage("toolslot.png")
	toolSlotActiveImage, _ := data.GetImage("toolslot-active.png")
	op := ebiten.DrawImageOptions{}
	if t.active {
		op.GeoM.Translate(float64(t.x-toolSlotActiveImage.Bounds().Dx()/2), float64(t.y-toolSlotActiveImage.Bounds().Dy()/2))
		screen.DrawImage(toolSlotActiveImage, &op)
	} else {
		op.GeoM.Translate(float64(t.x-toolSlotImage.Bounds().Dx()/2), float64(t.y-toolSlotImage.Bounds().Dy()/2))
		screen.DrawImage(toolSlotImage, &op)
	}
}

func (t *ToolbeltItem) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}

	// Move to the center of our item.
	op.GeoM.Translate(float64(t.x), float64(t.y))

	image := GetToolImage(t.tool, t.kind.Title)

	if image != nil {
		op.ColorM.Scale(data.GetPolarityColorScale(t.polarity))
		op.GeoM.Translate(-float64(image.Bounds().Dx()/2), -float64(image.Bounds().Dy()/2))
		screen.DrawImage(image, &op)
	}
}

// Cycles through available selections for the toolbelt item
func (t *ToolbeltItem) Cycle() {
	switch t.tool {
	// Abuse the fact that polarities have value
	case ToolGun:
		t.polarity++
		if t.polarity > data.PositivePolarity {
			t.polarity = data.NegativePolarity
		}
	case ToolTurret:
		t.polarity *= -1
	}
}

// Retrieves the image for a toolkind
// TODO: perhaps intialize toolbelt items with these instead?
func GetToolImage(t ToolKind, k string) *ebiten.Image {
	var image *ebiten.Image
	switch t {
	case ToolTurret:
		image = data.TurretConfigs[k].Images[0]
	case ToolDestroy:
		image, _ = data.GetImage("tool-destroy.png")
	case ToolGun:
		image, _ = data.GetImage("tool-gun.png")
	case ToolWall:
		image, _ = data.GetImage("wall.png")
	}
	return image
}
