package stragollum

import "fmt"

// Collection represents a connection to a specific collection in the database.
type Collection struct {
	apiEndpoint string
	name        string
	token       *string
	keyspace    string
	commander   *DataAPICommander
}

// Keyspace returns the keyspace associated with the Database.
func (co *Collection) Name() string {
	return co.name
}

// InsertOne inserts a single document into the collection.
// It takes any Go type that can be marshalled to JSON as the document.
// Returns the inserted document's ID as a string and an error if the operation failed.
func (co *Collection) InsertOne(document interface{}) (string, error) {
	// Create the request payload as per API requirements
	// The payload structure should be: {"insertOne": {"document": <the input doc>}}
	requestPayload := struct {
		InsertOne struct {
			Document interface{} `json:"document"`
		} `json:"insertOne"`
	}{
		InsertOne: struct {
			Document interface{} `json:"document"`
		}{
			Document: document,
		},
	}

	// Define the expected response structure
	// The API returns: {"status": {"insertedIds": [<the id>]}}
	var response struct {
		Status struct {
			InsertedIds []string `json:"insertedIds"`
		} `json:"status"`
	}

	// Send the request and parse the response
	err := co.commander.Request(requestPayload, &response)
	if err != nil {
		return "", err
	}

	// Validate the response
	if len(response.Status.InsertedIds) == 0 {
		return "", fmt.Errorf("no document ID returned after insertion")
	}

	// Return the inserted document ID
	return response.Status.InsertedIds[0], nil
}
