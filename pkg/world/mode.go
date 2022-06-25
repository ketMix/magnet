package world

import (
	"encoding/json"

	"github.com/kettek/ebijam22/pkg/net"
)

// WorldMode represents the type for representing the current game mode
type WorldMode int

// These are our various states.
const (
	PreGameMode  WorldMode = iota // PreGame leads to Build mode.
	BuildMode                     // Build leads to Wave mode.
	WaveMode                      // Wave leads to Wave, Loss, Victory, or PostGame.
	LossMode                      // Loss represents when the core 'splodes. Leads to a restart of the current.
	VictoryMode                   // Victory represents when all waves are finished. Leads to Travel state.
	PostGameMode                  // PostGame is... the final victory...?
)

func (m WorldMode) String() string {
	if m == PreGameMode {
		return "pre"
	} else if m == BuildMode {
		return "build"
	} else if m == WaveMode {
		return "wave"
	} else if m == LossMode {
		return "loss"
	} else if m == VictoryMode {
		return "victory"
	} else if m == PostGameMode {
		return "post"
	}
	return "dunno"
}

// SetModeRequest is used to both set the mode in world, as well as the data type send to the client to let them know we're moving along.
type SetModeRequest struct {
	Mode  WorldMode `json:"m"`
	local bool      // If it is a locally generated request.
}

func (r SetModeRequest) Type() net.TypedMessageType {
	return 500
}

type StartModeRequest struct {
}

func (r StartModeRequest) Type() net.TypedMessageType {
	return 501
}

func init() {
	net.AddTypedMessage(500, func(data json.RawMessage) net.Message {
		var m SetModeRequest
		json.Unmarshal(data, &m)
		return m
	})
	net.AddTypedMessage(501, func(data json.RawMessage) net.Message {
		var m StartModeRequest
		json.Unmarshal(data, &m)
		return m
	})

}
