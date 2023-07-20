package options

import "github.com/rabbitmq/amqp091-go"

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
