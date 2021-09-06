package blockchain

import (
	"reflect"
	"testing"
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
