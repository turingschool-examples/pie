package pie

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/turingschool-examples/pie/assets"
	"github.com/turingschool-examples/pie/pieql"
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
	h.mux.HandleFunc("/tables", h.serveCreateTable).Methods("POST")
	h.mux.HandleFunc("/tables/{name}", h.serveTable).Methods("GET")
	h.mux.HandleFunc("/visualize", h.serveVisualize).Methods("GET")
	h.mux.HandleFunc("/query", h.serveQuery).Methods("POST")

	return h
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// serveIndex processes a request to the root page.
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	Index(w, h.db.Tables())
}

// serveAsset serves an asset file by name.
func (h *Handler) serveAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// If in development mode then read from file system.
	// Otherwise read from the bundled assets.
	var b []byte
	if os.Getenv("PIE_ENV") == "development" {
		b, _ = ioutil.ReadFile(filepath.Join("assets", filename))
	} else {
		b, _ = assets.Asset(filename)
	}

	// Return a 404 if the file doesn't exist.
	if b == nil {
		http.NotFound(w, r)
		return
	}

	// Set content type.
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))

	// Write asset contents.
	w.Write(b)
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

	// Retrieve rows.
	rows, err := h.db.TableRows(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the table.
	TableShow(w, t, rows)
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
	if ext := path.Ext(name); ext != "" {
		name = name[0 : len(name)-len(ext)]
	}

	// Import file as CSV.
	i := NewCSVImporter()
	if err := i.Import(h.db, name, csv.NewReader(f)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// serveVisualize renders the visualization template.
func (h *Handler) serveVisualize(w http.ResponseWriter, r *http.Request) {
	Visualize(w)
}

// serveQuery executes a query against the database.
func (h *Handler) serveQuery(w http.ResponseWriter, r *http.Request) {
	// Parse the statement.
	stmt, err := pieql.NewParser(r.Body).Parse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute the statement.
	res, err := h.db.Execute(stmt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build header from fields.
	hdr := make([]string, len(stmt.Fields))
	for i, f := range stmt.Fields {
		hdr[i] = f.Name
	}

	// Write the results.
	cw := csv.NewWriter(w)
	cw.Write(hdr)
	cw.WriteAll(res)
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
