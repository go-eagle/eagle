package sign

import (
	"fmt"
	"testing"
)

func TestVerify_ParseQuery(t *testing.T) {
	requestURI := "/restful/api/numbers?app_id=9d8a121ce581499d&nonce_str=tempstring&city=beijing" +
		"&timestamp=1532585241&sign=0f5b8c97920bc95f1a8b893f41b42d9e"

	// 第一步：创建Verify校验类
	verifier := NewVerifier()

	// 假定从RequestUri中读取校验参数
	if err := verifier.ParseQuery(requestURI); nil != err {
		t.Fatal(err)
	}

	// 第二步：（可选）校验是否包含签名校验必要的参数
	if err := verifier.MustHasOtherKeys("city"); nil != err {
		t.Fatal(err)
	}

	// 第三步：检查时间戳是否超时。
	//if err := verifier.CheckTimeStamp(); nil != err {
	//	t.Fatal(err)
	//}

	// 第四步，创建Sign来重现客户端的签名信息：
	signer := NewSignerMd5()

	// 第五步：从Verify中读取所有请求参数
	signer.SetBody(verifier.GetBodyWithoutSign())

	// 第六步：从数据库读取AppID对应的SecretKey
	// appId := verifier.MustString("app_id")
	secretKey := "d93047a4d6fe6111"

	// 使用同样的WrapBody方式
	signer.SetAppSecretWrapBody(secretKey)

	// 生成
	sign := signer.GetSignature()
	t.Log("sign", sign)

	// 校验自己生成的和传递过来的是否一致
	if verifier.MustString("sign") != sign {
		t.Fatal("校验失败")
	}

	fmt.Println(sign)
}
