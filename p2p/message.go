package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/achung3071/gpcoin/blockchain"
	"github.com/achung3071/gpcoin/utils"
)

type MessageType int

type Message struct {
	Type    MessageType
	Payload []byte
}

const (
	MessageNewestBlock MessageType = iota // enumerate following constants
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

// NON-MUTATING FUNCTIONS
// Return message with the given type and payload in JSON format
func makeMessage(msgType MessageType, payload interface{}) []byte {
	m := Message{Type: msgType, Payload: utils.ToJSON(payload)}
	return utils.ToJSON(m)
}

// Handle an incoming message from a peer
func handleMessage(m *Message, p *peer) {
	switch m.Type {
	case MessageNewestBlock:
		var newestBlock blockchain.Block
		json.Unmarshal(m.Payload, &newestBlock)
		fmt.Println(newestBlock)
	}
}
