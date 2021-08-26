package main

import (
	"github.com/achung3071/gpcoin/api"
	"github.com/achung3071/gpcoin/webapp"
)

func main() {
	go webapp.Start(4000)
	api.Start(5000)
}
