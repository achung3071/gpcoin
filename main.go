package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/achung3071/gpcoin/api"
	"github.com/achung3071/gpcoin/webapp"
)

func instructions() {
	fmt.Printf("This is the GPCoin CLI.\n\n")
	fmt.Printf("Please use the following flags\n\n")
	fmt.Println("-mode:		Must be one of 'api', 'web'")
	fmt.Println("-port:		Set the port that the server should run on")
	os.Exit(0) // no error; Exit(1) is an eror
}

func main() {
	// automatically get flags from CLI and parse
	mode := flag.String("mode", "api", "Must be one of 'api', 'web'")
	port := flag.Int("port", 5000, "Set the port that the server should run on")
	flag.Parse()

	switch *mode {
	case "web":
		webapp.Start(*port)
	case "api":
		api.Start(*port)
	default:
		instructions()
	}
}
