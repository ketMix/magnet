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
	shouldSpawn bool
	wave        *data.Wave
	// spawnTargets []EnemyKind ???
	spawnElapsed int
}

func NewSpawnerEntity(p data.Polarity) *SpawnerEntity {
	return &SpawnerEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{
				polarity: p,
			},
		},
		floatTick:   rand.Intn(60), // Lightly randomize that start.
		shouldSpawn: true,
	}
}

func (e *SpawnerEntity) Update(world *World) (request Request, err error) {
	if e.wave != nil {
		if e.wave.Spawns != nil {
			if e.spawnElapsed >= e.wave.Spawns.Spawnrate {
				var spawnRequests MultiRequest
				// Spread out our spawns if more than 1 kind is to spawn.
				for i, k := range e.wave.Spawns.Kinds {
					// FIXME: Spread from center.
					spread := float64(i) / float64(len(e.wave.Spawns.Kinds))
					spreadX := spread * float64(data.CellWidth/2)
					spreadY := spread * float64(data.CellHeight/2)

					enemyConfig := data.EnemyConfigs[k]
					spawnRequests.requests = append(spawnRequests.requests,
						SpawnEnemyRequest{
							x:           e.physics.X + spreadX,
							y:           e.physics.Y + spreadY,
							enemyConfig: enemyConfig,
						},
					)
				}
				request = spawnRequests
				e.wave.Spawns.Count--
				if e.wave.Spawns.Count <= 0 {
					e.wave.Spawns = e.wave.Spawns.Next
				}
				e.spawnElapsed = 0
			}
		} else {
			e.wave = e.wave.Next
			e.spawnElapsed = 0
		}
	}
	e.spawnElapsed++

	e.floatTick++
	/*n := math.Sin(float64(e.floatTick)/30) * 2
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
	}*/
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

	if e.wave != nil {
		portalImg, _ := data.GetImage("spawner-portal.png")
		screen.DrawImage(portalImg, op)
	}

	screen.DrawImage(img, op)
}
