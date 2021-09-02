package p2p

import (
	"fmt"

	"github.com/achung3071/gpcoin/utils"
	"github.com/gorilla/websocket"
)

// Map of peers connected to this node (address -> peer)
var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn  *websocket.Conn
	inbox chan []byte // holds outgoing messages to peer
}

// NON-MUTATING FUNCTIONS
// Initialize a new peer with the given connection, ip, port
func initPeer(conn *websocket.Conn, address, port string) *peer {
	newPeer := &peer{conn, make(chan []byte)}
	nodeUrl := fmt.Sprintf("%s:%s", address, port)
	Peers[nodeUrl] = newPeer
	go read(newPeer)  // listen to incoming messages from peer
	go write(newPeer) // listen for new outgoing messages
	return newPeer
}

// Continue to read and print from a peer (place in goroutine)
func read(p *peer) {
	for {
		_, msg, err := p.conn.ReadMessage() // blocks until message comes
		if err != nil {
			break
		}
		fmt.Printf("%s", msg)
	}
}

// Whenever message lands in peer inbox, send to message to peer
func write(p *peer) {
	// Effective b/c sending messages doesn't block main thread
	for {
		m := <-p.inbox // blocks until message arrives in inbox
		utils.ErrorHandler(p.conn.WriteMessage(websocket.TextMessage, m))
	}
}
