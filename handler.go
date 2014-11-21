package pie

import (
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/turingschool-examples/pie/assets"
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
	h.mux.HandleFunc("/", h.serveIndex).Methods("GET")
	h.mux.HandleFunc("/assets/{filename}", h.serveAsset).Methods("GET")
	h.mux.HandleFunc("/tables", h.serveTables).Methods("GET")
	h.mux.HandleFunc("/tables", h.serveCreateTable).Methods("POST")
	h.mux.HandleFunc("/tables/{name}", h.serveTable).Methods("GET")

	return h
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// serveIndex processes a request to the root page.
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	Index(w)
}

// serveAsset serves an asset file by name.
func (h *Handler) serveAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Read asset from assets package.
	filename := vars["filename"]
	b, _ := assets.Asset(filename)
	if b == nil {
		http.NotFound(w, r)
		return
	}

	// Set content type.
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))

	// Write asset contents.
	w.Write(b)
}

// serveTables processes a request to list tables in the database.
func (h *Handler) serveTables(w http.ResponseWriter, r *http.Request) {
	TableIndex(w, h.db.Tables())
}

// serveTable serves the contents of the table.
func (h *Handler) serveTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Find table and return error if it doesn't exist.
	t := h.db.Table(name)
	if t == nil {
		http.NotFound(w, r)
		return
	}

	// Render the table.
	TableShow(w, t)
}

// serveCreateTable processes a request to create a table in the database.
func (h *Handler) serveCreateTable(w http.ResponseWriter, r *http.Request) {
	// Check for file in request body.
	f, hdr, err := r.FormFile("file")
	if err != nil {
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
