package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/achung3071/gpcoin/utils"
)

type mockDB struct {
	mockLoadBlockchain func() []byte
	mockFindBlock      func(hash string) []byte
}

func (m mockDB) FindBlock(hash string) []byte {
	return m.mockFindBlock(hash)
}
func (m mockDB) LoadBlockchain() []byte {
	return m.mockLoadBlockchain()
}
func (mockDB) SaveBlock(hash string, data []byte) {}
func (mockDB) SaveBlockchain(data []byte)         {}
func (mockDB) EmptyBlocks()                       {}

func TestBlockchain(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	t.Run("Blockchain() should create blockchain when blockchain is nil", func(t *testing.T) {
		once = *new(sync.Once) // ensure that code in Blockchain() can be run multiple times
		dbStorage = mockDB{mockLoadBlockchain: func() []byte { return nil }}
		b := Blockchain()
		if b.Height != 1 {
			t.Error("Blockchain() did not create a brand new blockchain")
		}
	})
	t.Run("Blockchain() return existing blockchain when available", func(t *testing.T) {
		once = *new(sync.Once) // ensure that code in Blockchain() can be run multiple times
		dbStorage = mockDB{mockLoadBlockchain: func() []byte {
			return utils.ToBytes(&blockchain{LastHash: "", Height: 2, CurrDifficulty: 1})
		}}
		b := Blockchain()
		if b.Height != 2 {
			t.Errorf("Expected blockchain of height 2, got height %d", b.Height)
		}
	})
}

func TestBlocks(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	t.Run("Blocks() should return slice of blocks", func(t *testing.T) {
		blocksAdded := 0
		dbStorage = mockDB{mockFindBlock: func(string) []byte {
			var block *Block
			if blocksAdded == 0 {
				block = &Block{Hash: "y", PrevHash: "x"}
			} else if blocksAdded == 1 {
				block = &Block{Hash: "x", PrevHash: ""}
			}
			blocksAdded++
			return utils.ToBytes(block)
		}}
		blocks := Blocks(&blockchain{LastHash: "y"})
		if reflect.TypeOf(blocks) != reflect.TypeOf([]*Block{}) {
			t.Error("Blocks() did not return a slice of blocks")
		} else if len(blocks) != 2 {
			t.Errorf("Expected Blocks() to return a slice of length 2, got %d", len(blocks))
		}
	})
}

func TestFindTx(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	t.Run("FindTx() should return nil when transaction doesn't exist", func(t *testing.T) {
		dbStorage = mockDB{mockFindBlock: func(string) []byte {
			block := &Block{Hash: "y", Transactions: []*Tx{}}
			return utils.ToBytes(block)
		}}
		tx := FindTx(&blockchain{LastHash: "y"}, "test")
		if tx != nil {
			t.Errorf("Expected transaction to be nil, got txId %s", tx.Id)
		}
	})
	t.Run("FindTx() should return existing transaction", func(t *testing.T) {
		dbStorage = mockDB{
			mockFindBlock: func(string) []byte {
				block := &Block{Hash: "y", Transactions: []*Tx{{Id: "test"}}}
				return utils.ToBytes(block)
			},
		}
		tx := FindTx(&blockchain{LastHash: "y"}, "test")
		if tx == nil {
			t.Error("Existing transaction not found.")
		} else if tx.Id != "test" {
			t.Errorf("Expected transaction with id 'test', got %s", tx.Id)
		}
	})
}
