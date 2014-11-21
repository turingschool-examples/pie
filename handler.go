package pie

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Handler represents the HTTP handler.
type Handler struct {
	db  *Database
	mux *mux.Router
}

// NewHandler returns a new instance of Handler associated with a database.
func NewHandler(db *Database) *Handler {
	// Initialize handler.
	h := &Handler{
		db:  db,
		mux: mux.NewRouter(),
	}

	// Setup request multiplexer.
	h.mux.HandleFunc("/tables", h.serveTables).Methods("GET")
	h.mux.HandleFunc("/tables", h.serveCreateTable).Methods("POST")

	return h
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// serveTables processes a request to list tables in the database.
func (h *Handler) serveTables(w http.ResponseWriter, r *http.Request) {
	for _, t := range h.db.tables {
		fmt.Fprintln(w, t.Name)
	}
}

// serveCreateTable processes a request to create a table in the database.
func (h *Handler) serveCreateTable(w http.ResponseWriter, r *http.Request) {
	// Check for file in request body.
	f, hdr, err := r.FormFile("file")
	if err != nil {
		warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	// Extract the filename.
	name := hdr.Filename

	// Import file as CSV.
	i := NewCSVImporter()
	if err := i.Import(h.db, name, csv.NewReader(f)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
