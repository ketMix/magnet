package world

import (
	"encoding/json"

	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
)

// Request represents results from an entity's action completion.
type Request interface {
	Type() net.TypedMessageType
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
	X           float64
	Y           float64
	Polarity    data.Polarity `json:"p"`
	enemyConfig data.EntityConfig
	Kind        string `json:"k"`
}

// MultiRequest is a container for multiple requests.
type MultiRequest struct {
	Requests []Request `json:"r"`
}

// Belt-related requests.

// SelectToolbeltItemRequest selects a given toolbelt item
type SelectToolbeltItemRequest struct {
	kind ToolKind
}

// DummyRequest is used to prevent action passthrough.
type DummyRequest struct {
}

// Here be code for networking again.
func (r MultiRequest) Type() net.TypedMessageType {
	return 300
}

func (r SpawnEnemyRequest) Type() net.TypedMessageType {
	return 301
}

func (r UseToolRequest) Type() net.TypedMessageType {
	return net.MissingMessageType
}

func (r SpawnProjecticleRequest) Type() net.TypedMessageType {
	return net.MissingMessageType
}

func (r SelectToolbeltItemRequest) Type() net.TypedMessageType {
	return net.MissingMessageType
}

func (r DummyRequest) Type() net.TypedMessageType {
	return net.MissingMessageType
}

func init() {
	net.AddTypedMessage(300, func(data json.RawMessage) net.Message {
		var m MultiRequest
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(301, func(data json.RawMessage) net.Message {
		var m SpawnEnemyRequest
		json.Unmarshal(data, &m)
		return m
	})
}
