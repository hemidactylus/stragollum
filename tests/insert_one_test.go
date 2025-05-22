package stragollum_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestCollection_InsertOne(t *testing.T) {
	// Test document to insert
	testDocument := map[string]interface{}{
		"name":  "Test User",
		"email": "test@example.com",
		"age":   30,
	}

	expectedDocumentID := "doc123"

	// Start a mock server for successful insert
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify content type and method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read and parse the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify the payload structure
		insertOne, ok := payload["insertOne"].(map[string]interface{})
		if !ok {
			t.Fatalf("Payload missing insertOne: %v", payload)
		}

		// Verify document is present in the request
		if _, ok := insertOne["document"].(map[string]interface{}); !ok {
			t.Errorf("Expected document to be a map, got %T", insertOne["document"])
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": {"insertedIds": ["%s"]}}`, expectedDocumentID)
	}))
	defer successServer.Close()

	// Setup client, database, and collection for testing
	token := "dummy"
	env := stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&env, &token)
	db := client.GetDatabase(successServer.URL, nil, "ks1")
	collection := db.GetCollection("test_collection", nil)

	// Test successful insertion
	docID, err := collection.InsertOne(testDocument)
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}
	if docID != expectedDocumentID {
		t.Errorf("Expected document ID %s, got %s", expectedDocumentID, docID)
	}

	// Test error: no inserted IDs returned
	emptyIdsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": {"insertedIds": []}}`)
	}))
	defer emptyIdsServer.Close()

	emptyCollection := client.GetDatabase(emptyIdsServer.URL, nil, "ks1").GetCollection("test_collection", nil)
	_, err = emptyCollection.InsertOne(testDocument)
	if err == nil {
		t.Error("Expected error when no IDs returned, got nil")
	}

	// Test error: server returns error status
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "Invalid document"}`)
	}))
	defer errorServer.Close()

	errorCollection := client.GetDatabase(errorServer.URL, nil, "ks1").GetCollection("test_collection", nil)
	_, err = errorCollection.InsertOne(testDocument)
	if err == nil {
		t.Error("Expected error when server returns error status, got nil")
	}
}
