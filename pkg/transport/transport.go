package transport

import "context"

// Server is transport server interface.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
