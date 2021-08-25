package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block // use pointers to prevent copying blocks
}

var b *blockchain // Holds singleton instance of blockchain
var once sync.Once

// This is the only function that should be used to access the blockchain (b).
func GetBlockchain() *blockchain {
	if b == nil {
		// Singleton pattern w/ sync.Once (prevent concurrent creation of multiple blockchains)
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("First block")
		})
	}
	return b
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) == 0 {
		return ""
	}
	return b.blocks[len(b.blocks)-1].Hash
}

func (b *block) calcHash() {
	// hash data with prev hash  using SHA256 - need to convert string to []byte
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash) // hex number as hash
}

func createBlock(data string) *block {
	chain := GetBlockchain()
	newBlock := block{data, "", chain.getLastHash()}
	newBlock.calcHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockchain) GetBlocks() []*block {
	return b.blocks
}
