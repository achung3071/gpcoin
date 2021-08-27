package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

// Save block in DB
func (b *Block) commit() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

// Load block data into block instance
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrBlockNotFound error = errors.New("Block with given hash not found")

// Find block from DB based on hash
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil { // non-existent
		return nil, ErrBlockNotFound
	}
	block := &Block{}         // init empty block
	block.restore(blockBytes) // load block data
	return block, nil
}

func createBlock(data string, prevHash string, height int) *Block {
	newBlock := &Block{data, "", prevHash, height}
	payload := newBlock.Data + newBlock.PrevHash + fmt.Sprint(newBlock.Height)
	newBlock.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	newBlock.commit()
	return newBlock
}
