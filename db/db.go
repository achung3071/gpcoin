package db

import (
	"github.com/achung3071/gpcoin/utils"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbName           string = "blockchain.db"
	dataBucketName   string = "data"
	dataBucketKey    string = "metadata"
	blocksBucketName string = "blocks"
)

func DB() *bolt.DB {
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
	return db
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucketName))
		err := blocksBucket.Put([]byte(hash), data) // updata db with block data
		return err
	})
	utils.ErrorHandler(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucketName))
		err := dataBucket.Put([]byte(dataBucketKey), data) // updata db with chain data
		return err
	})
	utils.ErrorHandler(err)
}

// For getting an existing blockchain from the db
func Blockchain() []byte {
	var data []byte // variable to store blockchain data in
	DB().View(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucketName))
		data = dataBucket.Get([]byte(dataBucketKey))
		return nil // no error here
	})
	return data
}

// Get an existing block from the db
func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucketName))
		data = blocksBucket.Get([]byte(hash))
		return nil
	})
	return data
}

// Close database
func Close() {
	DB().Close()
}