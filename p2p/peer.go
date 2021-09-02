package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Map of peers connected to this node (address -> peer)
var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	address string
	conn    *websocket.Conn
	inbox   chan []byte // holds outgoing messages to peer
	key     string
	port    string
}

// NON-MUTATING FUNCTIONS
// Initialize a new peer with the given connection, ip, port
func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	newPeer := &peer{
		address: address,
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		port:    port,
	}
	Peers[key] = newPeer
	go newPeer.read()  // listen to incoming messages from peer
	go newPeer.write() // listen for new outgoing messages
	return newPeer
}

// MUTATING FUNCTIONS
// Close a peer's connection and inbox channel + delete from peer list
func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key) // Will close inbox channel
}

// Continue to read and print from a peers
func (p *peer) read() {
	defer p.close() // close after function (after loop break)
	for {
		_, msg, err := p.conn.ReadMessage() // blocks until message comes
		if err != nil {
			break
		}
		fmt.Printf("%s", msg)
	}
}

// Whenever message lands in peer inbox, send to message to peer
func (p *peer) write() {
	defer p.close() // close after function (after loop break)
	for {
		m, ok := <-p.inbox // blocks until message arrives in inbox
		if !ok {           // channel no longer open
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}
