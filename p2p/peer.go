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

func initPeer(conn *websocket.Conn, address, port string) {
	newPeer := &peer{conn}
	nodeUrl := fmt.Sprintf("%s:%s", address, port)
	Peers[nodeUrl] = newPeer
}
