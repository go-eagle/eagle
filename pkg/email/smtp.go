package email

import (
	"time"

	"github.com/go-mail/mail"

	"github.com/1024casts/snake/pkg/log"
)

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Name      string // 发送者名称
	Address   string // 发送者地址
	ReplyTo   string // 回复地址
	Host      string // 服务器主机名
	Port      int    // 服务器端口
	Username  string // 用户名
	Password  string // 密码
	Keepalive int    // 连接保活时长, 单位秒
}

// SMTP 协议协议发送
type SMTP struct {
	Config SMTPConfig
	ch     chan *mail.Message
	chOpen bool
}

// NewSMTPClient 实例化一个SMTP客户端
func NewSMTPClient(config SMTPConfig) *SMTP {
	client := &SMTP{
		Config: config,
		ch:     make(chan *mail.Message, 30),
		chOpen: false,
	}

	return client
}

// Init 初始化发送队列
func (c *SMTP) Init() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.chOpen = false
				log.Error("send email queue err: %+v, retry after 10 second", err.(error))
				time.Sleep(time.Duration(10) * time.Second)
				c.Init()
			}
		}()

		d := mail.NewDialer(c.Config.Host, c.Config.Port, c.Config.Username, c.Config.Password)
		d.Timeout = time.Duration(c.Config.Keepalive+5) * time.Second
		c.chOpen = true

		var s mail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-c.ch:
				if !ok {
					log.Info("mail queue is close")
					c.chOpen = false
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := mail.Send(s, m); err != nil {
					log.Warnf("email send failed, %v", err)
				} else {
					log.Info("email has send")
				}
			case <-time.After(time.Duration(c.Config.Keepalive) * time.Second):
				if open {
					if err := s.Close(); err != nil {
						log.Warnf("can not close smtp conn, %v", err)
					}
				}
				open = false
			}
		}
	}()
}

// Send 发送邮件
func (c *SMTP) Send(to, subject, body string) error {
	if !c.chOpen {
		return ErrChanNotOpen
	}

	msg := mail.NewMessage()
	msg.SetAddressHeader("From", c.Config.Address, c.Config.Name)
	msg.SetAddressHeader("Reply-To", c.Config.ReplyTo, c.Config.Name)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("To", to)
	msg.SetBody("text/html", body)

	c.ch <- msg
	return nil
}

// Close 关闭队列
func (c *SMTP) Close() {
	if c.ch != nil {
		close(c.ch)
	}
}
