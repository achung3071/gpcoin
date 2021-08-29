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

// NON-MUTATING FUNCTIONS
// Only function that should be used to access the blockchain (b).
func Blockchain() *blockchain {
	// Singleton pattern w/ sync.Once (prevent concurrent creation of multiple blockchains)
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		chainData := db.Blockchain()
		if chainData == nil { // blockchain not in db
			// ensure AddBlock() does not call Blockchain() again,
			// or else it will result in a deadlock (circularity)
			b.AddBlock()
		} else {
			b.restore(chainData)
		}
	})
	return b
}

// Get sum of all transaction outputs for an address
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	balance := 0
	for _, txOut := range txOuts {
		balance += txOut.Amount
	}
	return balance
}

// Get all blocks
func Blocks(b *blockchain) []*Block {
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

// Save blockchain to DB
func commitBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

// Get difficulty of blockchain (i.e., how many 0s need to be in front of block hash)
func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		// no blocks yet
		return defaultDifficulty
	} else if b.Height%5 == 0 {
		// Time to recalculate & update difficulty!
		return recalculateDifficulty(b)
	} else {
		// 5 blocks not added since last update, so don't update
		return b.CurrDifficulty
	}
}

// Calculates difficulty based on whether time taken to create 5 blocks is
// too long (> 12 mins) or too short (< 8 mins)
func recalculateDifficulty(b *blockchain) int {
	blocks := Blocks(b)
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

// Get unspent transaction outputs (i.e., still valid for use as inputs) filtered by address
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut                       // holds unspent TxOuts by this address
	txsWithSpentTxOuts := make(map[string]bool) // transactions which created spent TxOuts
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, txIn := range tx.TxIns {
				// If transaction is initiated by address in question
				if txIn.Owner == address {
					// Earlier transaction (txIn.TxId) has an output that is now spent
					txsWithSpentTxOuts[txIn.TxId] = true
				}
			}
			for idx, txOut := range tx.TxOuts {
				if txOut.Owner == address {
					// Is this txOut spent (i.e., has the transaction generated a spent output)?
					outputSpent := txsWithSpentTxOuts[tx.Id]
					if !outputSpent { // output has yet to be spent
						uTxOut := UTxOut{
							TxId:   tx.Id,
							Index:  idx,
							Amount: txOut.Amount,
						}
						// Ensure output is not part of a pending tx (i.e., not on mempool)
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, &uTxOut)
						}
					}
					break // no other txOuts in this transaction belong to this address
				}
			}
		}
	}
	return uTxOuts
}

// MUTATING FUNCTIONS
// Load existing data into blockchain variable
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// Adds a new block to the blockchain & save in DB
func (b *blockchain) AddBlock() {
	newBlock := createBlock(b.LastHash, b.Height+1, getDifficulty(b))
	b.LastHash = newBlock.Hash
	b.Height = newBlock.Height
	// newBlock.Difficulty already updated using Blockchain().difficulty()
	b.CurrDifficulty = newBlock.Difficulty
	commitBlockchain(b)
}
