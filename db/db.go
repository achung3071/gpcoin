package db

import (
	"fmt"

	"github.com/achung3071/gpcoin/utils"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbName           string = "blockchain.db"
	dataBucketName   string = "data"
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
	fmt.Printf("Saving block %s.\nData: %b\n", hash, data)
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
		err := dataBucket.Put([]byte("metadata"), data) // updata db with chain data
		return err
	})
	utils.ErrorHandler(err)
}
