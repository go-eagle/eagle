本包主要提供对API请求的签名生成、签名校验等

主要包含一下两部分:

## 生成签名

生成的签名需要满足以下几点：

- 可变性：每次的签名必须是不一样的
- 时效性：每次请求的时效性，过期作废
- 唯一性：每次的签名是唯一的
- 完整性：能够对传入数据进行验证，防止篡改

如果增加了签名验证，需要再传递几个参数：

- app_id 表示App Id，用来识别调用方身份
- timestamp 表示时间戳，用来验证接口的时效性
- sign 表示签名加密串，用来验证数据的完整性，防止数据篡改

支持两种签名生成算法：

- MD5: 由 NewSingerMd5() 提供
- Sha1 + Hmac:  由 NewSingerHmac() 提供

以上两种如果都不满足，可以自定义签名算法，使用 `NewSigner(FUNC)` 指定实现签名生成算法的实现即可

### Usage

```go
signer := NewSignerMd5()

// 设置签名基本参数
signer.SetAppId("94857djfi49484")
signer.SetTimeStamp(1594294833)
signer.SetNonceStr("xiKdApRhbuxVckJa")

// 设置参与签名的其它参数
signer.AddBody("plate_number", "golang")

// AppSecretKey，前后包装签名体字符串
signer.SetAppSecretWrapBody("x90449dfde34d")

fmt.Println("生成签字字符串：" + signer.GetUnsignedString())
fmt.Println("输出URL字符串：" + signer.GetSignedQuery())
```

### Result

## 校验签名

sign.Verifier 工具类，用来校验签名参数的格式和时间戳。它与Signer一起使用，用于服务端校验API请求的签名信息。

### Usage

```go
requestUri := "/restful/api/numbers?app_id=9d8a121ce581499d&nonce_str=ibuaiVcKdpRxkhJA&plate_number=豫A66666" +
		"&time_stamp=1532585241&sign=072defd1a251dc58e4d1799e17ffe7a4"

	// 第一步：创建Verifier校验类
	verifier := NewVerifier()

	// 假定从RequestUri中读取校验参数
	if err := verifier.ParseQuery(requestUri); nil != err {
		t.Fatal(err)
	}

	// 或者使用verifier.ParseValues(Values)来解析。

	// 第二步：（可选）校验是否包含签名校验必要的参数
	if err := verifier.MustHasOtherFields("plate_number"); nil != err {
		t.Fatal(err)
	}

	// 第三步：检查时间戳是否超时。

	// 时间戳超时：5分钟
	verifier.SetTimeout(time.Minute * 5)
	if err := verifier.CheckTimeStamp(); nil != err {
		t.Fatal(err)
	}

	// 第四步: 创建Signer来重现客户端的签名信息
	signer := NewSignerMd5()

	// 第五步：从Verifier中读取所有请求参数
	signer.SetBody(verifier.GetBodyWithoutSign())

	// 第六步：从数据库读取AppID对应的SecretKey
	// appId := verifier.GetAppId()
	secretKey := "123abc456"

	// 使用同样的WrapBody方式
	signer.SetAppSecretWrapBody(secretKey)

	// 服务端根据客户端参数生成签名
	sign := signer.GetSignature()

    // 最后，比较服务端生成的签名信息，与客户端提供的签名是否一致即可。
	if verifier.MustString("sign") != sign {
		t.Fatal("校验失败")
	}
```