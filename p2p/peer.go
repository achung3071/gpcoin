package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peer struct {
	address string
	conn    *websocket.Conn
	inbox   chan []byte // holds outgoing messages to peer
	key     string
	port    string
}

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

// Map of peers connected to this node (address -> peer)
var Peers peers = peers{v: make(map[string]*peer)}

// NON-MUTATING FUNCTIONS
// Get a list of all peer addresses to return
func AllPeers(p *peers) []string {
	p.m.Lock() // Ensure peers are not updated while reading
	defer p.m.Unlock()
	peerList := []string{}
	for key := range p.v {
		peerList = append(peerList, key)
	}
	return peerList
}

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
	Peers.v[key] = newPeer
	go newPeer.read()  // listen to incoming messages from peer
	go newPeer.write() // listen for new outgoing messages
	return newPeer
}

// MUTATING FUNCTIONS
// Close a peer's connection and inbox channel + delete from peer list
func (p *peer) close() {
	Peers.m.Lock()         // ensure peer map is locked while updating (no data race)
	defer Peers.m.Unlock() // remove lock after map is updated
	p.conn.Close()
	delete(Peers.v, p.key) // Will close inbox channel
}

// Continue to read and print from peers
func (p *peer) read() {
	defer p.close() // close after function (after loop break)
	for {
		var m Message
		err := p.conn.ReadJSON(&m) // blocks until message comes
		if err != nil {
			break
		}
		fmt.Println(m.Type)
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
