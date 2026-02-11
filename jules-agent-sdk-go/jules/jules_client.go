package jules

// JulesClient is the main client for the Jules API
type JulesClient struct {
	baseClient *BaseClient
	Sessions   *SessionsAPI
	Activities *ActivitiesAPI
	Sources    *SourcesAPI
}

// NewClient creates a new JulesClient with the given API key
func NewClient(apiKey string) (*JulesClient, error) {
	config, err := NewConfig(apiKey)
	if err != nil {
		return nil, err
	}

	return NewClientWithConfig(config)
}

// NewClientWithConfig creates a new JulesClient with a custom configuration
func NewClientWithConfig(config *Config) (*JulesClient, error) {
	baseClient, err := NewBaseClient(config)
	if err != nil {
		return nil, err
	}

	return &JulesClient{
		baseClient: baseClient,
		Sessions:   NewSessionsAPI(baseClient),
		Activities: NewActivitiesAPI(baseClient),
		Sources:    NewSourcesAPI(baseClient),
	}, nil
}

// Close closes the client and releases all resources
func (c *JulesClient) Close() error {
	return c.baseClient.Close()
}

// Stats returns statistics about the client
func (c *JulesClient) Stats() map[string]int {
	return c.baseClient.Stats()
}
