package world

import "github.com/kettek/ebijam22/pkg/data"

// Request represents results from an entity's action completion.
type Request interface {
}

// UseToolRequest attempts to use the tool at a given cell.
type UseToolRequest struct {
	x, y     int
	kind     ToolKind // ???
	polarity data.Polarity
}

// SpawnProjecticleRequest attempts to spawn a projecticle at given location with given direction
type SpawnProjecticleRequest struct {
	x, y       float64 // Position
	projectile *ProjecticleEntity
}

type SpawnEnemyRequest struct {
	x, y        float64
	enemyConfig data.EntityConfig
}

// Belt-related requests.

// SelectToolbeltItemRequest selects a given toolbelt item
type SelectToolbeltItemRequest struct {
	kind ToolKind
}

// DummyRequest is used to prevent action passthrough.
type DummyRequest struct {
}
