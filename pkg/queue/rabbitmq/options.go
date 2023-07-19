package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type Config amqp091.Config

type ConnectionOptions struct {
	Uri      string  `json:"uri"`
	Scheme   string  `json:"scheme"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Port     int     `json:"port"`
	Host     string  `json:"host"`
	Vhost    string  `json:"vhost"`
	Config   *Config `json:"config"`
}
