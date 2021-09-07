# GPCoin

GPCoin is a basic cryptocurrency written in the Go language. It should be noted that this GPCoin does not have any
inherent value — it is simply used to illustrate the concept of blockchain currency. This is also why it is called
GPCoin, as GP stands for 거품 (guh-poom), which means "bubble" in Korean & is used to refer to something that looks
polished on the outside but has no actual value.

## Setup

### Downloading dependencies

Run `go mod download` in the root directory of this repository. Make sure to have Go version &ge; 1.16.

### Running the application

The most basic command is `go run main.go`. This will begin running the REST API for the blockchain on port 5000.

Here are the flags that can be used for the project:

- `-mode`: Can take on values `api` or `web`. `api` initializes the REST API endpoints for interacting with the
  cryptocurrency, while `web` hosts the HTML blockchain explorer web application (relevant files can be found in
  the [/webapp](webapp) folder).
- `-port`: Can take on any valid integer value for the port hosting the application. Default is `5000`.

Additionally, the `-race` flag can be used (e.g., `go run -race main.go -mode=api -port=4000`) to check for existing
data race conditions while the application is running.

## Usage

### Interacting with the blockchain

The genesis block in the blockchain is created when calling any of the API endpoints interacting with the blockchain
(e.g., `GET /blocks`). Other endpoints for actions such as mining a new block or adding a new transaction to the mempool
can be found in the HTTP response to `GET /` (see [api/endpoints.go](api/endpoints.go) for reference).

### P2P network

The code in this repository can be interpreted as the code for a single node in a P2P network. To simulate multiple peers,
one can run the following commands:

- Initialize nodes (run `go run main.go -port=xxxx`) on multiple ports (e.g., 2000, 3000, 4000, 5000).
- Ex. To make port 3000 send a websocket upgrade request to port 4000, send `POST /peers` to localhost:3000 with the body

      {address: "127.0.0.1", port: "4000"}.

  The nodes on port 3000 and 4000 will be peers.

- Ex. To make port 5000 be peers with 3000 and 4000, send `POST /peers` to localhost:3000 with the body

      {address: "127.0.0.1", port: "5000"}

  5000 will be added to the peers of 3000, and 3000 will broadcast 5000 to port 4000 so that the nodes are all conected.

- Once a P2P network is constructed, interactions with the blockchain (adding a transaction, mining a block, etc.) will be
  replicated across all peers in a synchronized manner.

### Running tests

Tests can be run by simply running the command `go test ./...`.

`go test -v -coverprofile=cover.out ./... && go tool cover -html=cover.out` will print all logs from the test cases, generate
a report of test coverage and display it in the browser as an HTML file.

## Remaining action items

- Refactor comments to give better API documentation in Godoc.
- Refactor and update web application to more widely interact with the blockchain.
- Make blockchain searching functions (e.g., UTxOutsByAddress) more performant.
- Handle errors better in a variety of ways based on situation (i.e., don't just log.Panic() on every error).
- Create a marshaler for checking whether HTTP request body data types are valid.
