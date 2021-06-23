package transport

// Server is transport server interface.
type Server interface {
	Start() error
	Stop() error
}
