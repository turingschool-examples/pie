package pie

import (
	"net/http"
)

// Handler represents the HTTP handler.
type Handler struct {
	db *Database
}

// NewHandler returns a new instance of Handler associated with a database.
func NewHandler(db *Database) *Handler {
	return &Handler{
		db: db,
	}
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented") // TODO
}
