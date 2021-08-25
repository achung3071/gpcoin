package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

// Main function
func main() {
	genesisBlock := block{"First block", "", ""}
	// hash data with prev hash (currently non-existent) using SHA256
	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
	genesisBlock.hash = fmt.Sprintf("%x", hash) // hex number as hash
	fmt.Println(genesisBlock)
}
