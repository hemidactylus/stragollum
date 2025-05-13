package stragollum_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestDataAPICommander(t *testing.T) {
	t.Run("WithToken", func(t *testing.T) {
		tok := "abc123"
		url := "https://api.example.com/api/json/v1/ks1"
		commander := stragollum.NewDataAPICommander(url, &tok)
		if commander.URL() != url {
			t.Errorf("URL() = %v; want %v", commander.URL(), url)
		}
		if commander.Token() == nil || *commander.Token() != tok {
			t.Errorf("Token() = %v; want %v", commander.Token(), tok)
		}
	})

	t.Run("NoToken", func(t *testing.T) {
		url := "https://api.example.com/api/json/v1/ks2"
		commander := stragollum.NewDataAPICommander(url, nil)
		if commander.URL() != url {
			t.Errorf("URL() = %v; want %v", commander.URL(), url)
		}
		if commander.Token() != nil {
			t.Errorf("Token() = %v; want nil", commander.Token())
		}
	})
}

func TestDataAPICommander_RawRequest(t *testing.T) {
	t.Run("SuccessfulRequestWithToken", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
			}
			expectedToken := "test-token"
			if r.Header.Get("Token") != expectedToken {
				t.Errorf("Expected Token '%s', got '%s'", expectedToken, r.Header.Get("Token"))
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"key":"value"}` {
				t.Errorf("Expected body '%s', got '%s'", `{"key":"value"}`, string(body))
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"status":"success"}`)
		}))
		defer server.Close()

		token := "test-token"
		commander := stragollum.NewDataAPICommander(server.URL, &token)
		payload := []byte(`{"key":"value"}`)

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		respBody, err := commander.RawRequest(payload, headers)

		if err != nil {
			t.Fatalf("RawRequest failed: %v", err)
		}
		if string(respBody) != `{"status":"success"}`+"\n" { // httptest server adds a newline
			t.Errorf("Expected response body '%s', got '%s'", `{"status":"success"}`+"\n", string(respBody))
		}
	})

	t.Run("SuccessfulRequestWithoutToken", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Token") != "" {
				t.Errorf("Expected no Token header, got '%s'", r.Header.Get("Token"))
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"status":"ok"}`)
		}))
		defer server.Close()

		commander := stragollum.NewDataAPICommander(server.URL, nil)
		payload := []byte(`{"data":"test"}`)

		// Empty headers
		headers := map[string]string{}

		respBody, err := commander.RawRequest(payload, headers)

		if err != nil {
			t.Fatalf("RawRequest failed: %v", err)
		}
		if string(respBody) != `{"status":"ok"}`+"\n" {
			t.Errorf("Expected response body '%s', got '%s'", `{"status":"ok"}`+"\n", string(respBody))
		}
	})

	t.Run("RequestFailedWithNon2xxStatus", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error":"bad request"}`)
		}))
		defer server.Close()

		commander := stragollum.NewDataAPICommander(server.URL, nil)
		payload := []byte(`{}`)
		_, err := commander.RawRequest(payload, nil)

		if err == nil {
			t.Fatal("RawRequest should have failed but did not")
		}
		expectedErrorMsg := "request failed with status 400: {\"error\":\"bad request\"}\n"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})

	t.Run("RequestFailsDueToServerShutdown", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Should not be called
		}))
		serverUrl := server.URL
		server.Close() // Close server immediately

		commander := stragollum.NewDataAPICommander(serverUrl, nil)
		payload := []byte(`{}`)
		_, err := commander.RawRequest(payload, nil)

		if err == nil {
			t.Fatal("RawRequest should have failed due to server being down, but it did not")
		}
		// The exact error message can vary depending on the OS and Go version for connection refused errors.
		// We check if it contains certain keywords.
		if !bytes.Contains([]byte(err.Error()), []byte("connect: connection refused")) &&
			!bytes.Contains([]byte(err.Error()), []byte("dial tcp")) &&
			!bytes.Contains([]byte(err.Error()), []byte("no such host")) { // macOS can give this for closed httptest server
			t.Errorf("Expected connection error, got: %v", err)
		}
	})
}

func TestDataAPICommander_Request(t *testing.T) {
	t.Run("SuccessfulJSONRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check headers
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
			}
			if r.Header.Get("Accept") != "application/json" {
				t.Errorf("Expected Accept: application/json, got %s", r.Header.Get("Accept"))
			}

			// Read and verify request body
			var requestBody map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}

			if requestBody["name"] != "test" || requestBody["value"] != float64(123) {
				t.Errorf("Unexpected request body: %v", requestBody)
			}

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"result":"success","count":42}`)
		}))
		defer server.Close()

		commander := stragollum.NewDataAPICommander(server.URL, nil)

		requestObj := struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}{
			Name:  "test",
			Value: 123,
		}

		responseObj := struct {
			Result string `json:"result"`
			Count  int    `json:"count"`
		}{}

		err := commander.Request(requestObj, &responseObj)

		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if responseObj.Result != "success" {
			t.Errorf("Expected result 'success', got '%s'", responseObj.Result)
		}
		if responseObj.Count != 42 {
			t.Errorf("Expected count 42, got %d", responseObj.Count)
		}
	})

	t.Run("InvalidRequestJSON", func(t *testing.T) {
		commander := stragollum.NewDataAPICommander("http://example.com", nil)

		// A function value cannot be marshaled to JSON
		requestObj := func() {}
		responseObj := struct{}{}

		err := commander.Request(requestObj, &responseObj)
		if err == nil {
			t.Fatal("Expected JSON marshal error, but got nil")
		}
	})

	t.Run("InvalidResponseJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"result": invalid json}`) // Invalid JSON
		}))
		defer server.Close()

		commander := stragollum.NewDataAPICommander(server.URL, nil)

		requestObj := struct{}{}
		responseObj := struct{}{}

		err := commander.Request(requestObj, &responseObj)
		if err == nil {
			t.Fatal("Expected JSON unmarshal error, but got nil")
		}
	})
}
