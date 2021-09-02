package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Map of peers connected to this node (address -> peer)
var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn *websocket.Conn
}

// NON-MUTATING FUNCTIONS
// Initialize a new peer with the given connection, ip, port
func initPeer(conn *websocket.Conn, address, port string) *peer {
	newPeer := &peer{conn}
	nodeUrl := fmt.Sprintf("%s:%s", address, port)
	Peers[nodeUrl] = newPeer
	go read(newPeer) // listen to peer in goroutine
	return newPeer
}

// Continue to read and print from a peer (place in goroutine)
func read(p *peer) {
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", msg)
	}
}
