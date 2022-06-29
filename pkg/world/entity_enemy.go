package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/goro/pathing"
)

type EnemyEntity struct {
	BaseEntity
	steps     []pathing.Step
	healthBar *ProgressBar
	speed     float64
	lifetime  float64
	// How many points are awarded upon defeat?
	points int
	//
	lastSync         int
	victoryAnimation Animation
	locked           bool // locked is used to lock the enemy entity when the mode changes to a loss.
	flies            bool
}

func NewEnemyEntity(config data.EntityConfig) *EnemyEntity {
	flies := false
	if config.Title == "flier" {
		flies = true
	}
	return &EnemyEntity{
		flies: flies,
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    config.WalkImages,
				frameTime: 30,
				speed:     1 - config.Speed,
			},
			health:    config.Health,
			maxHealth: config.Health,
			physics: PhysicsObject{
				polarity:       config.Polarity,
				magnetic:       config.Magnetic,
				magnetStrength: config.MagnetStrength,
				magnetRadius:   config.MagnetRadius,
				radius:         config.Radius,
			},
		},
		healthBar: NewProgressBar(
			7,
			1,
			color.RGBA{255, 0, 0, 1},
		),
		speed:  config.Speed,
		points: config.Points,
		victoryAnimation: Animation{
			images:    config.VictoryImages,
			frameTime: 30,
			speed:     1 - config.Speed,
		},
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	// Update animation.
	e.animation.Update()
	if e.locked {
		return
	}
	if e.health <= 0 {
		var requests MultiRequest
		requests.Requests = append(requests.Requests, TrashEntityRequest{
			NetID:  e.netID,
			entity: e,
			local:  true,
		})
		requests.Requests = append(requests.Requests, SpawnOrbRequest{
			X:     e.physics.X,
			Y:     e.physics.Y,
			Worth: e.points,
		})
		return requests, nil
	}

	e.lifetime++

	// Update healthbar
	e.healthBar.progress = float64(e.maxHealth) / float64(e.health)

	// Attempt to move along path to player's core
	if e.flies {
		// We receive steps, so let's just go to the last step (the core)
		if len(e.steps) != 0 {
			step := e.steps[len(e.steps)-1]
			tx := float64(step.X()*data.CellWidth + data.CellWidth/2)
			ty := float64(step.Y()*data.CellHeight + data.CellHeight/2)
			r := math.Atan2(e.physics.Y-ty, e.physics.X-tx)
			x := math.Cos(r) * e.speed * world.Speed
			y := math.Sin(r) * e.speed * world.Speed

			e.physics.X -= x
			e.physics.Y -= y

			if x > 0 {
				e.animation.mirror = false
			} else {
				e.animation.mirror = true
			}
			var requests MultiRequest
			for _, core := range world.cores {
				if e.IsCollided(core) {
					requests.Requests = append(requests.Requests, DamageCoreRequest{
						ID:     core.id,
						Damage: 1, // Should this be based upon some enemy damage value?
					})
				}
			}
			if len(requests.Requests) > 0 {
				requests.Requests = append(requests.Requests, TrashEntityRequest{
					NetID:  e.netID,
					entity: e,
					local:  true,
				})
				request = requests
			}
		}
	} else {
		if len(e.steps) != 0 {
			tx := float64(e.steps[0].X()*data.CellWidth + data.CellWidth/2)
			ty := float64(e.steps[0].Y()*data.CellHeight + data.CellHeight/2)
			if math.Abs(e.physics.X-float64(e.steps[0].X()*data.CellWidth+data.CellWidth/2)) < 1 && math.Abs(e.physics.Y-float64(e.steps[0].Y()*data.CellHeight+data.CellHeight/2)) < 1 {
				//
				e.steps = e.steps[1:]
			} else {
				r := math.Atan2(e.physics.Y-ty, e.physics.X-tx)
				x := math.Cos(r) * e.speed * world.Speed
				y := math.Sin(r) * e.speed * world.Speed

				e.physics.X -= x
				e.physics.Y -= y

				if x > 0 {
					e.animation.mirror = false
				} else {
					e.animation.mirror = true
				}
			}

			// TODO: move towards step[0], then remove it when near its center. If the last one is to be removed, then we have reached the core.
		} else {
			var requests MultiRequest
			for _, core := range world.cores {
				if e.IsCollided(core) {
					requests.Requests = append(requests.Requests, DamageCoreRequest{
						ID:     core.id,
						Damage: 1, // Should this be based upon some enemy damage value?
					})
				}
			}
			// No mo steppes
			requests.Requests = append(requests.Requests, TrashEntityRequest{
				NetID:  e.netID,
				entity: e,
				local:  true,
			})
			request = requests
			// We have reached a core, we should decrease the health on that core
			// Gotta figure out what core it reached...
		}
	}

	// Send periodic sync every 100 ticks. This is ignored during processing if the host is not set.
	e.lastSync++
	if e.lastSync > world.Game.GetOptions().SyncRate {
		e.lastSync = 0
		var r MultiRequest
		r.Requests = append(r.Requests, EntityPropertySync{
			X:      e.physics.X,
			Y:      e.physics.Y,
			NetID:  e.netID,
			Health: e.health,
		})
		if request != nil {
			r.Requests = append(r.Requests, request)
		}
		request = r
	}

	return request, nil
}

func (e *EnemyEntity) CanPathfind() bool {
	return true
}

func (e *EnemyEntity) SetSteps(s []pathing.Step) {
	e.steps = s
}

func (e *EnemyEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	if e.lifetime < 10 {
		op.ColorM.Scale(1, 1, 1, e.lifetime/10)
	}

	if e.physics.polarity != data.NeutralPolarity {
		r, g, b, a := data.GetPolarityColorScale(e.physics.polarity)
		op.ColorM.Scale(r*2, g*2, b*2, a)
	}

	// Draw animation.
	e.animation.Draw(screen, op)

	// Draw healthbar if less than max health
	if e.health < e.maxHealth {
		// Center the health bar horizontally and position it at the bottom of our image.
		// NOTE: we're using the first image, as if we use whatever the current frame is, it might be of different dimensions, which would lead to the health bar changing position.
		op.GeoM.Translate(-float64(e.animation.images[0].Bounds().Dx())/2, float64(e.animation.images[0].Bounds().Dy()))
		e.healthBar.Draw(screen, op)
	}
}
