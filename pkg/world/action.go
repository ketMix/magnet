package world

import (
	"encoding/json"

	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
)

// EntityAction is the interface to represent an entity's desired action.
type EntityAction interface {
	// Replaceble returns if the action can be immediately replaced by another.
	Replaceable() bool
	// Complete returns if the action is completed. If it is, then it is removed by the entity.
	Complete() bool
	GetNext() EntityAction
	Type() net.TypedMessageType
}

// EntityActionMove represents an action that should move the entity towards a target location.
type EntityActionMove struct {
	// x and y represents the target position to move to.
	X float64 `json:"x"`
	Y float64 `json:"y"`
	// distance represents the distance from the target that should be considered valid.
	Distance float64 `json:"d"`
	// relative represents if the movement is considered as relative to the entity's current position. Should this even be a thing?
	relative bool
	//
	complete bool
	//
	Next EntityAction `json:"n"`
}

func (a *EntityActionMove) Replaceable() bool {
	return true
}

func (a *EntityActionMove) Complete() bool {
	return a.complete
}

func (a *EntityActionMove) GetNext() EntityAction {
	return a.Next
}

type EntityActionPlace struct {
	// x and y are cell positions to place at.
	X        int
	Y        int
	complete bool
	Kind     ToolKind      `json:"k"`
	Polarity data.Polarity `json:"p"`
}

func (a *EntityActionPlace) Replaceable() bool {
	return true
}

func (a *EntityActionPlace) Complete() bool {
	return a.complete
}

func (a *EntityActionPlace) GetNext() EntityAction {
	return nil
}

type EntityActionShoot struct {
	// targetX and targetY represent the position to fire at.
	TargetX  float64       `json:"x"`
	TargetY  float64       `json:"y"`
	Polarity data.Polarity `json:"p"`
	complete bool
	Next     EntityAction `json:"n"`
}

func (a *EntityActionShoot) Replaceable() bool {
	return false
}

func (a *EntityActionShoot) Complete() bool {
	return a.complete
}

func (a *EntityActionShoot) GetNext() EntityAction {
	return a.Next
}

// Here be code to add network support to above actions.
func (a EntityActionMove) Type() net.TypedMessageType {
	return 200
}

func (a EntityActionPlace) Type() net.TypedMessageType {
	return 201
}

func (a EntityActionShoot) Type() net.TypedMessageType {
	return 202
}

func init() {
	net.AddTypedMessage(200, func(data json.RawMessage) net.Message {
		var m EntityActionMove
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(202, func(data json.RawMessage) net.Message {
		var m EntityActionShoot
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(203, func(data json.RawMessage) net.Message {
		var m EntityActionPlace
		json.Unmarshal(data, &m)
		return m
	})
}
