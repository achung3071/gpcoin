package blockchain

import (
	"time"

	"github.com/achung3071/gpcoin/utils"
)

// Basics of a transaction:
// Transaction inputs are the initial balance(s) of the party distributing the money.
// Transaction outputs are how the money has been distributed across parties after the transaction.
// A coinbase transaction is the first transaction from the blockchain to the miner, as a reward
// for mining the block (and verifying the transaction).

const (
	coinbaseAddress string = "COINBASE"
	minerAddress    string = "andrew"
	minerReward     int    = 50
)

// Holds info for one transaction
type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"TxIns"`
	TxOuts    []*TxOut `json:"TxOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

// Transaction input (how much distributing party has before transaction)
type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Transaction input (how much each party involved has after transaction)
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Creates a transaction from the blockchain that gives a reward to the miner.
func createCoinbaseTx() *Tx {
	txIns := []*TxIn{{coinbaseAddress, minerReward}}
	txOuts := []*TxOut{{minerAddress, minerReward}}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId() // attach an ID to the given transaction via hashing
	return &tx
}
