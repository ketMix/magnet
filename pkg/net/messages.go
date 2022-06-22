package net

type HandshakeMessage int

const (
	RegisterMessage HandshakeMessage = iota
	AwaitMessage
	ArrivedMessage
	HelloMessage
	PingMessage
)

type Message interface {
	Type() HandshakeMessage
}

type HenloMessage struct {
	Greeting string
}

func (m HenloMessage) Type() HandshakeMessage {
	return HelloMessage
}

func init() {
}
