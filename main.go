package main

import (
	"fmt"
	"log"
	"net/http"
)

const port string = ":8080"

// basic handler for route
func home(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "This is the home page.") // Format & write to response writer
}

func main() {
	fmt.Printf("Listening on http://localhost%s\n", port)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(port, nil)) // log when ListenAndServe returns an error
}
