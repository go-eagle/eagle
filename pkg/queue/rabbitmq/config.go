package rabbitmq

import (
	"time"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

type Config struct {
	AutoDeclare bool                       `yaml:"auto-declare"`
	Timeout     time.Duration              `yaml:"timeout"`
	Connection  *options.ConnectionOptions `yaml:"connection"`
	Exchange    *options.ExchangeOptions   `yaml:"exchange"`
	Queue       *options.QueueOptions      `yaml:"queue"`
	Bind        *options.BindOptions       `yaml:"bind"`
}
