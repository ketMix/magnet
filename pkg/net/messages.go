package net

import "encoding/json"

// HandshakeMessage represents the type for the handshake step of networking.
type HandshakeMessage int

// These are our base types for handshaking.
const (
	RegisterMessage HandshakeMessage = iota
	AwaitMessage
	ArrivedMessage
	HelloMessage
)

// TypedMessageType represents the contained type within a TypedMessage.
type TypedMessageType int

// These are our typed message types.
const (
	MissingMessageType TypedMessageType = iota
	HenloMessageType
)

// TypedMessage wraps a Message.
type TypedMessage struct {
	Type TypedMessageType `json:"t"`
	Data json.RawMessage  `json:"d"`
}

// Message returns the typed message's wrapped data as a Message.
func (t *TypedMessage) Message() Message {
	switch t.Type {
	case HenloMessageType:
		var m HenloMessage
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
}

// Type returns HenloMessage's corresponding type number.
func (m HenloMessage) Type() TypedMessageType {
	return HenloMessageType
}
