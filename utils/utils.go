package utils

import (
	"bytes"
	"encoding/gob"
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
