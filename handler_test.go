package pie_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie"
)

// Ensure we can retrieve a list of tables through the HTTP interface.
func TestHandler_Tables(t *testing.T) {
	db := pie.NewDatabase()
	h := pie.NewHandler(db)
	w := httptest.NewRecorder()

	// Create tables.
	db.CreateTable("bob", nil)
	db.CreateTable("susy", nil)

	// Retrieve list of tables.
	r, _ := http.NewRequest("GET", "/tables", nil)
	h.ServeHTTP(w, r)

	// Verify the request was successful.
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w)
	} else if !strings.Contains(w.Body.String(), `<li><a href="/tables/bob">bob</a></li>`) {
		t.Fatalf("table 'bob' not found")
	} else if !strings.Contains(w.Body.String(), `<li><a href="/tables/susy">susy</a></li>`) {
		t.Fatalf("table 'susy' not found")
	}
}

// Ensure we can create a table through the HTTP interface.
func TestHandler_CreateTable(t *testing.T) {
	db := OpenDatabase()
	defer db.Close()
	h := pie.NewHandler(db.Database)
	s := httptest.NewServer(h)
	defer s.Close()

	// Generate multipart form body.
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", "names.csv")
	fmt.Fprint(part, "fname,lname!!\n")
	fmt.Fprint(part, "bob,smith\n")
	fmt.Fprint(part, "susy,que\n")
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// Upload a file.
	resp, _ := http.Post(s.URL+"/tables", w.FormDataContentType(), &buf)
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	} else if string(body) != "" {
		t.Fatalf("unexpected body: %s", body)
	}

	// Verify table is created.
	if tbl := db.Table("names"); tbl == nil {
		t.Fatal("expected table")
	} else if len(tbl.Columns) != 2 {
		t.Fatalf("expected column count: %d", len(tbl.Columns))
	} else if rows, _ := db.TableRows("names"); len(rows) != 2 {
		t.Fatalf("expected row count: %d", len(rows))
	}
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
