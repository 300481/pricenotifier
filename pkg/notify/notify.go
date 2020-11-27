package notify

// Client is an interface for a notification client
type Client interface {
	Notify(message string) bool
}
