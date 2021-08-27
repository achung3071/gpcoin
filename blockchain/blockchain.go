package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/utils"
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
			chainData := db.GetBlockchain()
			if chainData == nil { // blockchain not in db
				b.AddBlock("Genesis block")
			} else {
				fmt.Println("Restoring blockchain...")
				b.restore(chainData)
			}
		})
	}
	fmt.Printf("LastHash: %s\nHeight:%d\n", b.LastHash, b.Height)
	return b
}

func (b *blockchain) restore(data []byte) {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(b) // restore blockchain data
	utils.ErrorHandler(err)
}

func (b *blockchain) commit() {
	db.SaveBlockchain(utils.ToBytes(b))
}

// Adds a new block to the blockchain & save in DB
func (b *blockchain) AddBlock(data string) {
	newBlock := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = newBlock.Hash
	b.Height = newBlock.Height
	b.commit()
}
