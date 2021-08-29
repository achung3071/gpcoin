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

// Transaction input (previous transaction output that is being spent)
type TxIn struct {
	TxId  string `json:"txId"`  // transaction which created the TxOut (spent as this input)
	Index int    `json:"index"` // index of TxOut within transaction
	Owner string `json:"owner"`
}

// Transaction output (how much each party involved has after transaction)
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Unspent transaction output
type UTxOut struct {
	TxId   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

// Mempool is where unconfirmed transactions are (before added to a block)
type mempool struct {
	Txs []*Tx `json:"txs"`
}

var Mempool *mempool = &mempool{}

// NON-MUTATING FUNCTIONS
// Creates a transaction from the blockchain that gives a reward to the miner
func createCoinbaseTx() *Tx {
	txIns := []*TxIn{{"", -1, coinbaseAddress}}
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

// Checks if a uTxOut is on the mempool already (so it isn't passed as an input again)
func isOnMempool(uTxOut UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxId == uTxOut.TxId && txIn.Index == uTxOut.Index {
				exists = true // uTxOut is already being used on the mempool
				break Outer   // using labels, break from outer for loop
			}
		}
	}
	return exists
}

// Create a new transaction from one address to another
func makeTx(from string, to string, amount int) (*Tx, error) {
	currBalance := BalanceByAddress(from, Blockchain())
	if currBalance < amount {
		return nil, errors.New("not enough funds to send specified amount")
	}
	txIns := []*TxIn{}
	txOuts := []*TxOut{}
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	// Append transaction inputs
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break // enoungh TxIns added
		}
		total += uTxOut.Amount
		txIns = append(txIns, &TxIn{uTxOut.TxId, uTxOut.Index, from})
	}
	// Create transaction outputs
	if change := total - amount; change > 0 {
		// give change back as a transaction output
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	// Return final transaction
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx, nil
}

// MUTATING FUNCTIONS
// Populates id field of a transaction
func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

// Add a transaction to a certain address on the mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(minerAddress, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

// Empties mempool and returns now-confirmed transactions
func (m *mempool) ConfirmTxs() []*Tx {
	// reward for mining new block & confirming transactions
	coinbaseTx := createCoinbaseTx()
	txs := m.Txs
	txs = append(txs, coinbaseTx)
	m.Txs = []*Tx{} // empty mempool
	return txs
}
