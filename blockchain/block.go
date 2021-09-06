package blockchain

import (
	"errors"
	"strings"
	"time"

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

var ErrBlockNotFound error = errors.New("Block with given hash not found")

// NON-MUTATING FUNCTIONS
// Save block in DB
func commitBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b))
}

// Create a new block (mine and add mempool transactions)
func createBlock(prevHash string, height int, diff int) *Block {
	// Initialize every new block added to chain w/ a coinbase transaction
	newBlock := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	newBlock.mine() // provide PoW
	// flush mempool and get confirmed transactions
	newBlock.Transactions = Mempool().ConfirmTxs()
	commitBlock(newBlock)
	return newBlock
}

// Find block from DB based on hash
func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.FindBlock(hash)
	if blockBytes == nil { // non-existent
		return nil, ErrBlockNotFound
	}
	block := &Block{}         // init empty block
	block.restore(blockBytes) // load block data
	return block, nil
}

// MUTATING FUNCTIONS
// Give proof of work (find nonce) to add block to blockchain
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

// Load block data into block instance
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}
