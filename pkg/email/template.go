package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/log"
)

// NewActivationEmail 发送激活邮件
func NewActivationEmail(username, activateURL string) (subject string, body string) {
	return "帐号激活链接", "Hi, " + username + "<br>请激活您的帐号： <a href = '" + activateURL + "'>" + activateURL + "</a>"
}

// ActiveUserMailData 激活用户模板数据
type ActiveUserMailData struct {
	HomeURL       string `json:"home_url"`
	WebsiteName   string `json:"website_name"`
	WebsiteDomain string `json:"website_domain"`
	ActivateURL   string `json:"activate_url"`
	Year          int    `json:"year"`
}

// NewActivationHTMLEmail 发送激活邮件 html
func NewActivationHTMLEmail(username, activateURL string) (subject string, body string) {
	mailData := ActiveUserMailData{
		HomeURL:       viper.GetString("website.domain"),
		WebsiteName:   viper.GetString("website.name"),
		WebsiteDomain: viper.GetString("website.domain"),
		ActivateURL:   activateURL,
		Year:          time.Now().Year(),
	}
	mailTplContent := getEmailHTMLContent("./templates/active-mail.html", mailData)
	return "帐号激活链接", mailTplContent
}

// ResetPasswordMailData 激活用户模板数据
type ResetPasswordMailData struct {
	HomeURL       string `json:"home_url"`
	WebsiteName   string `json:"website_name"`
	WebsiteDomain string `json:"website_domain"`
	ResetURL      string `json:"reset_url"`
	Year          int    `json:"year"`
}

// NewResetPasswordEmail 发送重置密码邮件
func NewResetPasswordEmail(username, resetURL string) (subject string, body string) {
	return "密码重置", "Hi, " + username + "<br>您的重置链接为： <a href = '" + resetURL + "'>" + resetURL + "</a>"
}

// NewResetPasswordHTMLEmail 发送重置密码邮件 html
func NewResetPasswordHTMLEmail(username, resetURL string) (subject string, body string) {
	mailData := ResetPasswordMailData{
		HomeURL:       viper.GetString("website.domain"),
		WebsiteName:   viper.GetString("website.name"),
		WebsiteDomain: viper.GetString("website.domain"),
		ResetURL:      resetURL,
		Year:          time.Now().Year(),
	}
	mailTplContent := getEmailHTMLContent("./templates/reset-mail.html", mailData)
	return "密码重置", mailTplContent
}

// getEmailHTMLContent 获取邮件模板
func getEmailHTMLContent(tplPath string, mailData interface{}) string {
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		log.Warnf("[util.email] read file err: %v", err)
		return ""
	}
	mailTpl := string(b)
	tpl, err := template.New("email tpl").Parse(mailTpl)
	if err != nil {
		log.Warnf("[util.email] template new err: %v", err)
		return ""
	}
	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		fmt.Println("exec err", err)
		log.Warnf("[util.email] execute template err: %v", err)
	}
	return buffer.String()
}
