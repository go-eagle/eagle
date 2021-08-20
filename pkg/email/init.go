package email

import (
	"errors"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
)

// Client 邮件发送客户端
var Client Driver

// Lock 读写锁
var Lock sync.RWMutex

var (
	// ErrChanNotOpen 邮件队列没有开启
	ErrChanNotOpen = errors.New("email queue does not open")
)

// Config email config
type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Name      string
	Address   string
	ReplyTo   string
	KeepAlive int
}

// Init 初始化客户端
func Init(cfg Config) {
	log.Info("email init")
	Lock.Lock()
	defer Lock.Unlock()

	// 确保是已经关闭的
	if Client != nil {
		Client.Close()
	}

	client := NewSMTPClient(SMTPConfig{
		Name:      cfg.Name,
		Address:   cfg.Address,
		ReplyTo:   cfg.ReplyTo,
		Host:      cfg.Host,
		Port:      cfg.Port,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Keepalive: cfg.KeepAlive,
	})

	Client = client
}
