package middleware

import (
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/app"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/sign"
)

// SignMd5Middleware md5 签名校验中间件
func SignMd5Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sn, err := verifySign(c)
		response := app.NewResponse()
		if err != nil {
			response.Error(c, errcode.ErrInternalServer)
			c.Abort()
			return
		}

		if sn != nil {
			response.Error(c, errcode.ErrSignParam)
			c.Abort()
			return
		}

		c.Next()
	}
}

// verifySign 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {
	requestURI := c.Request.RequestURI
	// 创建Verify校验器
	verifier := sign.NewVerifier()
	sn := verifier.GetSign()

	// 假定从RequestUri中读取校验参数
	if err := verifier.ParseQuery(requestURI); nil != err {
		return nil, err
	}

	// 检查时间戳是否超时。
	if err := verifier.CheckTimeStamp(); nil != err {
		return nil, errors.Errorf("%s error", sign.KeyNameTimeStamp)
	}

	// 验证签名
	localSign := genSign()
	if sn == "" || sn != localSign {
		return nil, errors.New(fmt.Sprintf("%s error", sign.KeyNameSign))
	}

	return nil, nil
}

// genSign 生成签名
func genSign() string {
	// todo: 读取配置
	signer := sign.NewSignerMd5()
	signer.SetAppID("123456")
	signer.SetTimeStamp(time.Now().Unix())
	signer.SetNonceStr("supertempstr")
	signer.SetAppSecretWrapBody("20200711")

	return signer.GetSignature()
}
