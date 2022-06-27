package world

import (
	"encoding/json"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
)

// WorldMode represents the type for representing the current game mode
type WorldMode interface {
	Type() net.TypedMessageType
	Init(w *World) error
	Update(w *World) (WorldMode, error)
	Draw(w *World, screen *ebiten.Image)
	String() string
	Local() bool
}

// PreGame leads to Build mode.
type PreGameMode struct {
	local bool
}

func (m PreGameMode) String() string {
	return "pre"
}
func (m PreGameMode) Type() net.TypedMessageType {
	return 500
}
func (m *PreGameMode) Init(w *World) error {
	for _, pl := range w.Game.Players() {
		pl.ReadyForWave = false
	}

	return nil
}
func (m *PreGameMode) Update(w *World) (next WorldMode, err error) {
	// Just immediately go to build mode.
	next = &BuildMode{local: true}
	return
}
func (m *PreGameMode) Draw(w *World, screen *ebiten.Image) {
}
func (m *PreGameMode) Local() bool {
	return m.local
}

// BuildMode leads to Wave mode.
type BuildMode struct {
	local bool
}

func (m BuildMode) String() string {
	return "build"
}
func (m BuildMode) Type() net.TypedMessageType {
	return 501
}
func (m *BuildMode) Init(w *World) error {
	w.CurrentWave++
	data.BGM.Set("build-phase.ogg")
	return nil
}
func (m *BuildMode) Update(w *World) (next WorldMode, err error) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		w.Game.Players()[0].ReadyForWave = true
		if w.Game.Net().Active() {
			w.Game.Net().SendReliable(StartModeRequest{})
		}
	}

	if w.ArePlayersReady() {
		next = &WaveMode{local: true}
	}
	return
}
func (m *BuildMode) Draw(w *World, screen *ebiten.Image) {
	// First draw/get pathing overlay for spawners.
	for _, e := range w.spawners {
		lastX, lastY := e.physics.X, e.physics.Y
		for _, s := range e.steps {
			x := float64(s.X()*data.CellWidth + data.CellWidth/2)
			y := float64(s.Y()*data.CellHeight + data.CellHeight/2)
			c := data.GetPolarityColor(e.physics.polarity)
			c.A = 128
			ebitenutil.DrawLine(screen, w.CameraX+lastX, w.CameraY+lastY, w.CameraX+float64(x), w.CameraY+float64(y), c)
			lastX = float64(x)
			lastY = float64(y)
		}
	}

	bounds := text.BoundString(data.NormalFace, "build mode")
	text.Draw(
		screen,
		"build mode",
		data.NormalFace,
		8,
		bounds.Dy()+24,
		color.White,
	)
	// Hmm.
	msg := "hit <spacebar> to start combat waves"
	bounds = text.BoundString(data.NormalFace, msg)
	text.Draw(
		screen,
		msg,
		data.NormalFace,
		ScreenWidth/2-bounds.Dx()/2,
		ScreenHeight-50,
		color.White,
	)
}
func (m *BuildMode) Local() bool {
	return m.local
}

// WaveMode leads to Wave, Loss, Victory, or PostGame.
type WaveMode struct {
	local bool
}

func (m WaveMode) String() string {
	return "wave"
}
func (m WaveMode) Type() net.TypedMessageType {
	return 502
}
func (m *WaveMode) Init(w *World) error {
	// Play that funky music.
	data.BGM.Set("wave-phase.ogg")
	// Reset the player's ready state.
	for _, pl := range w.Game.Players() {
		pl.ReadyForWave = false
	}
	// Unleash the spawners.
	for _, s := range w.spawners {
		s.heldWave = false
	}

	return nil
}
func (m *WaveMode) Update(w *World) (next WorldMode, err error) {
	if w.AreSpawnersHolding() && w.AreEnemiesDead() {
		next = &BuildMode{local: true}
	} else if w.AreCoresDead() {
		next = &LossMode{local: true}
	} else if w.AreWavesComplete() {
		next = &VictoryMode{local: true}
		// Move onto build mode...?
	}
	return
}
func (m *WaveMode) Draw(w *World, screen *ebiten.Image) {
}
func (m *WaveMode) Local() bool {
	return m.local
}

// LossMode represents when the core 'splodes. Leads to a restart of the current.
type LossMode struct {
	local bool
}

func (m LossMode) String() string {
	return "loss"
}
func (m LossMode) Type() net.TypedMessageType {
	return 503
}
func (m *LossMode) Init(w *World) error {
	return nil
}
func (m *LossMode) Update(w *World) (next WorldMode, err error) {
	return
}
func (m *LossMode) Draw(w *World, screen *ebiten.Image) {
}
func (m *LossMode) Local() bool {
	return m.local
}

// VictoryMode represents when all waves are finished. Leads to Travel state.
type VictoryMode struct {
	local bool
}

func (m VictoryMode) String() string {
	return "victory"
}
func (m VictoryMode) Type() net.TypedMessageType {
	return 504
}
func (m *VictoryMode) Init(w *World) error {
	return nil
}
func (m *VictoryMode) Update(w *World) (next WorldMode, err error) {
	return
}
func (m *VictoryMode) Draw(w *World, screen *ebiten.Image) {
}
func (m *VictoryMode) Local() bool {
	return m.local
}

// PostGameMode is... the final victory...?
type PostGameMode struct {
	local bool
}

func (m PostGameMode) String() string {
	return "post"
}
func (m PostGameMode) Type() net.TypedMessageType {
	return 505
}
func (m *PostGameMode) Init(w *World) error {
	return nil
}
func (m *PostGameMode) Update(w *World) (next WorldMode, err error) {
	return
}
func (m *PostGameMode) Draw(w *World, screen *ebiten.Image) {
}
func (m *PostGameMode) Local() bool {
	return m.local
}

type StartModeRequest struct {
}

func (r StartModeRequest) Type() net.TypedMessageType {
	return 510
}

func init() {
	net.AddTypedMessage(500, func(data json.RawMessage) net.Message {
		var m PreGameMode
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(501, func(data json.RawMessage) net.Message {
		var m BuildMode
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(502, func(data json.RawMessage) net.Message {
		var m WaveMode
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(503, func(data json.RawMessage) net.Message {
		var m LossMode
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(504, func(data json.RawMessage) net.Message {
		var m VictoryMode
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(505, func(data json.RawMessage) net.Message {
		var m PostGameMode
		json.Unmarshal(data, &m)
		return m
	})

	net.AddTypedMessage(510, func(data json.RawMessage) net.Message {
		var m StartModeRequest
		json.Unmarshal(data, &m)
		return m
	})

}
