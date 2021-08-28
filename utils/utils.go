package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(i)
	ErrorHandler(err)
	return buffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(i)
	ErrorHandler(err)
}

func Hash(i interface{}) string {
	iString := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(iString))
	return fmt.Sprintf("%x", hash)
}
