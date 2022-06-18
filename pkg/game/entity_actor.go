package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ActorEntity struct {
	BaseEntity
	player *Player
}

func NewActorEntity(player *Player) *ActorEntity {
	return &ActorEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
		player: player,
	}
}

func (e *ActorEntity) Update() (request Request, err error) {
	speed := 1.0
	switch a := e.action.(type) {
	case *EntityActionMove:
		if math.Abs(e.physics.X-a.x) < a.distance && math.Abs(e.physics.Y-a.y) < a.distance {
			a.complete = true
			break
		}

		// FIXME: Make this use actual physics resolution!
		r := math.Atan2(e.physics.Y-a.y, e.physics.X-a.x)
		x := math.Cos(r) * speed
		y := math.Sin(r) * speed

		e.physics.X -= x
		e.physics.Y -= y
	case *EntityActionPlace:
		a.complete = true
		request = UseToolRequest{
			x:    a.x,
			y:    a.y,
			kind: a.kind,
		}
	case *EntityActionShoot:
		// Can apply player's speed to action vector
		a.complete = true
		request = SpawnProjecticleRequest{
			x:        a.x,
			y:        a.y,
			vX:       a.vX,
			vY:       a.vY,
			polarity: a.polarity,
		}
	}

	// Separate action removal for now.
	if e.action != nil && e.action.Complete() {
		e.action = e.action.Next()
	}
	return request, nil
}

func (e *ActorEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	// Draw from center.
	// FIXME: We should probably use an explicit "originX" and "originY" variables.
	op.GeoM.Translate(
		-float64(playerImage.Bounds().Dx())/2,
		// Adjust Y to render from the "foot" of the image
		-float64(playerImage.Bounds().Dy()),
	)
	screen.DrawImage(playerImage, op)
	// NOTE: We _could_ draw something like a target marker for a moving action here.
}
