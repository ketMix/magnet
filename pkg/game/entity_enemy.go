package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/goro/pathing"
)

type EnemyEntity struct {
	BaseEntity
	path  pathing.Path
	steps []pathing.Step
}

func NewEnemyEntity(config EntityConfig) *EnemyEntity {
	return &EnemyEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    config.images,
				frameTime: 60,
				speed:     0.25,
			},
			health: config.health,
			physics: PhysicsObject{
				polarity:       config.polarity,
				magnetic:       config.magnetic,
				magnetStrength: config.magnetStrength,
				magnetRadius:   config.magnetRadius,
			},
		},
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	// Update animation.
	e.animation.Update()

	// Attempt to move along path to player's core
	if e.path != nil && len(e.steps) == 0 {
		cx, cy := world.GetClosestCellPosition(int(e.physics.X), int(e.physics.Y))
		e.steps = e.path.Compute(cx, cy, world.coreX, world.coreY)
	} else {
		// TODO: move towards step[0], then remove it when near its center. If the last one is to be removed, then we have reached the core.
	}

	return request, nil
}

func (e *EnemyEntity) CanPathfind() bool {
	return true
}

func (e *EnemyEntity) SetPath(p pathing.Path) {
	e.path = p
}

func (e *EnemyEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	// Draw animation.
	e.animation.Draw(screen, op)
}
