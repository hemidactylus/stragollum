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

func TestCreateCollection_Integration(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

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
	env := stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&env, &token)
	db := client.GetDatabase(apiEndpoint, nil, keyspace)

	const testCollectionName = "coll_create_test"
	const testRichCollectionName = "coll_create_test_rich"

	// 1. List collections before creation
	collections, err := db.ListCollectionNames()
	if err != nil {
		t.Fatalf("ListCollectionNames (before) failed: %v", err)
	}
	for _, name := range collections {
		if name == testCollectionName {
			t.Fatalf("Test collection %q already exists before creation!", testCollectionName)
		}
	}

	// 2. Create collection with minimal (empty) definition
	coll, err := db.CreateCollection(testCollectionName, stragollum.NewCollectionDefinition())
	if err != nil {
		t.Fatalf("CreateCollection failed: %v", err)
	}
	if coll.Name() != testCollectionName {
		t.Errorf("Created collection has wrong name: %q, expected %q", coll.Name(), testCollectionName)
	}

	// 2b. Create collection with a rich definition
	richDefinition := stragollum.NewCollectionDefinition().
		WithDefaultID("objectId").
		WithLexical("standard").
		WithIndexing(map[string]any{"allow": []string{"the_field"}}).
		WithRerank(
			&stragollum.RerankServiceOptions{
				Provider:  "nvidia",
				ModelName: "nvidia/llama-3.2-nv-rerankqa-1b-v2",
			},
			true,
		).
		WithVector(&stragollum.CollectionVectorOptions{
			Service: &stragollum.VectorServiceOptions{
				Provider:  "nvidia",
				ModelName: "NV-Embed-QA",
			},
		})
	richColl, err2 := db.CreateCollection(testRichCollectionName, richDefinition)
	if err2 != nil {
		t.Fatalf("Rich CreateCollection failed: %v", err2)
	}
	if richColl.Name() != testRichCollectionName {
		t.Errorf("Created rich collection has wrong name: %q, expected %q", richColl.Name(), testRichCollectionName)
	}

	// 3. List collections after creation
	collections, err3 := db.ListCollectionNames()
	if err3 != nil {
		t.Fatalf("ListCollectionNames (after) failed: %v", err3)
	}
	found := false
	foundRich := false
	for _, name := range collections {
		if name == testCollectionName {
			found = true
		}
		if name == testRichCollectionName {
			foundRich = true
		}
	}
	if !found {
		t.Errorf("Collection %q was not found after creation", testCollectionName)
	}
	if !foundRich {
		t.Errorf("Collection %q was not found after creation", testRichCollectionName)
	}

	// 4. Clean up: delete the test collection
	err = db.DropCollection(testCollectionName)
	if err != nil {
		t.Fatalf("DropCollection failed: %v", err)
	}
	// 4b. Clean up: delete the rich test collection
	err = db.DropCollection(testRichCollectionName)
	if err != nil {
		t.Fatalf("DropCollection failed: %v", err)
	}

	// 5. Verify deletion
	collections, err = db.ListCollectionNames()
	if err != nil {
		t.Fatalf("ListCollectionNames (after deletion) failed: %v", err)
	}
	found = false
	foundRich = false
	for _, name := range collections {
		if name == testCollectionName {
			found = true
		}
		if name == testRichCollectionName {
			foundRich = true
		}
	}
	if found {
		t.Errorf("Collection %q was found after deletion", testCollectionName)
	}
	if foundRich {
		t.Errorf("Collection %q was found after deletion", testRichCollectionName)
	}
}
