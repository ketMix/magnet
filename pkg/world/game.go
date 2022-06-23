package world

import "github.com/kettek/ebijam22/pkg/net"

type Game interface {
	Players() []*Player
	Net() *net.Connection
}
