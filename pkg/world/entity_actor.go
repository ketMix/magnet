package world

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type ActorEntity struct {
	BaseEntity
	player          *Player
	speed           float64
	walking         bool
	right           bool
	walkAnimation   Animation
	idleAnimation   Animation
	colorMultiplier [3]float64
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
		colorMultiplier: config.ColorMultiplier,
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

func (e *ActorEntity) Update(world *World) (request Request, err error) {
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

		targetX := e.physics.X - x
		targetY := e.physics.Y - y
		cellX := world.GetCell(world.GetClosestCellPosition(int(targetX), int(e.physics.Y)))
		cellY := world.GetCell(world.GetClosestCellPosition(int(e.physics.X), int(targetY)))
		if cellX != nil && cellX.kind != data.EmptyCell && cellX.kind != data.BlockedCell {
			e.physics.X = targetX
		}
		if cellY != nil && cellY.kind != data.EmptyCell && cellY.kind != data.BlockedCell {
			e.physics.Y = targetY
		}

		if x < 0 {
			e.right = true
		} else {
			e.right = false
		}
	case *EntityActionPlace:
		a.complete = true
		request = UseToolRequest{
			X:        a.X,
			Y:        a.Y,
			Tool:     a.Tool,
			Kind:     a.Kind,
			Polarity: a.Polarity,
			local:    true,
		}
	case *EntityActionShoot:
		image := e.animation.Image()

		// Get our player position for spawning.
		px := e.Physics().X
		py := e.Physics().Y - float64(image.Bounds().Dy())/2

		// Get direction vector from difference of player and target.
		vX, vY := GetDirection(px, py, float64(a.TargetX), float64(a.TargetY))

		a.complete = true

		if e.turret.projecticleNum == 1 {
			request = SpawnProjecticleRequest{
				X:        px,
				Y:        py,
				VX:       vX * e.turret.speed,
				VY:       vY * e.turret.speed,
				Polarity: a.Polarity,
				Damage:   e.turret.damage,
			}
		} else {
			var projecticleRequests MultiRequest
			const spreadArc = 45.0
			vectors := SplitVectorByDegree(spreadArc, vX, vY, e.turret.projecticleNum)
			for _, v := range vectors {
				req := SpawnProjecticleRequest{
					X:        px,
					Y:        py,
					VX:       v.vX * e.turret.speed,
					VY:       v.vY * e.turret.speed,
					Polarity: a.Polarity,
					Damage:   e.turret.damage,
				}
				projecticleRequests.Requests = append(projecticleRequests.Requests, req)
			}
			request = projecticleRequests
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

	return request, nil
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
