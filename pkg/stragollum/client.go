package stragollum

type AstraDBClient struct {
	environment string
	token       *string
}

func NewAstraDBClient(environment string, token string) *AstraDBClient {
	return &AstraDBClient{
		environment: environment,
		token:       &token,
	}
}

func (c *AstraDBClient) GetEnvironment() string {
	return c.environment
}
