package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/game"
	"github.com/thought-machine/go-flags"
)

func main() {
	g := &game.Game{}

	// Parse our initial command-line derived options.
	if _, err := flags.Parse(&g.Options); err != nil {
		return
	}

	if err := os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl"); err != nil {
		fmt.Println("WARNING: OpenGL backend could not be set, expect degraded performance if on DirectX.")
	}

	if err := g.Init(); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		if _, ok := err.(game.NoError); !ok {
			panic(err)
		}
	}
}
