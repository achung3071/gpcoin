package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achung3071/gpcoin/utils"
)

const port string = ":5000"

type URLDescription struct {
	URL         string
	Method      string
	Description string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	urls := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "Documentation of all endpoints",
		},
	}
	b, err := json.Marshal(urls) // returns JSON in byte format
	utils.ErrorHandler(err)
	fmt.Fprintf(rw, "%s", b)
}

func main() {
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
