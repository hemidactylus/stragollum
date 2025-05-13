package stragollum

// Environment type for specifying deployment environment
type Environment string

// Constants for Environment type
const (
	EnvironmentDev  Environment = "dev"
	EnvironmentTest Environment = "test"
	EnvironmentProd Environment = "prod"
)

type DataAPIClient struct {
	environment Environment // Use value, not pointer
	token       *string     // token can be nil
}

// NewDataAPIClient creates a new DataAPIClient.
// If environment is nil, it defaults to EnvironmentProd.
// If token is nil, it remains nil.
func NewDataAPIClient(environment *Environment, token *string) *DataAPIClient {
	env := EnvironmentProd
	if environment != nil {
		env = *environment
	}
	return &DataAPIClient{
		environment: env,
		token:       token,
	}
}

// Environment returns the client's configured environment.
func (c *DataAPIClient) Environment() Environment {
	return c.environment
}

// Token returns the client's configured token.
func (c *DataAPIClient) Token() *string {
	return c.token
}

// GetDatabase creates a Database instance with the given apiEndpoint, optional token, and optional keyspace.
// If token is nil, uses the DataAPIClient's token. If keyspace is empty, uses the default DefaultKeyspace.
func (c *DataAPIClient) GetDatabase(apiEndpoint string, token *string, keyspace string) *Database {
	finalToken := token
	if finalToken == nil {
		finalToken = c.token
	}
	finalKeyspace := keyspace
	if finalKeyspace == "" {
		finalKeyspace = DefaultKeyspace
	}
	return &Database{
		apiEndpoint: apiEndpoint,
		token:       finalToken,
		keyspace:    finalKeyspace,
	}
}
