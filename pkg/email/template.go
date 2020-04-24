package email

// NewActivationEmail 发送激活邮件
func NewActivationEmail(username, activateURL string) (subject string, body string) {
	return "帐号激活链接", "Hi, " + username + "<br>请激活您的帐号： <a href = '" + activateURL + "'>" + activateURL + "</a>"
}
