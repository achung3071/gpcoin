package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block // use pointers to prevent copying blocks
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

func (b *Block) calcHash() {
	// hash data with prev hash  using SHA256 - need to convert string to []byte
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash) // hex number as hash
}

func createBlock(data string) *Block {
	chain := GetBlockchain()
	newBlock := Block{data, "", chain.getLastHash(), len(chain.GetBlocks()) + 1}
	newBlock.calcHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockchain) GetBlocks() []*Block {
	return b.blocks
}

var ErrBlockNotFound error = errors.New("Block not found")

func (b *blockchain) GetBlock(height int) (*Block, error) {
	if height <= 0 || height > len(b.blocks) {
		return nil, ErrBlockNotFound
	}
	return b.blocks[height-1], nil
}
