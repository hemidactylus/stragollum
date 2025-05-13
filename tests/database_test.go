package stragollum_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestGetDatabase(t *testing.T) {
	clientToken := "client_token"
	var testEnv stragollum.Environment = stragollum.EnvironmentProd // Or some default
	client := stragollum.NewDataAPIClient(&testEnv, &clientToken)
	apiEndpoint := "https://api.example.com"

	t.Run("AllParamsProvided", func(t *testing.T) {
		customToken := "custom_token"
		customKeyspace := "ks1"
		db := client.GetDatabase(apiEndpoint, &customToken, customKeyspace)
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Keyspace() != customKeyspace {
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), customKeyspace)
		}
		if db.Token() == nil || *db.Token() != customToken {
			t.Errorf("Token() = %v; want %v", db.Token(), &customToken)
		}
		if db.ApiEndpoint() != apiEndpoint {
			t.Errorf("ApiEndpoint() = %v; want %v", db.ApiEndpoint(), apiEndpoint)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, customKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != customToken {
			t.Errorf("Commander().Token() = %v; want %v", commander.Token(), &customToken)
		}
	})

	t.Run("TokenOmitted", func(t *testing.T) {
		customKeyspace := "ks2"
		db := client.GetDatabase(apiEndpoint, nil, customKeyspace)
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token() = %v; want %v (client token)", db.Token(), &clientToken)
		}
		if db.Keyspace() != customKeyspace {
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), customKeyspace)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, customKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != clientToken {
			t.Errorf("Commander().Token() = %v; want %v (client token)", commander.Token(), &clientToken)
		}
	})

	t.Run("KeyspaceOmitted", func(t *testing.T) {
		customToken := "custom_token2"
		db := client.GetDatabase(apiEndpoint, &customToken, "")
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Keyspace() != stragollum.DefaultKeyspace { // Use the actual constant
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), stragollum.DefaultKeyspace)
		}
		if db.Token() == nil || *db.Token() != customToken {
			t.Errorf("Token() = %v; want %v", db.Token(), &customToken)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, stragollum.DefaultKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != customToken {
			t.Errorf("Commander().Token() = %v; want %v", commander.Token(), &customToken)
		}
	})

	t.Run("TokenAndKeyspaceOmitted", func(t *testing.T) {
		db := client.GetDatabase(apiEndpoint, nil, "")
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token() = %v; want %v (client token)", db.Token(), &clientToken)
		}
		if db.Keyspace() != stragollum.DefaultKeyspace { // Use the actual constant
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), stragollum.DefaultKeyspace)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, stragollum.DefaultKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != clientToken {
			t.Errorf("Commander().Token() = %v; want %v (client token)", commander.Token(), &clientToken)
		}
	})
}

func TestListCollectionNames(t *testing.T) {
	// Setup a mock server to respond to the findCollections request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check request payload
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// The findCollections payload should be an empty object
		if _, exists := requestBody["findCollections"]; !exists {
			t.Errorf("Expected findCollections in request payload, got %v", requestBody)
		}

		// Send a mock response with collection names
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": {"collections": ["collection1", "collection2", "collection3"]}}`)
	}))
	defer server.Close()

	// Create a DataAPIClient to generate a Database
	clientToken := "client_token"
	var testEnv stragollum.Environment = stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&testEnv, &clientToken)

	// Get a Database instance pointing to our mock server
	db := client.GetDatabase(server.URL, nil, "test_keyspace")

	// Call the method under test
	collections, err := db.ListCollectionNames()

	// Verify results
	if err != nil {
		t.Fatalf("ListCollectionNames failed with error: %v", err)
	}

	// Check we got the expected collection names
	expected := []string{"collection1", "collection2", "collection3"}
	if len(collections) != len(expected) {
		t.Fatalf("Expected %d collections, got %d", len(expected), len(collections))
	}

	for i, name := range expected {
		if collections[i] != name {
			t.Errorf("Collection at index %d: expected %s, got %s", i, name, collections[i])
		}
	}

	// Test error handling
	// Setup a server that returns an error response
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error": "Database unavailable"}`)
	}))
	defer errorServer.Close()

	// Get a Database instance pointing to our error server
	errorDB := client.GetDatabase(errorServer.URL, nil, "test_keyspace")

	// Call the method under test
	_, err = errorDB.ListCollectionNames()

	// Verify we got an error
	if err == nil {
		t.Fatal("Expected error when server returns non-2xx status, got nil")
	}

	// Test invalid response format
	// Setup a server that returns an invalid response format
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "not_the_right_format"}`)
	}))
	defer invalidServer.Close()

	// Get a Database instance pointing to our invalid response server
	invalidDB := client.GetDatabase(invalidServer.URL, nil, "test_keyspace")

	// Call the method under test
	collections, err = invalidDB.ListCollectionNames()

	// The error might be nil since the JSON unmarshaling could succeed,
	// but collections should be empty since the response doesn't match the expected structure
	if len(collections) != 0 {
		t.Errorf("Expected empty collection list for invalid response, got %v", collections)
	}
}
