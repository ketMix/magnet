/*
This file provides a _very_ simple handshaker service for use with basic UDP punching.
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	enet "github.com/kettek/ebijam22/pkg/net"
)

// AddressKey is a key that represents an ip address and port.
type AddressKey string

type MessageBox struct {
	name     string              // String for this messagebox
	wavingAt map[string]struct{} // waving at other names
}

var clientsMap map[AddressKey]*MessageBox = make(map[AddressKey]*MessageBox)

func IPToAddressKey(addr *net.UDPAddr) (a AddressKey) {
	return AddressKey(addr.String())
}

func AddressKeyToIP(a AddressKey) *net.UDPAddr {
	addr, _ := net.ResolveUDPAddr("udp", string(a))
	return addr
}

func main() {
	if len(os.Args) == 1 {
		os.Exit(1)
	}
	address := os.Args[1]
	fmt.Println("Starting handshaker...", address)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	localConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Begin the Eternal Listen (tm)
	for {
		buffer := make([]byte, 1024)
		bytesRead, remoteAddr, err := localConn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		clientKey := IPToAddressKey(remoteAddr)

		msg := string(buffer[0:bytesRead])
		parts := strings.Split(msg, " ")
		a, err := strconv.Atoi(parts[0])

		fmt.Println("[INCOMING]", msg)
		//if incoming.
		if enet.HandshakeMessage(a) == enet.RegisterMessage {
			if _, ok := clientsMap[clientKey]; !ok {
				clientsMap[clientKey] = new(MessageBox)
				clientsMap[clientKey].wavingAt = make(map[string]struct{})
				log.Printf("Registered %s as %s\n", clientKey, parts[1])
			}
			clientsMap[clientKey].name = parts[1]
			// Check if any clients are waiting this target and send arrival msg.
			for otherClientKey, mbox := range clientsMap {
				if _, ok := mbox.wavingAt[parts[1]]; ok {
					delete(mbox.wavingAt, parts[1])
					sendArrival(localConn, clientKey, otherClientKey)
					sendArrival(localConn, otherClientKey, clientKey)
				}
			}
		} else if enet.HandshakeMessage(a) == enet.AwaitMessage {
			fmt.Println("got someone seeking", clientKey, parts[1])
			mbox, ok := clientsMap[clientKey]
			if !ok {
				continue
			}
			var matched = false
			for otherClientKey, otherMbox := range clientsMap {
				if otherMbox.name == parts[1] {
					sendArrival(localConn, otherClientKey, clientKey)
					sendArrival(localConn, clientKey, otherClientKey)
					matched = true
				}
			}
			if !matched {
				if _, ok := mbox.wavingAt[parts[1]]; !ok {
					mbox.wavingAt[parts[1]] = struct{}{}
					log.Printf("Awaiting arrival of %s for %s\n", parts[1], mbox.name)
				}
			}
		}
	}
}

func sendArrival(conn *net.UDPConn, to, target AddressKey) {
	log.Printf("Sending arrival of %s to %s\n", target, to)
	toAddress := AddressKeyToIP(to)
	conn.WriteTo([]byte(fmt.Sprintf("%d %s", enet.ArrivedMessage, target)), toAddress)
}
