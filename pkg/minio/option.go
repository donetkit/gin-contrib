package minio

import (
	"context"
)

// Option for queue system
type Option func(*Client)

// WithContext set ctx function
func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.ctx = ctx
	}
}

// WithRegion set region function
func WithRegion(region string) Option {
	return func(c *Client) {
		c.region = region
	}
}

// WithEndpoint set endpoint function
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithAccess set accessKeyID, secretAccessKey function
func WithAccess(accessKeyID, secretAccessKey string) Option {
	return func(c *Client) {
		c.accessKeyID = accessKeyID
		c.secretAccessKey = secretAccessKey
	}
}

// WithUseSSL set useSSL function
func WithUseSSL(useSSL bool) Option {
	return func(c *Client) {
		c.useSSL = useSSL
	}
}

// WithEncryption set password function
func WithEncryption(password string) Option {
	return func(c *Client) {
		c.password = password
	}
}
