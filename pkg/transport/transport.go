package transport

import (
	"context"
	"net/url"
)

// Server is transport server interface.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Endpoint is registry endpoint.
type Endpoint interface {
	Endpoint() (*url.URL, error)
}
