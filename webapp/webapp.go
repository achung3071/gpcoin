package webapp

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/achung3071/gpcoin/blockchain"
)

const tempDir string = "webapp/templates/"

var templates *template.Template

// need uppercase fields to be able to access in template
type tempData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// basic handler for route
func home(rw http.ResponseWriter, r *http.Request) {
	data := tempData{"GPCoin Blockchain", blockchain.Blockchain().Blocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil) // no data to pass
	case "POST":
		r.ParseForm()              // form input will be passed in request
		data := r.Form.Get("data") // get block data
		blockchain.Blockchain().AddBlock(data)
		// Redirect client back to homepage, where they will see new block
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(portNum int) {
	// Ensure that diff. multiplexers (thing which calls handler funcs based on request url)
	// are used for web & rest API, so that there is no error regarding duplicate endpoints
	handler := http.NewServeMux()
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)

	// template.Must does automatic error handling (log.Panic) when parsing a template file
	templates = template.Must(template.ParseGlob(tempDir + "pages/*.html"))     // get pages
	templates = template.Must(templates.ParseGlob(tempDir + "partials/*.html")) // get partials

	port := fmt.Sprintf(":%d", portNum)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler)) // log when ListenAndServe returns an error
}
