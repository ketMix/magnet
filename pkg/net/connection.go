package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Name string

	// handshakerAddr is the target handshaker service.
	handshakerAddr *net.UDPAddr

	// conn is our own base connection.
	conn *net.UDPConn

	// otherConn is the peer we wish to play with.
	otherConn    *net.UDPConn
	otherAddress *net.UDPAddr

	//
	connected bool

	//
	Messages chan Message
}

func NewConnection(name string) Connection {
	return Connection{
		Name:     name,
		Messages: make(chan Message, 1000),
	}
}

// Connected returns if the connection is actually connected.
func (c *Connection) Connected() bool {
	return c.connected
}

func (c *Connection) AwaitHandshake(handshaker string, local string, target string) error {
	handshakerAddr, err := net.ResolveUDPAddr("udp", handshaker)
	if err != nil {
		return err
	}

	// Get a random local port.
	localAddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}
	log.Printf("Attempting to listen on %s\n", localAddr.String())

	// Start listening!
	localConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return err
	}

	c.handshakerAddr = handshakerAddr
	c.conn = localConn
	fmt.Println("listening on", localConn.LocalAddr().String())

	log.Println("Sending register message to handshaker service")
	_, err = localConn.WriteTo([]byte(fmt.Sprintf("%d %s", RegisterMessage, c.Name)), c.handshakerAddr)
	if err != nil {
		return err
	}

	// Wait for a character string response.
	for {
		buffer := make([]byte, 2)
		c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		bytesRead, fromAddr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		// Ignore sends from non-handhsaker.
		if fromAddr.String() != handshakerAddr.String() {
			continue
		}
		c.conn.SetReadDeadline(time.Time{})
		msg := string(buffer[0:bytesRead])
		a, err := strconv.Atoi(msg)
		if err != nil {
			return err
		}
		if a != int(HandshakerMessage) {
			return errors.New("incorrect handshake response")
		}
		break
	}

	if target != "" {
		log.Printf("Sending await message for %s to handshaker service\n", target)
		_, err := localConn.WriteTo([]byte(fmt.Sprintf("%d %s", AwaitMessage, target)), c.handshakerAddr)
		if err != nil {
			return err
		}
	}

	return c.awaitHandshake()
}

func (c *Connection) awaitHandshake() error {
	fmt.Println("entering main await")
	for {
		buffer := make([]byte, 1024)
		bytesRead, fromAddr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		msg := string(buffer[0:bytesRead])
		parts := strings.Split(msg, " ")
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		if a == int(ArrivedMessage) {
			otherAddr, err := net.ResolveUDPAddr("udp", parts[1])
			if err != nil {
				return err
			}
			_, err = c.conn.WriteTo([]byte(fmt.Sprintf("%d %s", HelloMessage, c.Name)), otherAddr)
			if err != nil {
				return err
			}
			c.otherAddress = otherAddr
			return nil
		} else if a == int(HelloMessage) {
			fmt.Println("got hello from self-declared", parts[1])
			fmt.Println(fromAddr.String())
			c.otherAddress = fromAddr
			return nil
		} else {
			return nil
			// BOGUS
		}
	}
}

// AwaitDirect attempts to set up a client-to-client connection without any handshaking.
func (c *Connection) AwaitDirect(local string, target string) error {
	localAddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}
	log.Printf("Attempting to listen on %s\n", localAddr.String())

	// Start listening!
	localConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return err
	}

	c.conn = localConn
	fmt.Println("listening on", localConn.LocalAddr().String())

	if target != "" {
		otherAddr, err := net.ResolveUDPAddr("udp", target)
		if err != nil {
			return err
		}
		_, err = c.conn.WriteTo([]byte(fmt.Sprintf("%d %s", HelloMessage, c.Name)), otherAddr)
		if err != nil {
			return err
		}
	}

	// Start the listen loop.
	for {
		buffer := make([]byte, 1024)
		bytesRead, fromAddr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		msg := string(buffer[0:bytesRead])
		parts := strings.Split(msg, " ")
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		if a == int(HelloMessage) {
			fmt.Println("got hello from self-declared", parts[1])
			fmt.Println(fromAddr.String())

			// Send hello back to let the other client that we're ready to rumble.
			_, err := c.conn.WriteTo([]byte(fmt.Sprintf("%d %s", HelloMessage, c.Name)), fromAddr)
			if err != nil {
				return err
			}

			c.otherAddress = fromAddr
			break
		} else {
			// BOGUS
		}
	}
	return nil
}

func (c *Connection) Loop() {
	fmt.Println("starting main loop with", c.otherAddress.String())
	if err := c.Send(HenloMessage{"hai from " + c.Name}); err != nil {
		panic(err)
	}
	c.connected = true
	for {
		var msg TypedMessage
		b := make([]byte, 10000)
		//c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, foreignAddr, err := c.conn.ReadFromUDP(b)
		if err != nil {
			c.connected = false
			fmt.Println(err)
			return
		}
		if foreignAddr.String() != c.otherAddress.String() {
			continue
		}
		b = b[:n]
		if err = json.Unmarshal(b, &msg); err != nil {
			fmt.Println(err)
		} else {
			m := msg.Message()
			if m != nil {
				c.Messages <- m
			}
		}
	}
}

// Send sends the given message interface to the other player.
func (c *Connection) Send(msg Message) error {
	var envelope TypedMessage

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	envelope.Type = msg.Type()
	envelope.Data = payload

	bytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	if bytes != nil {
		_, err = c.conn.WriteTo(bytes, c.otherAddress)
	}
	return err
}
