package rabbitmq

import (
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

type Config struct {
	AutoDeclare bool                       `yaml:"auto-declare"`
	Timeout     string                     `yaml:"timeout"`
	Connection  *options.ConnectionOptions `yaml:"connection"`
	Exchange    *options.ExchangeOptions   `yaml:"exchange"`
	Queue       *options.QueueOptions      `yaml:"queue"`
	Bind        *options.BindOptions       `yaml:"bind"`
}
