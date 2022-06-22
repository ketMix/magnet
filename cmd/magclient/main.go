package main

import (
	"fmt"
	"os"

	"github.com/kettek/ebijam22/pkg/net"
)

func main() {
	if len(os.Args) < 3 {
		help()
		os.Exit(1)
	}
	var handshaker string
	var addr string
	var name string
	var target string

	handshaker = os.Args[1]
	addr = os.Args[2]
	name = os.Args[3]

	if len(os.Args) > 4 {
		target = os.Args[4]
	}

	c := net.NewConnection(name)

	fmt.Println(addr, name, target)

	go func() {
		c.Await(handshaker, addr, target)
	}()
	for {
		select {
		case m := <-c.Messages:
			fmt.Println("got message", m)
		}
	}
}

func help() {
	fmt.Printf("Syntax: %s <address:port> <name> [<target>]\n", os.Args[0])
}
