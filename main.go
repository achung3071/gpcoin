package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":5000"

type URL string // custom type

// This is a built-in interface in the "encoding" package with
// a method that the json encoder uses to encode text. We can have
// the URL type implement this interface to change how it is encoded.
/* type TextMarshaler interface {
	MarshalText() (text []byte, err error)
} */
func (u URL) MarshalText() ([]byte, error) {
	fullUrl := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(fullUrl), nil
}

type URLDescription struct {
	URL         URL    `json:"url"` // struct field tag -> renames based on encoding
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omits this field when non-existent
}

// Another common interface: Stringer interface (implements String())
// "fmt" uses String() to display Go structs/instances when printed
/* func (u URLDescription) String() string {
	return "This is a url description"
} */

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
