package commands

import (
	"flag"
	"log"

	"github.com/ersonp/go-rest-shell/pkg/api"
)

// Execute is the entry point for the command line interface
func Execute() {
	host := flag.String("host", "localhost", "Host to run the server on")
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()
	api := api.New(*host, *port, log.Default())
	api.Start()
}
