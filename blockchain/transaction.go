package blockchain

import (
	"errors"
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
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
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

// Mempool is where unconfirmed transactions are (before added to a block)
type Mempool struct {
	Txs []*Tx `json:"txs"`
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

// Create a new transaction from one address to anotehr
func makeTx(from string, to string, amount int) (*Tx, error) {
	balance := Blockchain().BalanceByAddress(from)
	if balance < amount {
		return nil, errors.New("not enough money to send")
	}
	// Make new transaction inputs from prev transaction outputs
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	runningTotal := 0
	txIns := []*TxIn{}
	txOuts := []*TxOut{}
	for _, txOut := range oldTxOuts {
		if runningTotal >= amount {
			break // enough transaction inputs
		}
		// Else, add current txOut to txIns
		runningTotal += txOut.Amount
		newTxIn := TxIn{from, txOut.Amount}
		txIns = append(txIns, &newTxIn)
	}
	// Make transaction outputs
	change := runningTotal - amount
	if change > 0 { // change needs to be one of the transaction outputs
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	// Make final transaction
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId() // hash for transaction id
	return &tx, nil
}

// Add a transaction to a certain address on the mempool
func (m *Mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(minerAddress, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}
