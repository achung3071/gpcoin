package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/achung3071/gpcoin/blockchain"
	"github.com/achung3071/gpcoin/p2p"
	"github.com/achung3071/gpcoin/utils"
	"github.com/achung3071/gpcoin/wallet"
	"github.com/gorilla/mux"
)

var port string

type url string // custom type

// Response for /balance endpoint
type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// Request for /peers endpoint
type postPeersBody struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

// Request for /transactions endpoint
type postTransactionsBody struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type urlDescription struct {
	URL         url    `json:"url"` // struct field tag -> renames based on encoding
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omits this field when non-existent
}

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

// Another common interface: Stringer interface (implements String())
// "fmt" uses String() to display Go structs/instances when printed
/* func (u urlDescription) String() string {
	return "This is a url description"
} */

// HTTP HANDLER MIDDLEWARE
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

// Log endpoint each request is sent to
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request URL: %s\n", r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}

// HTTP HANDLER FUNCTIONS
// Get either TxOuts or total balance for given address/user
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	showTotal := r.URL.Query().Get("total")
	if showTotal == "true" {
		// Show total balance
		response := balanceResponse{address, blockchain.BalanceByAddress(address, blockchain.Blockchain())}
		utils.ErrorHandler(json.NewEncoder(rw).Encode(response))
	} else {
		// Show transaction outputs
		utils.ErrorHandler(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

// Get list of blocks (GET) | Mine a new block (POST)
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		block := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(block)
		rw.WriteHeader(http.StatusCreated)
	}
}

// Get a specific block based on the hash
func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	if err == blockchain.ErrBlockNotFound {
		rw.WriteHeader(404)
		json.NewEncoder(rw).Encode(errResponse{fmt.Sprint(err)})
	} else {
		json.NewEncoder(rw).Encode(block)
	}
}

// Check the current mempool
func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.ErrorHandler(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

// Get list of peers (GET) | Add a new peer via websocket (POST)
func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	case "POST":
		var data postPeersBody
		err := json.NewDecoder(r.Body).Decode(&data)
		utils.ErrorHandler(err)
		myPort := port[1:] // remove ":"
		// broadcast is true b/c peer added via API request (not broadcasted yet)
		p2p.AddPeer(data.Address, data.Port, myPort, true)
		rw.WriteHeader(http.StatusCreated)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Send blockchain metadata
func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

// Add a new transaction to mempool
func transactions(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var data postTransactionsBody
		json.NewDecoder(r.Body).Decode(&data) // get data
		// Add the new transaction to the blockchain mempool
		tx, err := blockchain.Mempool().AddTx(data.To, data.Amount)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(errResponse{err.Error()})
			return
		}
		p2p.BroadcastNewTx(tx)             // send new tx to all peers
		rw.WriteHeader(http.StatusCreated) // successfully created transaction
	}
}

// Returns address of wallet used by this node
func walletAddress(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(struct {
		Address string `json:"address"`
	}{Address: address})
}

func Start(portNum int) {
	// Use mux from gorilla to specify a new multiplexer
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)

	router.HandleFunc("/", Documentation).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/wallet-address", walletAddress).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")

	port = fmt.Sprintf(":%d", portNum)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router)) // log when ListenAndServe returns an error
}
