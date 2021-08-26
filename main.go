package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":5000"

type URLDescription struct {
	URL         string `json:"url"` // struct field tag -> renames based on encoding
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omits this field when non-existent
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	urls := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "Documentation of all endpoints",
			Payload:     "",
		},
		{
			URL:         "/blocks",
			Method:      "POST",
			Description: "Add a block",
			Payload:     "string for data field",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(urls) // easy way to send json to writer
}

func main() {
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
