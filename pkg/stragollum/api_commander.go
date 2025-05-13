package stragollum

// DataAPICommander is a helper for making HTTP POST requests to the Data API.
type DataAPICommander struct {
	url   string
	token *string
}

// NewDataAPICommander creates a new DataAPICommander with the given URL and optional token.
func NewDataAPICommander(url string, token *string) *DataAPICommander {
	return &DataAPICommander{
		url:   url,
		token: token,
	}
}

// URL returns the commander's URL.
func (c *DataAPICommander) URL() string {
	return c.url
}

// Token returns the commander's token (may be nil).
func (c *DataAPICommander) Token() *string {
	return c.token
}
