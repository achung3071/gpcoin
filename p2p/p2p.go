package p2p

import (
	"fmt"
	"net/http"

	"github.com/achung3071/gpcoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{}

// Upgrade http request to websocket connection
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Normally you want to be more careful in allowing certain origins
	// in making a websocket connection, but here we allow all connections
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(rw, r, nil) // return ws connection
	utils.ErrorHandler(err)
	// Read every message that comes in - continue listening
	for {
		_, msg, err := conn.ReadMessage()
		utils.ErrorHandler(err)
		fmt.Printf("New message:  %s\n\n", msg)
	}
}
