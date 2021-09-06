package db

import (
	"fmt"

	"github.com/achung3071/gpcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dataBucketName   string = "data"
	dataBucketKey    string = "metadata"
	blocksBucketName string = "blocks"
)

// Struct to implement "storage" interface from blockchain pkg.
type BoltDB struct{}

func (BoltDB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (BoltDB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (BoltDB) EmptyBlocks() {
	emptyBlocks()
}
func (BoltDB) SaveBlockchain(data []byte) {
	saveBlockchain(data)
}
func (BoltDB) LoadBlockchain() []byte {
	return loadBlockchain()
}

var db *bolt.DB
var dbName string = "blockchain.db"

// DB name reset to include port when cli.Start() called
func SetDBName(port int) {
	dbName = fmt.Sprintf("blockchain_%d.db", port)
}

// Initialize database connection on program start
func InitDB() {
	if db == nil {
		// create db (chmod 0600 is read-write permission)
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.ErrorHandler(err)
		db = dbPointer
		err = db.Update(func(t *bolt.Tx) error {
			// returns *bucket and error (no need for bucket right now)
			_, err := t.CreateBucketIfNotExists([]byte(dataBucketName))
			utils.ErrorHandler(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucketName))
			return err
		})
		utils.ErrorHandler(err)
	}
}

// Close database connection
func Close() {
	db.Close()
}

// Remove blocks from blocks bucket in db
func emptyBlocks() {
	err := db.Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket([]byte(blocksBucketName))
		if err != nil {
			return err
		}
		_, err = t.CreateBucket([]byte(blocksBucketName))
		return err
	})
	utils.ErrorHandler(err)
}

// Get an existing block from the db
func findBlock(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucketName))
		data = blocksBucket.Get([]byte(hash))
		return nil
	})
	return data
}

// Load blockchain metadata from db
func loadBlockchain() []byte {
	var data []byte // variable to store blockchain data in
	db.View(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucketName))
		data = dataBucket.Get([]byte(dataBucketKey))
		return nil // no error here
	})
	return data
}

// Save blockchain metadata to db
func saveBlockchain(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucketName))
		err := dataBucket.Put([]byte(dataBucketKey), data) // updata db with chain data
		return err
	})
	utils.ErrorHandler(err)
}

// Save a single block to the db
func saveBlock(hash string, data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucketName))
		err := blocksBucket.Put([]byte(hash), data) // updata db with block data
		return err
	})
	utils.ErrorHandler(err)
}
