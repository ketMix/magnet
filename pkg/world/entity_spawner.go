package world

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type SpawnerEntity struct {
	BaseEntity
	floatTick   int
	active      bool
	shouldSpawn bool
	// spawnTargets []EnemyKind ???
}

func NewSpawnerEntity(p data.Polarity) *SpawnerEntity {
	return &SpawnerEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{
				polarity: p,
			},
		},
		floatTick:   rand.Intn(60), // Lightly randomize that start.
		active:      true,
		shouldSpawn: true,
	}
}

func (e *SpawnerEntity) Update(world *World) (request Request, err error) {
	// TODO: after some duration, attempt a spawn. The request should be handled such that it uses pathfinding to find the best spot towards the player's core.
	e.floatTick++
	n := math.Sin(float64(e.floatTick)/30) * 2
	max := 1.0
	if n < -max && e.shouldSpawn {
		e.shouldSpawn = false
		var enemyConfig data.EntityConfig
		switch e.physics.polarity {
		case data.PositivePolarity:
			enemyConfig = data.EnemyConfigs["walker-positive"]
		case data.NegativePolarity:
			enemyConfig = data.EnemyConfigs["walker-negative"]
		}
		request = SpawnEnemyRequest{
			x:           e.physics.X,
			y:           e.physics.Y,
			enemyConfig: enemyConfig,
		}
	} else if n > max {
		e.shouldSpawn = true
	}
	return request, nil
}

func (e *SpawnerEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	var img *ebiten.Image
	if e.physics.polarity == data.NegativePolarity {
		img, _ = data.GetImage("spawner-negative.png")
	} else if e.physics.polarity == data.PositivePolarity {
		img, _ = data.GetImage("spawner-positive.png")
	} else {
		img, _ = data.GetImage("spawner.png")
	}

	op.GeoM.Translate(
		-float64(img.Bounds().Dx())/2,
		0,
	)

	// Draw shadow
	{
		shadowImg, _ := data.GetImage("spawner-shadow.png")
		sop := &ebiten.DrawImageOptions{}
		sop.GeoM.Concat(op.GeoM)
		sop.GeoM.Translate(
			float64(shadowImg.Bounds().Dx())/2,
			0,
		)
		screen.DrawImage(shadowImg, sop)
	}

	// Draw from center.
	op.GeoM.Translate(
		0,
		-float64(img.Bounds().Dy())/2-math.Sin(float64(e.floatTick)/30)*2,
	)

	if e.active {
		portalImg, _ := data.GetImage("spawner-portal.png")
		screen.DrawImage(portalImg, op)
	}

	screen.DrawImage(img, op)
}
