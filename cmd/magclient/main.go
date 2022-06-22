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

	var c net.Connection

	if os.Args[1] == "host" {
		addr = os.Args[2]
		name = os.Args[3]

		c = net.NewConnection(name)

		go func() {
			err := c.AwaitDirect(addr, target)
			if err != nil {
				panic(err)
			}
		}()
	} else if os.Args[1] == "join" {
		target = os.Args[2]
		name = os.Args[3]

		c = net.NewConnection(name)

		go func() {
			err := c.AwaitDirect(addr, target)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		handshaker = os.Args[1]
		addr = os.Args[2]
		name = os.Args[3]

		if len(os.Args) > 4 {
			target = os.Args[4]
		}

		c = net.NewConnection(name)

		fmt.Println(addr, name, target)

		go func() {
			err := c.AwaitHandshake(handshaker, addr, target)
			if err != nil {
				panic(err)
			}
		}()
	}
	for {
		select {
		case m := <-c.Messages:
			fmt.Println("got message", m)
		}
	}
}

func help() {
	fmt.Printf("Syntax: %s <signaler> <address:port> <name> [<target>]\n", os.Args[0])
	fmt.Printf("Syntax: %s join <target address:port> <name>\n", os.Args[0])
	fmt.Printf("Syntax: %s host <address:port> <name>\n", os.Args[0])
}
