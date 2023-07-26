package options

import "github.com/rabbitmq/amqp091-go"

var DefaultConnectionOptions = ConnectionOptions{
	Scheme:   "amqp",
	Host:     "localhost",
	Port:     5672,
	Username: "guest",
	Password: "guest",
	Vhost:    "/",
}

type ConnectionOption func(*ConnectionOptions)

type ConnectionOptions struct {
	URI      string          `json:"uri"`
	Scheme   string          `json:"scheme"`
	Username string          `json:"username"`
	Password string          `json:"password"`
	Port     int             `json:"port"`
	Host     string          `json:"host"`
	Vhost    string          `json:"vhost"`
	Config   *amqp091.Config `json:"config"`
}

func NewConnectionOptions(opts ...ConnectionOption) (*ConnectionOptions, error) {
	options := DefaultConnectionOptions
	for _, opt := range opts {
		opt(&options)
	}

	var (
		uri = amqp091.URI{
			Scheme:   options.Scheme,
			Host:     options.Host,
			Port:     options.Port,
			Username: options.Username,
			Password: options.Password,
			Vhost:    options.Vhost,
		}
		err error
	)

	if options.URI != "" {
		uri, err = amqp091.ParseURI(options.URI)
		if err != nil {
			return nil, err
		}
		options.Scheme = uri.Scheme
		options.Host = uri.Host
		options.Port = uri.Port
		options.Vhost = uri.Vhost
		options.Username = uri.Username
		options.Password = uri.Password
	}

	options.URI = uri.String()

	return &options, nil
}

func WithURI(uri string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.URI = uri
	}
}
