package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/world"
)

type HelpOverlay struct {
}

func (o *HelpOverlay) Draw(screen *ebiten.Image) {
	infoColor := color.RGBA{64, 96, 200, 255}

	// Draw tool usage information.
	x := 8
	y := world.ScreenHeight - world.ScreenHeight/5
	text.Draw(screen, "Tools and turrets are here.", data.NormalFace, x, y, infoColor)
	y += 12
	text.Draw(screen, "Constructing costs points.", data.NormalFace, x, y, infoColor)

	// Draw wave and mode information.
	x = 32
	y = world.ScreenHeight / 4
	text.Draw(screen, "Waves and upcoming enemies", data.NormalFace, x, y, infoColor)
	y += 12
	text.Draw(screen, "are shown here.", data.NormalFace, x, y, infoColor)

	// Draw orb information and player info
	x = world.ScreenWidth - 30
	y = world.ScreenHeight / 6
	bounds := text.BoundString(data.NormalFace, "Players, their points,")
	text.Draw(screen, "Players, their points,", data.NormalFace, x-bounds.Dx(), y, infoColor)
	bounds = text.BoundString(data.NormalFace, "and ready status are shown here.")
	y += 12
	text.Draw(screen, "and ready status are shown here.", data.NormalFace, x-bounds.Dx(), y, infoColor)

	// Draw controls
	x = world.ScreenWidth / 2
	y = world.ScreenHeight / 6
	data.DrawStaticText("Controls", data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
	y += 16
	data.DrawStaticText("WASD : Move", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Shift : Sprint", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Left Mouse : Shoot/Construct", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Right Mouse : Move/Deconstruct", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Middle Mouse / Tab : Invert Tool/Turret Polarity", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Mousewheel / 1-9: Select Tool/Turret", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Alt : Show Turret Range", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("R : Restart", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("F : Fullscreen", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Escape : Escape Menu", data.NormalFace, x, y, color.White, screen, true)
	y += 32

	data.DrawStaticText("Objectives", data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
	y += 16
	data.DrawStaticText("Build turrets opposite the portals' polarities!", data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticText("Defend the crystal at all costs!", data.NormalFace, x, y, color.White, screen, true)
	y += 32

	data.DrawStaticText("Press H or F1 to toggle this screen!", data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
}
