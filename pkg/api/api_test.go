package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCmdHandler(t *testing.T) {
	api := New("localhost", 8080, log.Default())
	requestBody := bytes.NewBuffer([]byte(`{"command":"echo hello"}`))
	req, err := http.NewRequest("POST", "/api/cmd", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test server with the same handlers as your actual server
	server := httptest.NewServer(api.Mux)
	defer server.Close()

	// Update the request URL to the test server's URL
	req.URL, err = url.Parse(server.URL + "/api/cmd")
	if err != nil {
		t.Fatal(err)
	}
	// Make the request to the test server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body
	expected := strings.TrimSpace("hello")
	actual := strings.TrimSpace(string(body))
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}
