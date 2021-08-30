package blockchain

import (
	"errors"
	"time"

	"github.com/achung3071/gpcoin/utils"
	"github.com/achung3071/gpcoin/wallet"
)

// Basics of a transaction:
// Transaction inputs are the initial balance(s) of the party distributing the money.
// Transaction outputs are how the money has been distributed across parties after the transaction.
// A coinbase transaction is the first transaction from the blockchain to the miner, as a reward
// for mining the block (and verifying the transaction).

const (
	coinbaseAddress string = "COINBASE"
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
	TxId      string `json:"txId"`      // transaction which created the TxOut (spent as this input)
	Index     int    `json:"index"`     // index of TxOut within transaction
	Signature string `json:"signature"` // signature by person creating the transaction
}

// Transaction output (how much each party involved has after transaction)
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
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
var errNoMoney error = errors.New("not enough funds to send specified amount")
var errInvalidTx error = errors.New("inputs are not valid txOuts for the given wallet")

// NON-MUTATING FUNCTIONS
// Creates a transaction from the blockchain that gives a reward to the miner
func createCoinbaseTx() *Tx {
	txIns := []*TxIn{{"", -1, coinbaseAddress}}
	txOuts := []*TxOut{{wallet.Wallet().Address, minerReward}}
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
		return nil, errNoMoney
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
	tx.getId()             // hash transaction to populate id
	tx.sign()              // sign all inputs in transaction
	valid := validate(&tx) // ensure transaction inputs are valid
	if !valid {
		return nil, errInvalidTx
	}
	return &tx, nil
}

// Validate a transaction (i.e., that the wallet owner owns
// the transaction outputs that are now used as inputs)
func validate(tx *Tx) bool {
	for _, txIn := range tx.TxIns {
		// Find the transaction that created the transaction input
		prevTx := FindTx(Blockchain(), txIn.TxId)
		if prevTx == nil {
			// Fake input: not on the blockchain
			return false
		}
		prevTxOut := prevTx.TxOuts[txIn.Index]
		address := prevTxOut.Address
		// If the public key (address) cannot verify the signature that I just
		// created w/ my wallet, that means the TxOuts/funds are not actually mine
		if !wallet.Verify(tx.Id, txIn.Signature, address) {
			return false
		}
	}
	return true // All transaction inputs verfied
}

// MUTATING FUNCTIONS
// Add a transaction to a certain address on the mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
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

// Populates id field of a transaction
func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

// Sign all transaction inputs in a transaction
func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		// sign transaction id, which is a hash of the transaction
		txIn.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}
