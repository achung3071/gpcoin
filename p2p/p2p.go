package p2p

import (
	"fmt"
	"net/http"

	"github.com/achung3071/gpcoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{}

// Upgrade http request to websocket connection
// (e.g., :4000 accepts a websocket upgrade request from :5000)
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Normally you want to be more careful in allowing certain origins
	// in making a websocket connection, but here we allow all connections
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(rw, r, nil) // return ws connection
	utils.ErrorHandler(err)
	initPeer(conn, "xx", "xx") // TODO: change to use address/port of request
}

// Add a peer (initiate a websocket connection with another node)
// (e.g., :5000 requests a websocket upgrade to :4000)
func AddPeer(address, port string) {
	url := fmt.Sprintf("ws://%s:%s/ws", address, port)
	// Request a websocket upgrade from the other node
	// (2nd argument (nil) is request header, usually w/ credentials/cookies)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.ErrorHandler(err)
	initPeer(conn, address, port) // add to list of active peers
}
