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
	var addr string
	var name string
	var target string

	addr = os.Args[1]
	name = os.Args[2]

	if len(os.Args) > 3 {
		target = os.Args[3]
	}

	c := net.NewConnection(name)

	fmt.Println(addr, name, target)

	c.Await("gamu.group:20220", addr, target)
}

func help() {
	fmt.Printf("Syntax: %s <address:port> <name> [<target>]\n", os.Args[0])
}
