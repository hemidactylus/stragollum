package stragollum

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
