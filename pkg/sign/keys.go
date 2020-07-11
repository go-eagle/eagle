package sign

const (
	KeyNameTimeStamp = "timestamp"
	KeyNameNonceStr  = "nonce_str"
	KeyNameAppId     = "app_id"
	KeyNameSign      = "sign"
)

// DefaultKeyName 签名需要用到的字段
type DefaultKeyName struct {
	Timestamp string
	NonceStr  string
	AppId     string
	Sign      string
}

func newDefaultKeyName() *DefaultKeyName {
	return &DefaultKeyName{
		Timestamp: KeyNameTimeStamp,
		NonceStr:  KeyNameNonceStr,
		AppId:     KeyNameAppId,
		Sign:      KeyNameSign,
	}
}

func (d *DefaultKeyName) SetKeyNameTimestamp(name string) {
	d.Timestamp = name
}

func (d *DefaultKeyName) SetKeyNameNonceStr(name string) {
	d.NonceStr = name
}

func (d *DefaultKeyName) SetKeyNameAppId(name string) {
	d.AppId = name
}

func (d *DefaultKeyName) SetKeyNameSign(name string) {
	d.Sign = name
}
