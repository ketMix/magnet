package world

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type ActorEntity struct {
	BaseEntity
	player        *Player
	speed         float64
	walking       bool
	right         bool
	walkAnimation Animation
	idleAnimation Animation
}

func NewActorEntity(player *Player, config data.EntityConfig) *ActorEntity {
	fmt.Println(len(config.WalkImages))
	return &ActorEntity{
		speed: config.Speed,
		walkAnimation: Animation{
			images:    config.WalkImages,
			frameTime: 5,
			speed:     1,
		},
		idleAnimation: Animation{
			images: config.Images,
		},
		BaseEntity: BaseEntity{
			animation: Animation{
				images: config.Images,
			},
			health: config.Health,
			physics: PhysicsObject{
				polarity: config.Polarity,
			},
			turret: Turret{
				damage:         config.Damage,
				speed:          config.ProjecticleSpeed,
				rate:           config.AttackRate,
				projecticleNum: config.ProjecticleNum,
			},
		},
		player: player,
	}
}

func (e *ActorEntity) Update(world *World) (requests MultiRequest, err error) {
	e.animation.Update()
	switch a := e.action.(type) {
	case *EntityActionMove:
		if math.Abs(e.physics.X-a.X) < a.Distance && math.Abs(e.physics.Y-a.Y) < a.Distance {
			a.complete = true
			break
		}

		// FIXME: Make this use actual physics resolution!
		r := math.Atan2(e.physics.Y-a.Y, e.physics.X-a.X)
		x := math.Cos(r) * e.speed
		y := math.Sin(r) * e.speed

		e.physics.X -= x
		e.physics.Y -= y

		if x < 0 {
			e.right = true
		} else {
			e.right = false
		}
	case *EntityActionPlace:
		a.complete = true
		requests.Requests = append(requests.Requests, UseToolRequest{
			x:          a.X,
			y:          a.Y,
			tool:       a.Tool,
			toolConfig: data.TurretConfigs[a.Kind],
			polarity:   a.Polarity,
		})
	case *EntityActionShoot:
		image := e.animation.Image()
		// Get our player position for spawning.
		px := e.Physics().X
		py := e.Physics().Y - float64(image.Bounds().Dy())/2

		// Get direction vector from difference of player and target.
		vX, vY := GetDirection(px, py, float64(a.TargetX), float64(a.TargetY))

		// Can apply player's speed to action vector
		a.complete = true

		const spreadArc = 45.0
		var vectors = SplitVectorByDegree(spreadArc, vX, vY, e.turret.projecticleNum)
		for _, v := range vectors {
			projecticle := &ProjecticleEntity{
				BaseEntity: BaseEntity{
					physics: PhysicsObject{
						vX:       v.vX * e.turret.speed,
						vY:       v.vY * e.turret.speed,
						polarity: e.player.Toolbelt.activeItem.polarity,
					},
				},
				lifetime: 500,
				damage:   e.turret.damage,
			}
			request := SpawnProjecticleRequest{
				x:          px,
				y:          py,
				projectile: projecticle,
			}
			requests.Requests = append(requests.Requests, request)
		}
	}

	// Separate action removal for now.
	if e.action != nil && e.action.Complete() {
		e.action = e.action.GetNext()
	}

	// Set our animation to idle if nothing else is doin.
	if e.action == nil {
		if e.walking {
			e.walking = false
			e.animation = e.idleAnimation
		}
	} else {
		if !e.walking {
			e.walking = true
			e.animation = e.walkAnimation
		}
	}

	return requests, nil
}

func (e *ActorEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	image := e.animation.Image()
	op := &ebiten.DrawImageOptions{}

	// Draw from center.
	// FIXME: We should probably use an explicit "originX" and "originY" variables.
	op.GeoM.Translate(
		-float64(image.Bounds().Dx())/2,
		// Adjust Y to render from the "foot" of the image
		-float64(image.Bounds().Dy()),
	)

	// Mirror image if we're moving right
	if e.right {
		op.GeoM.Scale(-1, 1)
	}

	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	screen.DrawImage(image, op)
	// NOTE: We _could_ draw something like a target marker for a moving action here.
}
