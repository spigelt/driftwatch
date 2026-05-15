package docker

import (
	"context"

	"github.com/docker/docker/client"
)

// Client wraps the Docker SDK client.
type Client struct {
	c *client.Client
}

// NewClient creates a new Docker client using environment variables
// (DOCKER_HOST, DOCKER_TLS_VERIFY, etc.).
func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &Client{c: c}, nil
}

// Close releases the underlying Docker client resources.
func (dc *Client) Close() error {
	return dc.c.Close()
}

// Ping verifies the daemon is reachable.
func (dc *Client) Ping(ctx context.Context) error {
	_, err := dc.c.Ping(ctx)
	return err
}
