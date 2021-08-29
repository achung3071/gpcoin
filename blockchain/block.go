package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
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
		fmt.Printf("\n\nTarget: %s\nHash: %s\nNonce: %d\n\n", target, hash, b.Nonce)
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

func createBlock(prevHash string, height int) *Block {
	// Initialize every new block added to chain w/ a coinbase transaction
	newBlock := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
		Nonce:      0,
	}
	newBlock.mine() // provide PoW
	// flush mempool and get confirmed transactions
	newBlock.Transactions = Mempool.ConfirmTxs()
	newBlock.commit()
	return newBlock
}
