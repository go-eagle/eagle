package transport

type Server interface {
	Start() error
	Stop() error
}
