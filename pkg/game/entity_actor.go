package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ActorEntity struct {
	BaseEntity
	player *Player
	speed  float64
}

func NewActorEntity(player *Player, config EntityConfig) *ActorEntity {
	return &ActorEntity{
		speed: config.speed,
		BaseEntity: BaseEntity{
			animation: Animation{
				images: config.images,
			},
			health: config.health,
			physics: PhysicsObject{
				polarity: config.polarity,
			},
			turret: Turret{
				damage: config.damage,
				speed:  config.projecticleSpeed,
				rate:   config.attackRate,
			},
		},
		player: player,
	}
}

func (e *ActorEntity) Update(world *World) (request Request, err error) {
	switch a := e.action.(type) {
	case *EntityActionMove:
		if math.Abs(e.physics.X-a.x) < a.distance && math.Abs(e.physics.Y-a.y) < a.distance {
			a.complete = true
			break
		}

		// FIXME: Make this use actual physics resolution!
		r := math.Atan2(e.physics.Y-a.y, e.physics.X-a.x)
		x := math.Cos(r) * e.speed
		y := math.Sin(r) * e.speed

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
		image := e.animation.Image()
		// Get our player position for spawning.
		px := e.Physics().X
		py := e.Physics().Y - float64(image.Bounds().Dy())/2

		// Get direction vector from difference of player and target.
		vX, vY := GetDirection(px, py, float64(a.targetX), float64(a.targetY))
		xSide := 1.0
		if vX < 0 {
			xSide = -xSide
		}

		// Can apply player's speed to action vector
		a.complete = true
		request = SpawnProjecticleRequest{
			x:        px + (float64(image.Bounds().Dx()/2) * xSide),
			y:        py,
			vX:       vX * e.turret.speed,
			vY:       vY * e.turret.speed,
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
	image := e.animation.Image()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	// Draw from center.
	// FIXME: We should probably use an explicit "originX" and "originY" variables.
	op.GeoM.Translate(
		-float64(image.Bounds().Dx())/2,
		// Adjust Y to render from the "foot" of the image
		-float64(image.Bounds().Dy()),
	)
	screen.DrawImage(image, op)
	// NOTE: We _could_ draw something like a target marker for a moving action here.
}
