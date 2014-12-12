package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/turingschool-examples/pie"
)

const (
	// DefaultBindAddress is the default port that pie listens on.
	DefaultBindAddress = ":19876"
)

func main() {
	log.SetFlags(0)

	// Extract command from OS args.
	var cmd string
	var args []string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	// Set default command and update args list.
	if strings.HasPrefix(cmd, "-") || cmd == "" {
		cmd = "server"
		args = os.Args[1:]
	} else {
		args = os.Args[2:]
	}

	// Handle commands.
	switch cmd {
	case "server":
		runServer(args)
	case "exec", "execute":
		runExecute(args)
	default:
		log.Fatalf("invalid command: %s", cmd)
	}
}

func runServer(args []string) {
	// Parse command line flags.
	fs := flag.NewFlagSet("pie", flag.ExitOnError)
	dir := fs.String("d", "", "data directory")
	addr := fs.String("addr", DefaultBindAddress, "bind address")
	fs.Parse(args)

	// Set data directory to user directory if not set.
	if *dir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		*dir = filepath.Join(usr.HomeDir, ".pie")
	}

	// Open database.
	db := pie.NewDatabase()
	if err := db.Open(*dir); err != nil {
		log.Fatalf("open: %s", err)
	}
	defer db.Close()

	// Initialize handler.
	h := pie.NewHandler(db)

	// Start HTTP handler.
	log.Printf("Listening on http://localhost%s", *addr)
	log.SetFlags(log.LstdFlags)
	http.ListenAndServe(*addr, h)
}

func runExecute(args []string) {
	// Parse command line flags.
	fs := flag.NewFlagSet("pie", flag.ExitOnError)
	addr := fs.String("addr", DefaultBindAddress, "bind address")
	fs.Parse(args)

	// Read query string from arguments.
	str := strings.Join(fs.Args(), " ")

	// Execute POST against remote pie.
	u := fmt.Sprintf("http://localhost%s/query", *addr)
	resp, err := http.Post(u, "application/pieql", strings.NewReader(str))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Report non-200 status code.
	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stderr, resp.Body)
		os.Exit(-1)
	}

	// Write out response body.
	io.Copy(os.Stdout, resp.Body)
}
