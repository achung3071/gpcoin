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

type blockchain struct {
	blocks []block
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) == 0 {
		return ""
	}
	return b.blocks[len(b.blocks)-1].hash
}

func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	// hash data with prev hash  using SHA256 - need to convert string to []byte
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hash) // hex number as hash
	b.blocks = append(b.blocks, newBlock)
}

func (b blockchain) listBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("Prev hash: %s\n\n", block.prevHash)
	}
}

// Main function
func main() {
	b := blockchain{}
	b.addBlock("First block")
	b.addBlock("Second block")
	b.addBlock("Third block")
	b.listBlocks()
}
