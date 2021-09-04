package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/achung3071/gpcoin/blockchain"
	"github.com/achung3071/gpcoin/utils"
)

type MessageType int

// Information to send when broadcasting peer
type BroadcastPeerInfo struct {
	NewPeerAddress string
	NewPeerPort    string
	ReceivingPort  string
}

type Message struct {
	Type    MessageType
	Payload []byte
}

const (
	MessageNewestBlock MessageType = iota // enumerate following constants
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNotifyNewBlock
	MessageNotifyNewPeer
	MessageNotifyNewTx
)

// NON-MUTATING FUNCTIONS
// Return message with the given type and payload in JSON format
func makeMessage(msgType MessageType, payload interface{}) []byte {
	m := Message{Type: msgType, Payload: utils.ToJSON(payload)}
	return utils.ToJSON(m)
}

// Broadcast a newly mined block to all peers
func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		m := makeMessage(MessageNotifyNewBlock, b)
		p.inbox <- m
	}
}

// Broadcast a newly added peer to all other peers, so they add new peer as well
func BroadcastNewPeer(newPeer *peer) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		if p.key != newPeer.key { // if peer is not the newly added peer
			payload := &BroadcastPeerInfo{
				NewPeerAddress: newPeer.address,
				NewPeerPort:    newPeer.port,
				ReceivingPort:  p.port,
			}
			m := makeMessage(MessageNotifyNewPeer, payload)
			p.inbox <- m
		}
	}
}

// Broadcast a newly posted transaction to all peers
func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		m := makeMessage(MessageNotifyNewTx, tx)
		p.inbox <- m
	}
}

// Handle an incoming message from a peer
func handleMessage(m *Message, p *peer) {
	switch m.Type {
	case MessageNewestBlock:
		fmt.Printf("Received newest block from %s.\n", p.key)
		var payload blockchain.Block
		utils.ErrorHandler(json.Unmarshal(m.Payload, &payload))
		latestBlock, err := blockchain.FindBlock(blockchain.Blockchain().LastHash)
		utils.ErrorHandler(err)
		if payload.Height >= latestBlock.Height { // our node is behind, so request blocks
			requestAllBlocks(p)
		} else { // our node is ahead, so send newest block to let them know they are behind
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("Received a request for all blocks from %s.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all blocks from the blockchain of %s.\n", p.key)
		var payload []*blockchain.Block
		utils.ErrorHandler(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)
	case MessageNotifyNewBlock:
		var payload *blockchain.Block
		utils.ErrorHandler(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddBlockFromPeer(payload)
	case MessageNotifyNewPeer:
		var payload BroadcastPeerInfo
		utils.ErrorHandler(json.Unmarshal(m.Payload, &payload))
		// broadcast is false b/c new peer has already been broadcasted to other peers
		AddPeer(payload.NewPeerAddress, payload.NewPeerPort, payload.ReceivingPort, false)
	case MessageNotifyNewTx:
		var payload *blockchain.Tx
		utils.ErrorHandler(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddTxFromPeer(payload)
	}

}

// Request all blocks from peer
func requestAllBlocks(p *peer) {
	fmt.Printf("Requesting %s for all blocks...\n", p.key)
	msgJson := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- msgJson
}

// Send all blocks to peer
func sendAllBlocks(p *peer) {
	fmt.Printf("Sending %s all blocks in our blockchain...\n", p.key)
	blocks := blockchain.Blocks(blockchain.Blockchain())
	msgJson := makeMessage(MessageAllBlocksResponse, blocks)
	p.inbox <- msgJson
}

// Send newest block to the peer
func sendNewestBlock(p *peer) {
	fmt.Printf("Sending %s the newest block in our blockchain...\n", p.key)
	newestBlock, err := blockchain.FindBlock(blockchain.Blockchain().LastHash)
	utils.ErrorHandler(err)
	msgJson := makeMessage(MessageNewestBlock, newestBlock)
	p.inbox <- msgJson
}
