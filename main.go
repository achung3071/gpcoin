package main

import (
	"github.com/achung3071/gpcoin/blockchain"
	"github.com/achung3071/gpcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}
