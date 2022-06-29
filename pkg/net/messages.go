package net

import (
	"encoding/json"
	"time"
)

// HandshakeMessage represents the type for the handshake step of networking.
type HandshakeMessage int

// These are our base types for handshaking.
const (
	RegisterMessage HandshakeMessage = iota
	AwaitMessage
	ArrivedMessage
	HelloMessage
	HandshakerMessage
	BroadcastMessage
)

// TypedMessageType represents the contained type within a TypedMessage.
type TypedMessageType int

// These are our typed message types.
const (
	MissingMessageType TypedMessageType = iota
	HenloMessageType
	PingMessageType
	TravelMessageType
	MoveMessageType
)

var topType TypedMessageType = 100

type ReliableTypedMessage struct {
	TypedMessage
	InboundID  int `json:"i"`
	OutboundID int `json:"o"`
}

// TypedMessage wraps a Message.
type TypedMessage struct {
	Type     TypedMessageType `json:"t"`
	Data     json.RawMessage  `json:"d"`
	lastSent time.Time        // Locally used value to determine re-send time.d time.
}

// TypedMessageHandler represents a dynamically defined unmarshaller for a message.
type TypedMessageHandler struct {
	Unwrap func(data json.RawMessage) Message
}

// TypedMessageMap is a map of all defined types->TypedMessageHandlers
var TypedMessageMap = make(map[TypedMessageType]TypedMessageHandler)

// AddTypedMessage adds the given unwrapper func as a Message unmarshaler. If typeIndex is 0, then an automatic index is generated and returned.
func AddTypedMessage(typeIndex int, unwrapper func(data json.RawMessage) Message) TypedMessageType {
	if typeIndex == 0 {
		topType++
		typeIndex = int(topType)
	}
	TypedMessageMap[TypedMessageType(typeIndex)] = TypedMessageHandler{
		Unwrap: unwrapper,
	}
	return TypedMessageType(typeIndex)
}

// Message returns the typed message's wrapped data as a Message.
func (t *TypedMessage) Message() Message {
	if tt, ok := TypedMessageMap[t.Type]; ok {
		return tt.Unwrap(t.Data)
	}
	switch t.Type {
	case HenloMessageType:
		var m HenloMessage
		json.Unmarshal(t.Data, &m)
		return m
	case PingMessageType:
		var m PingMessage
		json.Unmarshal(t.Data, &m)
		return m
	case TravelMessageType:
		var m TravelMessage
		json.Unmarshal(t.Data, &m)
		return m
	}
	return nil
}

// Message represents a message that can be sent as a typed message's data.
type Message interface {
	Type() TypedMessageType
}

// HenloMessage is our basic greeting message.
type HenloMessage struct {
	Greeting string `json:"g"`
	Name     string `json:"n"`
}

// Type returns HenloMessage's corresponding type number.
func (m HenloMessage) Type() TypedMessageType {
	return HenloMessageType
}

// PingMessage is used to send periodic pings.
type PingMessage struct {
}

// Type returns PingMessage's corresponding type number.
func (m PingMessage) Type() TypedMessageType {
	return PingMessageType
}

// TravelMessage is sent by the host to clients to enforce travel.
type TravelMessage struct {
	Destination string `json:"d"`
}

// Type returns TravelMessage's corresponding type number.
func (m TravelMessage) Type() TypedMessageType {
	return TravelMessageType
}
