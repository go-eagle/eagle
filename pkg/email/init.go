package email

import (
	"errors"
	"sync"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/1024casts/snake/pkg/log"
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
func Init() {
	log.Info("email init")
	Lock.Lock()
	defer Lock.Unlock()

	// 确保是已经关闭的
	if Client != nil {
		Client.Close()
	}

	client := NewSMTPClient(SMTPConfig{
		Name:      conf.Conf.Email.Name,
		Address:   conf.Conf.Email.Address,
		ReplyTo:   conf.Conf.Email.ReplyTo,
		Host:      conf.Conf.Email.Host,
		Port:      conf.Conf.Email.Port,
		Username:  conf.Conf.Email.Username,
		Password:  conf.Conf.Email.Password,
		Keepalive: conf.Conf.Email.KeepAlive,
	})

	Client = client
}
