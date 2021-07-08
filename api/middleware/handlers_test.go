package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

	expected := `[{"id":1,"title":"test","author":"test","publisher":"test","publishDate":"2020-01-01T00:00:00Z","rating":1,"status":"CheckedIn"}]`
	if resp.Body.String() != expected {
		t.Errorf("Wrong result returned: got %v want %v", resp.Body.String(), expected)
	}
}
