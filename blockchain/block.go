package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/utils"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
}

// Save block in DB
func (b *Block) commit() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

// Load block data into block instance
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

// Can only add block to blockchain when proof of work exists (finding nonce)
func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty) // num. zeros hash must start with
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++ // increment nonce
		}
	}
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
	newBlock := &Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
		Nonce:      0,
	}
	newBlock.mine() // provide PoW
	newBlock.commit()
	return newBlock
}
