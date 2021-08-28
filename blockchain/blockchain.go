package blockchain

import (
	"sync"

	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/utils"
)

const (
	defaultDifficulty      int = 2
	updateIntervalInBlocks int = 5 // how often we should update difficulty
	expectedMinsPerBlock   int = 2 // num. mins expected for a block to be created
	updateWindowInMins     int = 2 // difficulty changes only when actual - expected time exceeds this window
)

type blockchain struct {
	LastHash       string
	Height         int
	CurrDifficulty int
}

var b *blockchain // Holds singleton instance of blockchain
var once sync.Once

// Only function that should be used to access the blockchain (b).
func Blockchain() *blockchain {
	if b == nil {
		// Singleton pattern w/ sync.Once (prevent concurrent creation of multiple blockchains)
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			chainData := db.Blockchain()
			if chainData == nil { // blockchain not in db
				b.AddBlock("Genesis block")
			} else {
				b.restore(chainData)
			}
		})
	}
	return b
}

// Calculates difficulty based on whether time taken to create x (5) blocks is
// too long (> 12 mins) or too short (< 8 mins)
func (b *blockchain) recalculateDifficulty() int {
	blocks := b.Blocks()
	newestBlock := blocks[0]
	lastUpdatedBlock := blocks[updateIntervalInBlocks-1]
	// convert from seconds -> minutes
	timeSinceLastUpdate := (newestBlock.Timestamp - lastUpdatedBlock.Timestamp) / 60
	expectedTime := updateIntervalInBlocks * expectedMinsPerBlock
	if timeSinceLastUpdate < expectedTime-updateWindowInMins {
		return b.CurrDifficulty + 1 // increase difficulty
	} else if timeSinceLastUpdate > expectedTime+updateWindowInMins {
		return b.CurrDifficulty - 1 // lower difficulty
	} else {
		return b.CurrDifficulty
	}
}

// Get difficulty of blockchain (i.e., how many 0s need to be in front of block hash)
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		// no blocks yet
		return defaultDifficulty
	} else if b.Height%5 == 0 {
		// Time to recalculate & update difficulty!
		return b.recalculateDifficulty()
	} else {
		// 5 blocks not added since last update, so don't update
		return b.CurrDifficulty
	}
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) commit() {
	db.SaveBlockchain(utils.ToBytes(b))
}

// Adds a new block to the blockchain & save in DB
func (b *blockchain) AddBlock(data string) {
	newBlock := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = newBlock.Hash
	b.Height = newBlock.Height
	// newBlock.Difficulty already updated using Blockchain().difficulty()
	b.CurrDifficulty = newBlock.Difficulty
	b.commit()
}

// Get all blocks
func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	currHash := b.LastHash
	for {
		block, err := FindBlock(currHash)
		utils.ErrorHandler(err)
		blocks = append(blocks, block)
		if block.PrevHash != "" { // not the first block
			currHash = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
