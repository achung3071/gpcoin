package blockchain

import (
	"sync"
)

type blockchain struct {
	LastHash string
	Height   int
}

var b *blockchain // Holds singleton instance of blockchain
var once sync.Once

// Only function that should be used to access the blockchain (b).
func Blockchain() *blockchain {
	if b == nil {
		// Singleton pattern w/ sync.Once (prevent concurrent creation of multiple blockchains)
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("First block")
		})
	}
	return b
}

// Adds a new block to the blockchain & save in DB
func (b *blockchain) AddBlock(data string) {
	newBlock := createBlock(data, b.LastHash, b.Height)
	b.LastHash = newBlock.Hash
	b.Height = newBlock.Height
}
