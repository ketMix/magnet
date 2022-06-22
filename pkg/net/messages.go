package net

type HandshakeMessage int

const (
	RegisterMessage HandshakeMessage = iota
	AwaitMessage
	ArrivedMessage
	HelloMessage
)

func init() {
}
