//go:build integration
// +build integration

package stragollum_test

import (
	"log"
	"os"
	"stragollum/pkg/stragollum"
	"testing"

	"github.com/joho/godotenv"
)

func TestListCollectionNames_Integration(t *testing.T) {
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

	// List collections
	collections, err := db.ListCollectionNames()
	if err != nil {
		t.Fatalf("ListCollectionNames failed: %v", err)
	}

	// Basic validation of the response
	if collections == nil {
		t.Error("Expected non-nil slice of collection names, got nil")
	}

	// Log the found collections (useful for debugging)
	log.Printf("Found %d collections: %v", len(collections), collections)

	// Validate that each collection name is a non-empty string
	for i, name := range collections {
		if name == "" {
			t.Errorf("Collection at index %d has empty name", i)
		}
	}
}
