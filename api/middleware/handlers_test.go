package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetBooks(t *testing.T) {

	req, err := http.NewRequest("GET", "/book", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllBooks)
	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Different Status Code than Expected: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"title":"test","author":"test","publisher":"test","publishDate":"2020-01-01T00:00:00Z","rating":1,"status":"CheckedIn"},{"id":2,"title":"test2","author":"test2","publisher":"test2","publishDate":"2021-01-02T00:00:00Z","rating":3,"status":"CheckedOut"}]`
	if strings.TrimSpace(resp.Body.String()) != expected {
		t.Errorf("Wrong result returned: got %v want %v", resp.Body.String(), expected)
	}
}

func TestGetBooksByID(t *testing.T) {
	path := fmt.Sprintf("/book/%s", "2")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc("/book/{id}", GetBook)
	router.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Different Status Code than Expected: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":2,"title":"test2","author":"test2","publisher":"test2","publishDate":"2021-01-02T00:00:00Z","rating":3,"status":"CheckedOut"}`
	if strings.TrimSpace(resp.Body.String()) != expected {
		t.Errorf("Wrong result returned: got %v want %v", resp.Body.String(), expected)
	}
}

func TestGetBooksByInvalidID(t *testing.T) {
	path := fmt.Sprintf("/book/%s", "-1")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc("/book/{id}", GetBook)
	router.ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateBook(t *testing.T) {

	var jsonStr = []byte(`{"id":0,"title":"Dune","author":"Frank Herbert","publisher":"Dune Publisher","publishDate":"1965-01-01T00:00:00Z","rating":3,"status":"CheckedIn"}`)

	req, err := http.NewRequest("POST", "/newbook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(AddBook)
	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"id":3,"message":"Book was created seccessfully"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body.String(), expected)
	}
}
