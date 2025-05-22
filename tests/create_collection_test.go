package stragollum_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestDatabase_CreateCollection(t *testing.T) {
	// Expected payload for assertion
	expectedName := "my_collection"
	definition := stragollum.NewCollectionDefinition().
		WithDefaultID("string").
		WithLexical("standard").
		WithIndexing(map[string]any{"foo": "bar"})

	// Start a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		cc, ok := payload["createCollection"].(map[string]any)
		if !ok {
			t.Fatalf("Payload missing createCollection: %v", payload)
		}
		if cc["name"] != expectedName {
			t.Errorf("Expected name %q, got %v", expectedName, cc["name"])
		}
		// Check options is present and is a map
		if _, ok := cc["options"].(map[string]any); !ok {
			t.Errorf("Expected options to be a map, got %T", cc["options"])
		}
		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": {"ok": 1}}`)
	}))
	defer server.Close()

	// Setup client and database
	token := "dummy"
	env := stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&env, &token)
	db := client.GetDatabase(server.URL, nil, "ks1")

	// Should succeed
	collection, err := db.CreateCollection(expectedName, definition)
	if err != nil {
		t.Fatalf("CreateCollection failed: %v", err)
	}
	if collection == nil {
		t.Fatalf("Expected non-nil collection, got nil")
	}

	// Test error: server returns ok != 1
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": {"ok": 0}}`)
	}))
	defer failServer.Close()
	failDB := client.GetDatabase(failServer.URL, nil, "ks1")
	_, err = failDB.CreateCollection(expectedName, definition)
	if err == nil {
		t.Error("Expected error when status.ok != 1, got nil")
	}

	// Test error: server returns missing status
	missingStatusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"notstatus": {}}`)
	}))
	defer missingStatusServer.Close()
	missingStatusDB := client.GetDatabase(missingStatusServer.URL, nil, "ks1")
	_, err = missingStatusDB.CreateCollection(expectedName, definition)
	if err == nil {
		t.Error("Expected error when status is missing, got nil")
	}

	// Test error: server returns missing ok
	missingOkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": {}}`)
	}))
	defer missingOkServer.Close()
	missingOkDB := client.GetDatabase(missingOkServer.URL, nil, "ks1")
	_, err = missingOkDB.CreateCollection(expectedName, definition)
	if err == nil {
		t.Error("Expected error when status.ok is missing, got nil")
	}
}
