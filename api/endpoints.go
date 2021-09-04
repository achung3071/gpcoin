package api

import (
	"encoding/json"
	"net/http"
)

func Documentation(rw http.ResponseWriter, r *http.Request) {
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
			Description: "Mine a block and add to blockchain",
			Payload:     "",
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
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get transaction outputs or balance(?total=true) at address",
			Payload:     "",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "Get the current mempool",
			Payload:     "",
		},
		{
			URL:         url("/transactions"),
			Method:      "POST",
			Description: "Post a new transaction to the mempool",
			Payload:     "{to: string, amount: int}",
		},
		{
			URL:         url("/wallet-address"),
			Method:      "GET",
			Description: "Get address of wallet used to post transactions",
			Payload:     "",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to websocket connection",
			Payload:     "",
		},
		{
			URL:         url("/peers"),
			Method:      "POST",
			Description: "Add a peer via websocket connection",
			Payload:     "{address: string, port: string}",
		},
	}
	json.NewEncoder(rw).Encode(urls) // easy way to send json to writer
}
