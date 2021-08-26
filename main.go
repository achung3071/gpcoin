package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/achung3071/gpcoin/blockchain"
)

const port string = ":8080"
const tempDir string = "templates/"

var templates *template.Template

// need uppercase fields to be able to access in template
type tempData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// basic handler for route
func home(rw http.ResponseWriter, r *http.Request) {
	chain := blockchain.GetBlockchain()
	templates.ExecuteTemplate(rw, "home", tempData{"GPCoin Blockchain", chain.GetBlocks()})
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil) // no data to pass
	case "POST":
		r.ParseForm()              // form input will be passed in request
		data := r.Form.Get("data") // get block data
		blockchain.GetBlockchain().AddBlock(data)
		// Redirect client back to homepage, where they will see new block
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func main() {
	// template.Must does automatic error handling (log.Panic) for errors
	// that are returned when parsing a template file
	templates = template.Must(template.ParseGlob(tempDir + "pages/*.html"))     // get pages
	templates = template.Must(templates.ParseGlob(tempDir + "partials/*.html")) // get partials
	fmt.Printf("Listening on http://localhost%s\n", port)
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(port, nil)) // log when ListenAndServe returns an error
}
