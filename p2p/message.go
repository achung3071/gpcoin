package p2p

import (
	"encoding/json"

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
	m := Message{Type: msgType}
	m.addPayload(payload)
	mBytes, err := json.Marshal(m)
	utils.ErrorHandler(err)
	return mBytes
}

// MUTATING FUNCTIONS
// Add payload in bytes to message instance
func (m *Message) addPayload(payload interface{}) {
	bytes, err := json.Marshal(payload)
	utils.ErrorHandler(err)
	m.Payload = bytes
}
