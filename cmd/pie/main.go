package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/turingschool-examples/pie"
)

func main() {
	// Parse command line flags.
	addr := flag.String("addr", ":19876", "bind address")
	flag.Parse()
	log.SetFlags(0)

	// Open database.
	db := pie.NewDatabase()

	// Initialize handler.
	h := pie.NewHandler(db)

	// Start HTTP handler.
	log.Printf("Listening on http://localhost%s", *addr)
	log.SetFlags(log.LstdFlags)
	http.ListenAndServe(*addr, h)
}
