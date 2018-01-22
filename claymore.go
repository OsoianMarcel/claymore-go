package claymore

// Service structure
type Client struct {
	Conn Connection
}

// Create new instance of client
func NewClient(serverAddress string) Client {
	conn := NewConnection(serverAddress)

	return Client{conn}
}

// Get claymore stats
func (c Client) GetStats() (StatsModel, error) {
	return NewStatsService(c.Conn).Execute()
}