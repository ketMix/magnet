package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
	"github.com/kettek/ebijam22/pkg/world"
)

type PlayState struct {
	game          *Game
	levelDataName string
	level         data.Level
	world         world.World
	messages      []Message
	clickables    []data.UIComponent
	framebuffer   *ebiten.Image
}

func (s *PlayState) Init() error {
	// Create our framebuffer so we can do some nicer fx.
	s.framebuffer = ebiten.NewImage(world.ScreenWidth, world.ScreenHeight)

	s.world.Game = s.game // Eww
	s.world.Speed = s.game.Options.Speed

	// Add players here...?
	s.game.players = append(s.game.players, world.NewPlayer())
	s.game.players[0].Local = true
	// Add other player!
	if s.game.net.Active() {
		s.game.players = append(s.game.players, world.NewPlayer())
		// Set player names if networked.
		s.game.players[0].Name = s.game.net.Name
		s.game.players[1].Name = s.game.net.OtherName
	}

	// Build the level.
	if err := s.world.BuildFromLevel(s.level); err != nil {
		return err
	}

	s.world.Mode = &world.PreGameMode{}
	s.clickables = []data.UIComponent{
		data.NewBGMIcon(),
		data.NewSFXIcon(),
	}
	return nil
}

func (s *PlayState) Dispose() error {
	// Remove player entity reference.
	for _, p := range s.game.players {
		p.Entity = nil
	}
	// Remove players. Should this be moved to a preplay state? Something between menu and travel for setting up players.
	s.game.players = make([]*world.Player, 0)

	// Dispose framebuffer. This is unnecessary afaik.
	if s.framebuffer != nil {
		s.framebuffer.Dispose()
	}
	return nil
}

