package main

import (
	"github.com/achung3071/gpcoin/cli"
	"github.com/achung3071/gpcoin/db"
)

func main() {
	defer db.Close() // close db connection when program exits
	db.InitDB()
	cli.Start()
}
