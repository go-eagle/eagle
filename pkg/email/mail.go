package email

// Driver 邮件发送驱动接口定义
type Driver interface {
	// Send 发送邮件
	Send(to, subject, body string) error
	// Close 关闭链接
	Close()
}

// Send 发送邮件
func Send(to, subject, body string) error {
	Lock.RLock()
	defer Lock.RUnlock()

	if Client == nil {
		return nil
	}

	return Client.Send(to, subject, body)
}
