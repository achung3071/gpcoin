package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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

// Convert block to bytes
func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer
	err := gob.NewEncoder(&blockBuffer).Encode(b)
	utils.ErrorHandler(err)
	return blockBuffer.Bytes()
}

// Save block in DB
func (b *Block) commit() {
	db.SaveBlock(b.Hash, b.toBytes())
}

func createBlock(data string, prevHash string, height int) *Block {
	newBlock := &Block{data, "", prevHash, height + 1}
	payload := newBlock.Data + newBlock.PrevHash + fmt.Sprint(newBlock.Height)
	newBlock.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	newBlock.commit()
	return newBlock
}
