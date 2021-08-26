package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achung3071/gpcoin/blockchain"
)

var port string

type url string // custom type

// This is a built-in interface in the "encoding" package with
// a method that the json encoder uses to encode text. We can have
// the url type implement this interface to change how it is encoded.
/* type TextMarshaler interface {
	MarshalText() (text []byte, err error)
} */
func (u url) MarshalText() ([]byte, error) {
	fullUrl := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(fullUrl), nil
}

type urlDescription struct {
	URL         url    `json:"url"` // struct field tag -> renames based on encoding
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omits this field when non-existent
}

// Another common interface: Stringer interface (implements String())
// "fmt" uses String() to display Go structs/instances when printed
/* func (u urlDescription) String() string {
	return "This is a url description"
} */

func documentation(rw http.ResponseWriter, r *http.Request) {
	urls := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "Documentation of all endpoints",
			Payload:     "",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "Get all blocks",
			Payload:     "",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "{data: string}",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(urls) // easy way to send json to writer
}

type postBlocksBody struct {
	Data string
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().GetBlocks())
	case "POST":
		var blockData postBlocksBody
		// save body in blockData variable (decoder automatically maps lowercase data -> Data)
		json.NewDecoder(r.Body).Decode(&blockData)
		blockchain.GetBlockchain().AddBlock(blockData.Data) // add to blockchain
		rw.WriteHeader(http.StatusCreated)                  // response 201
	}
}

func Start(portNum int) {
	// Ensure that diff. multiplexers (thing which calls handler funcs based on request url)
	// are used for web & rest API, so that there is no error regarding duplicate endpoints
	handler := http.NewServeMux()

	handler.HandleFunc("/", documentation)
	handler.HandleFunc("/blocks", blocks)

	port = fmt.Sprintf(":%d", portNum)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler)) // log when ListenAndServe returns an error
}
