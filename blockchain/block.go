package blockchain

import (
	"crypto/sha256"
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

func createBlock(data string, prevHash string, height int) *Block {
	newBlock := &Block{data, "", prevHash, height}
	payload := newBlock.Data + newBlock.PrevHash + fmt.Sprint(newBlock.Height)
	newBlock.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	newBlock.commit()
	return newBlock
}
