package blockchain

import (
	"reflect"
	"testing"

	"github.com/achung3071/gpcoin/utils"
)

func TestCreateBlock(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	dbStorage = mockDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	t.Run("createBlock() should return a block", func(t *testing.T) {
		if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
			t.Error("createBlock() did not return an instance of a block")
		}
	})
	t.Run("createBlock() should include mempool and coinbase transactions", func(t *testing.T) {
		if len(b.Transactions) != 2 {
			t.Errorf("Expected %d transactions in block, received %d", 2, len(b.Transactions))
		}
	})
}

func TestFindBlock(t *testing.T) {
	oldStorage := dbStorage
	defer func() { dbStorage = oldStorage }()
	t.Run("FindBlock() should error when block doesn't exist", func(t *testing.T) {
		dbStorage = mockDB{mockFindBlock: func(string) []byte { return nil }}
		_, err := FindBlock("xx")
		if err == nil {
			t.Error("FindBlock() did not error even though block does not exist")
		}
	})
	t.Run("FindBlock() should return a block with the correct data", func(t *testing.T) {
		dbStorage = mockDB{mockFindBlock: func(string) []byte {
			b := &Block{Height: 1}
			return utils.ToBytes(b)
		}}
		b, _ := FindBlock("xx")
		if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
			t.Error("FindBlock() did not return a block instance")
		} else if b.Height != 1 {
			t.Errorf("Expected block of height 1, got %d", b.Height)
		}
	})
}
