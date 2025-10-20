package jules

// JulesClient is the main client for the Jules API.
type JulesClient struct {
	client     *Client
	Sessions   *SessionsService
	Activities *ActivitiesService
	Sources    *SourcesService
}

// NewJulesClient creates a new JulesClient.
func NewJulesClient(apiKey string) *JulesClient {
	baseClient := NewClient(apiKey)
	return &JulesClient{
		client:     baseClient,
		Sessions:   NewSessionsService(baseClient),
		Activities: NewActivitiesService(baseClient),
		Sources:    NewSourcesService(baseClient),
	}
}