package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/data/assets/lang"
	"github.com/kettek/ebijam22/pkg/world"
)

type HelpOverlay struct {
}

func (o *HelpOverlay) Draw(screen *ebiten.Image) {
	infoColor := color.RGBA{64, 96, 200, 255}

	// Draw tool usage information.
	x := 8
	y := world.ScreenHeight - world.ScreenHeight/5
	data.DrawStaticTextByCode(lang.HelpToolsTurrets, data.NormalFace, x, y, infoColor, screen, false)
	y += 12
	data.DrawStaticTextByCode(lang.HelpCost, data.NormalFace, x, y, infoColor, screen, false)

	// Draw wave and mode information.
	x = 32
	y = world.ScreenHeight / 4
	data.DrawStaticTextByCode(lang.HelpWaves, data.NormalFace, x, y, infoColor, screen, false)
	y += 12
	data.DrawStaticTextByCode(lang.HelpShown, data.NormalFace, x, y, infoColor, screen, false)

	// Draw orb information and player info
	x = world.ScreenWidth - 30
	y = world.ScreenHeight / 6
	textString := data.GiveMeString(lang.HelpPlayers)
	bounds := text.BoundString(data.NormalFace, textString)
	data.DrawStaticText(textString, data.NormalFace, x-bounds.Dx(), y, infoColor, screen, false)

	y += 12
	textString = data.GiveMeString(lang.HelpReady)
	bounds = text.BoundString(data.NormalFace, textString)
	data.DrawStaticText(textString, data.NormalFace, x-bounds.Dx(), y, infoColor, screen, false)
	data.DrawStaticText(textString, data.NormalFace, x-bounds.Dx(), y, infoColor, screen, false)

	// Draw controls
	x = world.ScreenWidth / 2
	y = world.ScreenHeight / 6
	data.DrawStaticTextByCode(lang.HelpControls, data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpMove, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpSprint, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpShoot, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpDeconstruct, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpInvert, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpSelect, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpShowRange, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpRestart, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpFullscreen, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpEscape, data.NormalFace, x, y, color.White, screen, true)
	y += 32

	data.DrawStaticTextByCode(lang.HelpObjectives, data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpBuildTurrets, data.NormalFace, x, y, color.White, screen, true)
	y += 16
	data.DrawStaticTextByCode(lang.HelpDefend, data.NormalFace, x, y, color.White, screen, true)
	y += 32

	data.DrawStaticTextByCode(lang.HelpToggleHelp, data.BoldFace, x, y, color.RGBA{255, 255, 0, 255}, screen, true)
}
