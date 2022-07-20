package cronjob

// ServerOption is cron server option.
type ServerOption func(o *Server)

// WithAddress with server address.
func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.clientOpt.Addr = addr
	}
}
