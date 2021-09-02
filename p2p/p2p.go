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
	// Get ip address (port in r.RemoteAddr is not the open port, so no use)
	originIp := utils.Splitter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")
	// Don't allow connection if invalid ip or no open port
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return originIp != "" && openPort != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil) // return ws connection
	utils.ErrorHandler(err)
	initPeer(conn, originIp, openPort)
}

// Add a peer (initiate a websocket connection with another node)
// (e.g., :5000 requests a websocket upgrade to :4000)
func AddPeer(address, port, myPort string) {
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, myPort)
	// Request a websocket upgrade from the other node
	// (2nd argument (nil) is request header, usually w/ credentials/cookies)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.ErrorHandler(err)
	initPeer(conn, address, port) // add to list of active peers
}
