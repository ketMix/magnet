package net

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	Name    string
	Encoder *gob.Encoder
	Decoder *gob.Decoder

	// handshakerAddr is the target handshaker service.
	handshakerAddr *net.UDPAddr

	// conn is our own base connection.
	conn *net.UDPConn

	// otherConn is the peer we wish to play with.
	otherConn    *net.UDPConn
	otherAddress *net.UDPAddr

	//
	Messages chan Message
}

func NewConnection(name string) Connection {
	return Connection{
		Name:     name,
		Messages: make(chan Message, 1000),
	}
}

func (c *Connection) Await(handshaker string, local string, target string) {
	handshakerAddr, _ := net.ResolveUDPAddr("udp", handshaker)

	// Get a random local port.
	localAddr, _ := net.ResolveUDPAddr("udp", local)
	log.Printf("Attempting to listen on %s\n", localAddr.String())

	// Start listening!
	localConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal(err)
	}

	c.handshakerAddr = handshakerAddr
	c.conn = localConn
	fmt.Println("listening on", localConn.LocalAddr().String())

	log.Println("Sending register message to handshaker service")
	_, err = localConn.WriteTo([]byte(fmt.Sprintf("%d %s", RegisterMessage, c.Name)), c.handshakerAddr)
	if err != nil {
		panic(err)
	}

	if target != "" {
		log.Printf("Sending await message for %s to handshaker service\n", target)
		_, err := localConn.WriteTo([]byte(fmt.Sprintf("%d %s", AwaitMessage, target)), c.handshakerAddr)
		if err != nil {
			panic(err)
		}
	}

	c.await()
}

func (c *Connection) await() {
	fmt.Println("entering main await")
	for {
		buffer := make([]byte, 1024)
		bytesRead, fromAddr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		msg := string(buffer[0:bytesRead])
		parts := strings.Split(msg, " ")
		a, _ := strconv.Atoi(parts[0])
		if a == int(ArrivedMessage) {
			otherAddr, _ := net.ResolveUDPAddr("udp", parts[1])
			c.conn.WriteTo([]byte(fmt.Sprintf("%d %s", HelloMessage, c.Name)), otherAddr)
			c.loop(otherAddr)
			return
		} else if a == int(HelloMessage) {
			fmt.Println("got hello from self-declared", parts[1])
			fmt.Println(fromAddr.String())
			c.loop(fromAddr)
			return
		} else {
			// BOGUS
		}
	}
}

func (c *Connection) loop(otherAddress *net.UDPAddr) {
	fmt.Println("starting main loop with", otherAddress.String())
	gob.Register(HenloMessage{})
	c.otherAddress = otherAddress
	c.Encoder = gob.NewEncoder(c)
	c.Decoder = gob.NewDecoder(c)
	if err := c.Send(HenloMessage{"hai from " + c.Name}); err != nil {
		panic(err)
	}
	for {
		var msg Message
		if err := c.Decoder.Decode(&msg); err != nil {
			panic(err)
		}
		switch msg := msg.(type) {
		case HenloMessage:
			fmt.Println("got henlo", msg)
			c.Messages <- msg
		default:
			fmt.Println("got unhandled message!", msg)
		}
	}
}

func (c *Connection) Read(p []byte) (n int, err error) {
	n, foreignAddr, err := c.conn.ReadFromUDP(p)
	p = p[:n]
	if foreignAddr.String() != c.otherAddress.String() {
		p = p[:0]
		return 0, nil
	}
	return
}

func (c *Connection) Write(p []byte) (n int, err error) {
	n, err = c.conn.WriteTo(p, c.otherAddress)
	return n, err
}

func (c *Connection) Send(msg Message) error {
	fmt.Println("sending", msg)
	return c.Encoder.Encode(&msg)
}
