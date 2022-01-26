package service

import (
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"

	"github.com/go-eagle/eagle/internal/repository"
)

// SMSService define interface func
type SMSService interface {
	SendSMS(phoneNumber string, verifyCode int) error
}

type smsService struct {
	repo repository.Repository
}

var _ SMSService = (*smsService)(nil)

func newSMS(svc *service) *smsService {
	return &smsService{repo: svc.repo}
}

// Send 发送短信
func (s *smsService) SendSMS(phoneNumber string, verifyCode int) error {
	// 校验参数的正确性
	if phoneNumber == "" || verifyCode == 0 {
		return errors.New("param phone or verify_code error")
	}

	// 调用第三方发送服务
	return sendViaQiNiu(phoneNumber, verifyCode)
}

// sendViaQiNiu 调用七牛短信服务
func sendViaQiNiu(phoneNumber string, verifyCode int) error {
	accessKey := ""
	secretKey := ""

	mac := auth.New(accessKey, secretKey)
	manager := sms.NewManager(mac)

	args := sms.MessagesRequest{
		SignatureID: "",
		TemplateID:  "",
		Mobiles:     []string{phoneNumber},
		Parameters: map[string]interface{}{
			"code": verifyCode,
		},
	}

	ret, err := manager.SendMessage(args)
	if err != nil {
		return errors.Wrap(err, "send sms message error")
	}

	if len(ret.JobID) == 0 {
		return errors.New("send sms message job id error")
	}

	return nil
}
