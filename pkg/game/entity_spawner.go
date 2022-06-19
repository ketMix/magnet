package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpawnerEntity struct {
	BaseEntity
	floatTick int
	active    bool
	// spawnTargets []EnemyKind ???
}

func NewSpawnerEntity() *SpawnerEntity {
	return &SpawnerEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
		floatTick: rand.Intn(60), // Lightly randomize that start.
	}
}

func (e *SpawnerEntity) Update(world *World) (request Request, err error) {
	// TODO: after some duration, attempt a spawn. The request should be handled such that it uses pathfinding to find the best spot towards the player's core.
	e.floatTick++
	return request, nil
}

func (e *SpawnerEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	op.GeoM.Translate(
		-float64(spawnerImage.Bounds().Dx())/2,
		0,
	)

	// Draw shadow
	{
		sop := &ebiten.DrawImageOptions{}
		sop.GeoM.Concat(op.GeoM)
		sop.GeoM.Translate(
			float64(spawnerShadowImage.Bounds().Dx())/2,
			0,
		)
		screen.DrawImage(spawnerShadowImage, sop)
	}

	// Draw from center.
	op.GeoM.Translate(
		0,
		-float64(spawnerImage.Bounds().Dy())/2-math.Sin(float64(e.floatTick)/30)*2,
	)

	// TODO: Make an "active" mode that has an alternative image or an image underlay.
	screen.DrawImage(spawnerImage, op)
}
