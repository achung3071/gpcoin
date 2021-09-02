package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/achung3071/gpcoin/api"
	"github.com/achung3071/gpcoin/db"
	"github.com/achung3071/gpcoin/webapp"
)

func displayUsage() {
	fmt.Printf("This is the GPCoin CLI.\n\n")
	fmt.Printf("Please use the following flags\n\n")
	fmt.Println("-mode:		Must be one of 'api', 'web'")
	fmt.Println("-port:		Set the port that the server should run on")
	runtime.Goexit() // ensure deferred calls (db.Close) are honored even when exiting
}

func Start() {
	// automatically get flags from CLI and parse
	mode := flag.String("mode", "api", "Must be one of 'api', 'web'")
	port := flag.Int("port", 5000, "Set the port that the server should run on")
	flag.Parse()
	db.SetDBName(*port)

	switch *mode {
	case "web":
		webapp.Start(*port)
	case "api":
		api.Start(*port)
	default:
		displayUsage()
	}
}