func (s *PlayState) Update() error {
	// You may judge me for this, but I leave myself in the arms of the Gopher.
	if s.framebuffer.Bounds().Dx() != world.ScreenWidth || s.framebuffer.Bounds().Dy() != world.ScreenHeight {
		s.framebuffer.Dispose()
		s.framebuffer = ebiten.NewImage(world.ScreenWidth, world.ScreenHeight)
	}

	// Update the clickables
	for _, c := range s.clickables {
		c.Update()
	}

	// World mode handling.
	switch s.world.Mode.(type) {
	case *world.LossMode:
		// TODO: Show hit "R" to restart or something. Also maybe stats.
		// if s.game.net.Hosting() || !s.game.net.Active() {
		// 	s.game.SetState(&TravelState{
		// 		game:        s.game,
		// 		targetLevel: s.levelDataName,
		// 	})
		// }
	case *world.VictoryMode:
		// TODO: Show end game stats, if possible! Then some sort of "hit okay" to travel button/key.
		if s.level.Next != "" {
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				if s.game.net.Hosting() || !s.game.net.Active() {
					fmt.Println("TRAVELING TO NEXT")
					s.game.SetState(&TravelState{
						game:        s.game,
						targetLevel: s.level.Next,
					})
				}
			}
		} else {
			// We beat the video game.
		}
	}

	// If we're the host/solo and we hit R, restart the level. If we're the client, send a request.
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: s.levelDataName, // ???
			})
		} else {
			s.game.net.Send(net.TravelMessage{})
		}
		return nil
	}

	// Handle our network updates.
	for _, msg := range s.game.net.Messages() {
		switch msg := msg.(type) {
		case net.TravelMessage:
			if !s.game.net.Hosting() {
				s.game.SetState(&TravelState{
					game:        s.game,
					targetLevel: msg.Destination,
					fromLive:    true,
				})
			} else {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to restart! Hit 'r' to conform.", s.game.net.OtherName),
				})
			}
		case world.StartModeRequest:
			s.game.players[1].ReadyForWave = true
			if !s.world.ArePlayersReady() {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to start! Hit '<spacebar>' to conform.", s.game.net.OtherName),
				})
			}
		default:
			// Send unhandled messages to the world.
			s.world.ProcessNetMessage(msg)
		}
	}

	// Process our local messages.
	t := s.messages[:0]
	for _, m := range s.messages {
		m.lifetime++
		if m.lifetime < m.deathtime {
			t = append(t, m)
		}
	}
	s.messages = t

	// Update our players.
	for _, p := range s.game.players {
		action, err := p.Update(&s.world)
		if err != nil {
			return err
		}
		if action != nil {
			// If our net is active, send our desired action to the other.
			if s.game.net.Active() {
				s.game.net.Send(action)
			}
			p.Entity.SetAction(action)
		}
	}

	// Update our world.
	if err := s.world.Update(); err != nil {
		return err
	}

	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	// Draw our world.
	s.world.Draw(s.framebuffer)

	// Draw the framebuffer.
	framebufferOp := ebiten.DrawImageOptions{}
	screen.DrawImage(s.framebuffer, &framebufferOp)

	// Draw level text centered at top of screen for now.
	data.DrawStaticText(
		s.level.Title,
		data.BoldFace,
		world.ScreenWidth/2,
		5,
		color.White,
		screen,
		true,
	)

	// Draw our messages from most recent to oldest, bottom to top.
	mx := 8
	my := world.ScreenHeight - 80
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		bounds := text.BoundString(data.NormalFace, m.content)

		d := float64(m.lifetime) / float64(m.deathtime)
		c := color.RGBA{
			255,
			255,
			255,
			255 - uint8(255*d),
		}

		text.Draw(
			screen,
			m.content,
			data.NormalFace,
			mx,
			my,
			c,
		)
		my -= bounds.Dy() + 2
	}

	// Draw mode.
	s.world.Mode.Draw(&s.world, screen)

	// Draw the waves and current points.
	mx = 8
	my = 16
	offset := 16
	t := fmt.Sprintf("wave: %d/%d", s.world.CurrentWave, s.world.MaxWave)
	bounds := text.BoundString(data.NormalFace, t)
	data.DrawStaticText(
		t,
		data.NormalFace,
		mx,
		my,
		color.White,
		screen,
		false,
	)

	mx = bounds.Dx() + offset

	// Draw players and points. (don't judge me)
	for i, pl := range s.game.players {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(world.ScreenWidth)-12, float64(i)*32)

		// Draw our player image and name from right to left.
		imgs := data.Player2Init.Images
		if !s.game.net.Active() || s.game.net.Hosting() {
			if i == 0 {
				imgs = data.PlayerInit.Images
			}
		} else {
			if i == 1 {
				imgs = data.PlayerInit.Images
			}
		}
		op.GeoM.Translate(-float64(imgs[0].Bounds().Dx()/2), float64(imgs[0].Bounds().Dy()))
		screen.DrawImage(imgs[0], op)

		bounds := text.BoundString(data.BoldFace, pl.Name)
		text.Draw(screen, pl.Name, data.BoldFace, int(op.GeoM.Element(0, 2))-bounds.Dx()-4, int(op.GeoM.Element(1, 2))+8, color.White)

		// Move down and draw our points.
		orb, _ := data.GetImage("orb-large.png")
		op.GeoM.Translate(1, 16)

		screen.DrawImage(orb, op)

		op.GeoM.Translate(-float64(orb.Bounds().Dx()), 0)

		t := fmt.Sprint(pl.Points)
		bounds = text.BoundString(data.NormalFace, t)
		op.GeoM.Translate(-float64(bounds.Dx()), 0)
		text.Draw(screen, t, data.NormalFace, int(op.GeoM.Element(0, 2)), int(op.GeoM.Element(1, 2))+8, color.White)
	}
	// Draw current points
	/*t = fmt.Sprint(s.world.Points)
	bounds = text.BoundString(data.NormalFace, t)
	data.DrawStaticText(
		fmt.Sprint(s.world.Points),
		data.NormalFace,
		mx,
		my,
		color.White,
		screen,
		false,
	)
	mx += bounds.Dx() + 5
	orb, _ := data.GetImage("orb-large.png")
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mx), float64(my-orb.Bounds().Dy()))
	screen.DrawImage(orb, &op)

	mx += orb.Bounds().Dx() + offset*/

	// Draw our clickables
	if s.clickables != nil {
		for i, c := range s.clickables {
			c.SetPos(mx+(i+1)*offset, 12)
			c.Draw(screen, &ebiten.DrawImageOptions{})
		}
	}

	// Draw our player's belt!
	s.game.players[0].Toolbelt.Draw(screen)
}

func (s *PlayState) AddMessage(m Message) {
	if m.deathtime <= 0 || m.deathtime >= 1000 {
		m.deathtime = 300
	}
	s.messages = append(s.messages, m)
}
