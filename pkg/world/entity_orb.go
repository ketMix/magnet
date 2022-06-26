package world

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type OrbEntity struct {
	BaseEntity
	elapsed  int
	worth    int
	lifetime int
}

func NewOrbEntity(worth int) *OrbEntity {
	// Get our animation images.
	var images []*ebiten.Image
	if worth <= 10 {
		images = data.OrbSmallImages
	} else if worth <= 15 {
		images = data.OrbMediumImages
	} else {
		images = data.OrbLargeImages
	}

	return &OrbEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    images,
				frameTime: 20,
				speed:     1,
			},
		},
		worth:    worth,
		lifetime: 2000,
	}
}

// The orb ponders
func (e *OrbEntity) Update(world *World) (request Request, err error) {
	e.elapsed++

	if e.elapsed >= e.lifetime {
		request = TrashEntityRequest{
			NetID:  e.netID,
			entity: e,
			local:  true,
		}
		return
	}

	// Update animation.
	e.animation.Update()

	// Meander towards players.
	for _, pl := range world.Game.Players() {
		if pl.Entity == nil {
			continue
		}

		// Check if we actually have intersected with a player.

		vX, vY := GetDistanceVector(e.physics.X, e.physics.Y, pl.Entity.Physics().X, pl.Entity.Physics().Y)
		if math.Abs(vX) < 10 && math.Abs(vY) < 10 {
			var r MultiRequest
			r.Requests = append(r.Requests, TrashEntityRequest{
				NetID:  e.netID,
				entity: e,
				local:  true,
			})
			r.Requests = append(r.Requests, CollectOrbRequest{
				Worth:     e.worth,
				Collector: pl.Name,
			})
			return r, nil
		}
		magnitude := GetMagnitude(vX, vY)
		vX, vY = Normalize(vX, vY, math.Max(magnitude, 1))
		affect := 1 / (magnitude * magnitude)

		e.physics.vX += vX * affect * world.Speed
		e.physics.vY += vY * affect * world.Speed
	}
	e.physics.X += e.physics.vX
	e.physics.Y += e.physics.vY

	return
}

func (e *OrbEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)

	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	if e.lifetime-e.elapsed < 500 {
		op.ColorM.Scale(1, 1, 1, float64(e.lifetime-e.elapsed)/500)
	}

	// Draw animation.
	e.animation.Draw(screen, op)
}
