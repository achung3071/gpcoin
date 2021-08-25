package main

import (
	"fmt"

	"github.com/achung3071/gpcoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")
	for _, block := range chain.GetBlocks() {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Prev hash: %s\n\n", block.PrevHash)
	}
}
