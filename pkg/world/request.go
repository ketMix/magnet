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
	X, Y     int
	Tool     ToolKind      `json:"t"`
	Kind     string        `json:"k"`
	Polarity data.Polarity `json:"p"`
	NetID    int           `json:"i"` // Yeah, yeah, we shouldn't have NetID here, but it's easier to reuse UseToolRequest rather than implement some new SpawnTurret/SpawnWall/RemoveWall Request set.
	Owner    string        `json:"o"` // The owner's name. This is a little excessive to send, but it's easier than mucking about with client/server index checking. Also enables more players if we ever want that.
	local    bool          // Used to determine if the result of this tool use should be considered the server's or the client's.
}

// SpawnProjecticleRequest attempts to spawn a projecticle at given location with given direction
type SpawnProjecticleRequest struct {
	X, Y       float64 // Position
	projectile *ProjecticleEntity
	NetID      int `json:"i"`
	VX, VY     float64
	Polarity   data.Polarity `json:"p"`
	Damage     int           `json:"d"`
}

type SpawnEnemyRequest struct {
	X           float64
	Y           float64
	Polarity    data.Polarity `json:"p"`
	enemyConfig data.EntityConfig
	Kind        string `json:"k"`
	NetID       int    `json:"i"`
}

// SpawnToolEntityRequest is used to tell the client to spawn an entity tied to a tool.
type SpawnToolEntityRequest struct {
	X, Y     int
	Tool     ToolKind      `json:"t"`
	Kind     string        `json:"k"`
	Polarity data.Polarity `json:"p"`
	NetID    int           `json:"i"`
}

// TrashEntityRequest is send from the server to client(s) to let them know to delete the given entity.
type TrashEntityRequest struct {
	NetID  int    `json:"i"`
	entity Entity // Used locally to just trash.
	local  bool   // Used to determine in the trash request is local.
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

func (r SpawnProjecticleRequest) Type() net.TypedMessageType {
	return 302
}

func (r UseToolRequest) Type() net.TypedMessageType {
	return 303
}

func (r SpawnToolEntityRequest) Type() net.TypedMessageType {
	return 304
}

func (r SelectToolbeltItemRequest) Type() net.TypedMessageType {
	return net.MissingMessageType
}

func (r TrashEntityRequest) Type() net.TypedMessageType {
	return 310
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
	net.AddTypedMessage(302, func(data json.RawMessage) net.Message {
		var m SpawnProjecticleRequest
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(303, func(data json.RawMessage) net.Message {
		var m UseToolRequest
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(304, func(data json.RawMessage) net.Message {
		var m SpawnToolEntityRequest
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(310, func(data json.RawMessage) net.Message {
		var m TrashEntityRequest
		json.Unmarshal(data, &m)
		return m
	})
}
