package blockchain

import (
	"sync"
	"testing"

	"github.com/achung3071/gpcoin/utils"
)

type mockDB struct {
	mockLoadChain func() []byte
	mockFindBlock func(hash string) []byte
}

func (m mockDB) FindBlock(hash string) []byte {
	return m.mockFindBlock(hash)
}
func (m mockDB) LoadBlockchain() []byte {
	return m.mockLoadChain()
}
func (mockDB) SaveBlock(hash string, data []byte) {}
func (mockDB) SaveBlockchain(data []byte)         {}
func (mockDB) EmptyBlocks()                       {}

func TestBlockchain(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	t.Run("Blockchain() should create blockchain when blockchain is nil", func(t *testing.T) {
		once = *new(sync.Once) // ensure that code in Blockchain() can be run multiple times
		dbStorage = mockDB{mockLoadChain: func() []byte { return nil }}
		b := Blockchain()
		if b.Height != 1 {
			t.Error("Blockchain() did not create a brand new blockchain")
		}
	})
	t.Run("Blockchain() return existing blockchain when available", func(t *testing.T) {
		once = *new(sync.Once) // ensure that code in Blockchain() can be run multiple times
		dbStorage = mockDB{mockLoadChain: func() []byte {
			return utils.ToBytes(&blockchain{LastHash: "", Height: 2, CurrDifficulty: 1})
		}}
		b := Blockchain()
		if b.Height != 2 {
			t.Errorf("Expected blockchain of height 2, got height %d", b.Height)
		}
	})
}
