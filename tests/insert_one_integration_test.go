//go:build integration
// +build integration

package stragollum_test

import (
	"os"
	"stragollum/pkg/stragollum"
	"testing"

	"github.com/joho/godotenv"
)

func TestInsertOne_Integration(t *testing.T) {
	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	apiEndpoint := os.Getenv("ASTRA_DB_API_ENDPOINT")
	if apiEndpoint == "" {
		t.Fatal("ASTRA_DB_API_ENDPOINT environment variable is required")
	}

	token := os.Getenv("ASTRA_DB_APPLICATION_TOKEN")
	if token == "" {
		t.Fatal("ASTRA_DB_APPLICATION_TOKEN environment variable is required")
	}

	keyspace := os.Getenv("ASTRA_DB_KEYSPACE")
	if keyspace == "" {
		t.Fatal("ASTRA_DB_KEYSPACE environment variable is required")
	}

	// Create a client with the token
	env := stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&env, &token)

	// Get a database instance
	db := client.GetDatabase(apiEndpoint, nil, keyspace)

	// Collection name for testing
	const testCollectionName = "insertone_integration_test"

	// Create collection
	definition := stragollum.NewCollectionDefinition()
	collection, err := db.CreateCollection(testCollectionName, definition)
	if err != nil {
		t.Fatalf("Failed to create test collection: %v", err)
	}
	t.Logf("Created collection: %s", testCollectionName)

	// Test document to insert
	testDocument := map[string]interface{}{
		"name":        "Integration Test User",
		"email":       "integration@example.com",
		"age":         35,
		"isActive":    true,
		"preferences": []string{"reading", "coding"},
		"address": map[string]interface{}{
			"city":    "Test City",
			"country": "Test Country",
		},
	}

	// Insert the document
	docID, err := collection.InsertOne(testDocument)
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}

	// Validate the document ID
	if docID == "" {
		t.Error("Expected non-empty document ID, got empty string")
	} else {
		t.Logf("Successfully inserted document with ID: %s", docID)
	}

	// Cleanup: We don't delete the collection as it might be used for future tests
	// and deleting/creating collections too often might hit rate limits
}
