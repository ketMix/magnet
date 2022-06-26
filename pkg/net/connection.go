package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	OtherName    string

	//
	connected bool
	active    bool
	hosting   bool

	//
	lastReceived time.Time
	lastSent     time.Time

	//
	messages             chan Message
	reliableID           int
	sendingReliables     []ReliableTypedMessage
	receivedReliables    []ReliableTypedMessage
	confirmedReliableIDs []int
	reliableLock         sync.Mutex
}

func NewConnection(name string) Connection {
	return Connection{
		Name:     name,
		messages: make(chan Message, 1000),
		active:   true,
	}
}

// Connected returns if the connection is actually connected.
func (c *Connection) Connected() bool {
	return c.connected
}

// Active returns if the connection should be connected.
func (c *Connection) Active() bool {
	return c.active
}

// Hosting returns if the connection is acting as a host.
func (c *Connection) Hosting() bool {
	return c.hosting
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
	} else {
		c.hosting = true
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
			c.otherAddress = fromAddr
			return nil
		} else {
			fmt.Println("unhandled message from", fromAddr.String())
			continue
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
	} else {
		c.hosting = true
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
	if err := c.Send(HenloMessage{
		Name:     c.Name,
		Greeting: "hai",
	}); err != nil {
		panic(err)
	}
	c.connected = true
	c.lastReceived = time.Now()
	c.lastSent = time.Now()
	for {
		t := time.Now()
		// More than 5 seconds have passed since last receive, presume failure.
		if t.Sub(c.lastReceived) > 5*time.Second {
			fmt.Println("lost connection")
			c.connected = false
			return
		}
		// Send a ping every 3 seconds.
		if t.Sub(c.lastSent) > 3*time.Second {
			c.Send(PingMessage{})
		}

		// Attempt to read any pending messages, with a 2 second deadline.
		b := make([]byte, 1000)
		c.conn.SetReadDeadline(t.Add(2 * time.Second))
		n, foreignAddr, err := c.conn.ReadFromUDP(b)
		if err != nil && !os.IsTimeout(err) {
			fmt.Println("disconnect")
			c.connected = false
			fmt.Println(err)
			return
		}
		if foreignAddr.String() != c.otherAddress.String() {
			continue
		}
		b = b[:n]
		var msg ReliableTypedMessage
		if err = json.Unmarshal(b, &msg); err != nil {
			fmt.Println(err)
		} else {
			c.lastReceived = time.Now()
			// Handle reliable messages such that we send a response whenever we receive it.
			if msg.InboundID != 0 {
				c.reliableLock.Lock()
				// See if we've already confirmed sending this particular message.
				for _, id := range c.confirmedReliableIDs {
					if id == msg.InboundID {
						c.reliableLock.Unlock()
						// Ignore it if we've already confirmed it.
						goto loopEnd
					}
				}
				for i, m := range c.sendingReliables {
					if m.OutboundID == msg.InboundID {
						c.confirmedReliableIDs = append(c.confirmedReliableIDs, m.OutboundID)
						c.sendingReliables = append(c.sendingReliables[:i], c.sendingReliables[i+1:]...)
						c.reliableLock.Unlock()
						goto loopEnd
					}
				}
				c.reliableLock.Unlock()
			} else if msg.OutboundID != 0 {
				c.reliableLock.Lock()
				// Always send an empty response with the ID.
				c.sendReliableWithIDs(nil, msg.OutboundID, 0)
				// Only continue processing if we haven't processed this specific message yet.
				for _, m := range c.receivedReliables {
					if m.OutboundID == msg.InboundID {
						c.reliableLock.Unlock()
						goto loopEnd
					}
				}
				c.receivedReliables = append(c.receivedReliables, msg)
				c.reliableLock.Unlock()
			}
			switch m := msg.Message().(type) {
			case HenloMessage:
				c.OtherName = m.Name
			case PingMessage:
			default:
				c.messages <- m
			}
		}
	loopEnd:
		// This is a bit of a bad spot for this, but resend any unconfirmed reliable messages at this point.
		for _, m := range c.sendingReliables {
			n := time.Now()
			if n.Sub(m.lastSent) >= 500*time.Millisecond {
				c.reliableLock.Lock()
				c.sendReliableWithIDs(m.Message(), 0, m.OutboundID)
				c.reliableLock.Unlock()
				m.lastSent = n
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
	c.lastSent = time.Now()
	return err
}

// SendReliable sends the given message with special resending until a confirmation is received.
func (c *Connection) SendReliable(msg Message) error {
	c.reliableLock.Lock()
	c.reliableID++
	env, err := c.sendReliableWithIDs(msg, 0, c.reliableID)
	env.lastSent = time.Now()
	c.sendingReliables = append(c.sendingReliables, env)
	c.reliableLock.Unlock()
	return err
}

func (c *Connection) sendReliableWithIDs(msg Message, inboundID, outboundID int) (ReliableTypedMessage, error) {
	var envelope ReliableTypedMessage

	payload, err := json.Marshal(msg)
	if err != nil {
		return envelope, err
	}

	if msg != nil {
		envelope.Type = msg.Type()
	}
	envelope.Data = payload
	envelope.InboundID = inboundID
	envelope.OutboundID = outboundID

	bytes, err := json.Marshal(envelope)
	if err != nil {
		return envelope, err
	}

	if bytes != nil {
		_, err = c.conn.WriteTo(bytes, c.otherAddress)
	}
	c.lastSent = time.Now()
	return envelope, err
}

// Messages returns the current contents of the messages channel as a slice.
func (c *Connection) Messages() (m []Message) {
	for {
		done := false
		select {
		case msg := <-c.messages:
			m = append(m, msg)
		default:
			done = true
			break
		}
		if done {
			break
		}
	}
	return m
}
