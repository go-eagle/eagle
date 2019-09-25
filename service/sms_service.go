package service

import (
	"os"

	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"
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

	accessKey := os.Getenv("accessKey")
	secretKey := os.Getenv("secretKey")

	mac := auth.New(accessKey, secretKey)
	manager := sms.NewManager(mac)

	// CreateTemplate
	tempArgs := sms.TemplateRequest{
		Name:     "Test",
		Type:     sms.VerificationType,
		Template: "验证码: ${code}， 10分钟内有效。",
	}
	tempRet, err := manager.CreateTemplate(tempArgs)
	if err != nil {
		return errors.New("create template error")
	}
	if len(tempRet.TemplateID) == 0 {
		return errors.New("template id is empty")
	}

	// CreateSignature
	signArgs := sms.SignatureRequest{
		Signature: "Test",
		Source:    sms.APP,
	}
	signRet, err := manager.CreateSignature(signArgs)
	if err != nil {
		return errors.New("create signature error")
	}
	if len(signRet.SignatureID) == 0 {
		return errors.New("signature id is empty")
	}

	args := sms.MessagesRequest{
		SignatureID: signRet.SignatureID,
		TemplateID:  tempRet.TemplateID,
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
