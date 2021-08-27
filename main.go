package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/achung3071/gpcoin/api"
)

func instructions() {
	fmt.Printf("This is the GPCoin CLI.\n\n")
	fmt.Printf("Use one of the following commands:\n\n")
	fmt.Println("web:	Run the HTML web application")
	fmt.Println("api:	Run the REST API server")
	os.Exit(0) // no error; Exit(1) is an eror
}

func main() {
	if len(os.Args) < 2 {
		instructions()
	}

	// can specify set of flags & how to handle flag errors
	apiFlagSet := flag.NewFlagSet("api", flag.ExitOnError)
	portFlag := apiFlagSet.Int("port", 5000, "Sets the port of the REST API server")

	switch os.Args[1] {
	case "web":
		fmt.Println("Starting web app...")
	case "api":
		apiFlagSet.Parse(os.Args[2:]) // parse flags from 2nd arg onwards
	default:
		instructions()
	}

	if apiFlagSet.Parsed() {
		api.Start(*portFlag)
	}
}
