package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achung3071/gpcoin/blockchain"
	"github.com/gorilla/mux"
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
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "Get a specific block",
			Payload:     "",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "Check status of blockchain",
			Payload:     "",
		},
	}
	json.NewEncoder(rw).Encode(urls) // easy way to send json to writer
}

type postBlocksBody struct {
	Data string
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var blockData postBlocksBody
		// save body in blockData variable (decoder automatically maps lowercase data -> Data)
		json.NewDecoder(r.Body).Decode(&blockData)
		blockchain.Blockchain().AddBlock(blockData.Data) // add to blockchain
		rw.WriteHeader(http.StatusCreated)               // response 201
	}
}

type errResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash) // get block based on height
	if err == blockchain.ErrBlockNotFound {
		rw.WriteHeader(404)
		json.NewEncoder(rw).Encode(errResponse{fmt.Sprint(err)})
	} else {
		json.NewEncoder(rw).Encode(block)
	}
}

// Send blockchain metadata
func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

// Attach application/json to every response
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	/* Normally, http.Handler is an interface having the ServeHTTP function.
	http.HandlerFunc is a TYPE that implements this interface. It defines
	the ServeHTTP function as simply calling the function f(rw, r) passed into
	it. Thus, it allows us to use a custom function as an http handler (i.e.,
	use custom middleware). As a result, we do not need to define our own struct
	and receiver function for http.Handler, and can use HandlerFunc(f) as an
	ADAPTER which allows us to easily define a Handler using our own function.
	(see source code for http.HandlerFunc for more info)
	*/
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(portNum int) {
	// Use mux from gorilla to specify a new multiplexer
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")

	port = fmt.Sprintf(":%d", portNum)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router)) // log when ListenAndServe returns an error
}
