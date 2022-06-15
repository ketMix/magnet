package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/game"
)

func main() {
	g := &game.Game{}

	if err := g.Init(); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
