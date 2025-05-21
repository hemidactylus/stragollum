package stragollum

import "fmt"

// Database represents a connection to a specific database/keyspace via the Data API.
type Database struct {
	apiEndpoint string
	token       *string
	keyspace    string
	commander   *DataAPICommander // Added commander field
}

// Keyspace returns the keyspace associated with the Database.
func (db *Database) Keyspace() string {
	return db.keyspace
}

// ApiEndpoint returns the API endpoint associated with the Database.
func (db *Database) ApiEndpoint() string {
	return db.apiEndpoint
}

// Token returns the token associated with the Database.
func (db *Database) Token() *string {
	return db.token
}

// Commander returns the DataAPICommander instance associated with the Database.
func (db *Database) Commander() *DataAPICommander {
	return db.commander
}

// ListCollectionNames retrieves the collection names in the database/keyspace.
// It returns a slice of strings containing the collection names, or an error if the request fails.
func (db *Database) ListCollectionNames() ([]string, error) {
	// Create the request payload as per API requirements
	requestPayload := struct {
		FindCollections struct{} `json:"findCollections"`
	}{
		FindCollections: struct{}{},
	}

	// Define the response structure based on the expected format
	responseData := struct {
		Status struct {
			Collections []string `json:"collections"`
		} `json:"status"`
	}{}

	// Send the request and parse the response
	err := db.commander.Request(requestPayload, &responseData)
	if err != nil {
		return nil, err
	}

	// Return the extracted collection names
	return responseData.Status.Collections, nil
}

// CreateCollection creates a new collection with the given name and definition (as options).
// Returns an error if the API response is not {"status": {"ok": 1}} or if the request fails.
func (db *Database) CreateCollection(name string, definition *CollectionDefinition) error {
	// Prepare the payload as per API spec
	type inner struct {
		Name    string                `json:"name"`
		Options *CollectionDefinition `json:"options"`
	}
	payload := struct {
		CreateCollection inner `json:"createCollection"`
	}{
		CreateCollection: inner{
			Name:    name,
			Options: definition,
		},
	}

	// Define the expected response structure
	var response struct {
		Status struct {
			Ok *int `json:"ok"`
		} `json:"status"`
	}

	err := db.commander.Request(payload, &response)
	if err != nil {
		return err
	}

	// Defensive: check for missing status or ok fields
	if response.Status.Ok == nil || *response.Status.Ok != 1 {
		return fmt.Errorf("unexpected response: expected status.ok == 1, got: %+v", response)
	}

	return nil
}
