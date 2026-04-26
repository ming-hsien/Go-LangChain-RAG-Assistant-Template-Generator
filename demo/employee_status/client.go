package employee_status

import (
	"math/rand"
	"time"
)

// Client simulates an API client for an Employee Status Server
type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// FetchStatus retrieves the real-time status of an employee from the "server"
func (c *Client) FetchStatus(name string) (string, error) {
	statuses := []string{
		"Available and ready for calls.",
		"In a focused work session (do not disturb).",
		"In a meeting until 16:00.",
		"Out for lunch, will be back in 30 minutes.",
		"Currently on annual leave (vacation).",
		"Working remotely, available via Slack/Email.",
		"In a customer support call.",
		"At an offline event, response may be delayed.",
	}

	// Simulation: logic moved here from the tool
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(len(name))))
	idx := r.Intn(len(statuses))

	return statuses[idx], nil
}
