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
		blocks := []*Block{{PrevHash: "x"}, {PrevHash: ""}}
		currBlock := 0
		dbStorage = mockDB{mockFindBlock: func(string) []byte {
			defer func() { currBlock++ }()
			return utils.ToBytes(blocks[currBlock])
		}}
		blocksResult := Blocks(&blockchain{LastHash: "y"})
		if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
			t.Error("Blocks() did not return a slice of blocks")
		} else if len(blocksResult) != 2 {
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

func TestGetDifficulty(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: ""},
	}
	// Needed b/c recalculateDifficulty calls Blocks(), which calls FindBlock()
	currBlock := 0
	dbStorage = mockDB{mockFindBlock: func(string) []byte {
		defer func() { currBlock++ }()
		return utils.ToBytes(blocks[currBlock])
	}}
	t.Run("Should only update difficulty when height is a multiple of the update interval", func(t *testing.T) {
		type test struct {
			height         int
			expectedOutput int
		}
		tests := []test{
			{height: 0, expectedOutput: defaultDifficulty},
			{height: 2, expectedOutput: defaultDifficulty},
			{height: 5, expectedOutput: defaultDifficulty + 1},
		}
		for _, tc := range tests {
			bc := &blockchain{Height: tc.height, CurrDifficulty: defaultDifficulty}
			result := getDifficulty(bc)
			if result != tc.expectedOutput {
				t.Errorf("getDifficulty() should return %d got %d", tc.expectedOutput, result)
			}
		}
	})
}

func TestAddBlockFromPeer(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	dbStorage = mockDB{}

	bc := &blockchain{Height: 1, CurrDifficulty: 1, LastHash: "xx"}
	Mempool().Txs["test"] = &Tx{} // ensure this tx is removed from mempool
	newBlock := &Block{Difficulty: 2, Hash: "test", Transactions: []*Tx{{Id: "test"}}}
	bc.AddBlockFromPeer(newBlock)

	t.Run("AddBlockFromPeer() should update the blockchain", func(t *testing.T) {
		if bc.CurrDifficulty != 2 || bc.Height != 2 || bc.LastHash != "test" {
			t.Error("AddBlockFromPeer() did not update the blockchain with new block's data")
		}
	})
	t.Run("AddBlockFromPeer() should remove transactions from the mempool", func(t *testing.T) {
		_, ok := Mempool().Txs["test"]
		if ok {
			t.Error("AddBlockFromPeer() should have removed transaction id 'test' from mempool")
		}
	})

}

func TestReplace(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	dbStorage = mockDB{}
	t.Run("Replace() should mutate the blockchain", func(t *testing.T) {
		bc := &blockchain{Height: 1, CurrDifficulty: 1, LastHash: "xx"}
		blocks := []*Block{{Difficulty: 2, Hash: "test"}, {Difficulty: 2, Hash: "test"}}
		bc.Replace(blocks)
		if bc.CurrDifficulty != 2 || bc.Height != 2 || bc.LastHash != "test" {
			t.Error("Replace() did not update the blockchain with the new blocks")
		}
	})
}
