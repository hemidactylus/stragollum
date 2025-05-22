package stragollum

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
