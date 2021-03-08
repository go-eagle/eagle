package email

import (
	"errors"
	"sync"

	"github.com/1024casts/snake/config"

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
		Name:      config.Conf.Email.Name,
		Address:   config.Conf.Email.Address,
		ReplyTo:   config.Conf.Email.ReplyTo,
		Host:      config.Conf.Email.Host,
		Port:      config.Conf.Email.Port,
		Username:  config.Conf.Email.Username,
		Password:  config.Conf.Email.Password,
		Keepalive: config.Conf.Email.KeepAlive,
	})

	Client = client
}
