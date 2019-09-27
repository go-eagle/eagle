package service

import (
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"
	"github.com/spf13/viper"
)

// 直接初始化，可以避免在使用时再实例化
// 使用七牛云短信服务
var SmsService = NewSmsService()

// 校验码服务，生成校验码和获得校验码
type smsService struct {
}

func NewSmsService() *smsService {
	return &smsService{}
}

// 发送短信
func (srv *smsService) Send(phoneNumber string, verifyCode int) error {
	// 校验参数的正确性
	if phoneNumber == "" || verifyCode == 0 {
		return errors.New("param error")
	}

	// 调用第三方发送服务
	return srv._sendViaQiNiu(phoneNumber, verifyCode)
}

// 调用第三方发送服务
func (srv *smsService) _sendViaQiNiu(phoneNumber string, verifyCode int) error {

	accessKey := viper.GetString("qiniu.access_key")
	secretKey := viper.GetString("qiniu.secret_key")

	mac := auth.New(accessKey, secretKey)
	manager := sms.NewManager(mac)

	args := sms.MessagesRequest{
		SignatureID: viper.GetString("qiniu.signature_id"),
		TemplateID:  viper.GetString("qiniu.template_id"),
		Mobiles:     []string{phoneNumber},
		Parameters: map[string]interface{}{
			"code": verifyCode,
		},
	}

	ret, err := manager.SendMessage(args)
	if err != nil {
		return errors.New("send sms message error")
	}

	if len(ret.JobID) == 0 {
		return errors.New("send sms message job id error")
	}

	return nil
}
