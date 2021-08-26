package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/achung3071/gpcoin/blockchain"
)

const port string = ":8080"

// need uppercase fields to be able to access in template
type tempData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// basic handler for route
func home(rw http.ResponseWriter, r *http.Request) {
	// template.Must does automatic error handling (log.Panic) for errors
	// that are returned when parsing a template file
	temp := template.Must(template.ParseFiles("templates/home.gohtml"))
	chain := blockchain.GetBlockchain()
	temp.Execute(rw, tempData{"GPCoin", chain.GetBlocks()})
}

func main() {
	fmt.Printf("Listening on http://localhost%s\n", port)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(port, nil)) // log when ListenAndServe returns an error
}
