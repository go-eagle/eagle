package sign

const (
	KeyNameTimeStamp = "time_stamp"
	KeyNameNonceStr  = "nonce_str"
	KeyNameAppId     = "app_id"
	KeyNameSign      = "sign"
)

var (
	gKeyNameTimestamp = KeyNameTimeStamp
	gKeyNameNonceStr  = KeyNameNonceStr
	gKeyNameAppId     = KeyNameAppId
	gKeyNameSign      = KeyNameSign
)

func SetKeyNameTimestamp(name string) {
	gKeyNameTimestamp = name
}

func SetKeyNameNonceStr(name string) {
	gKeyNameNonceStr = name
}

func SetKeyNameAppId(name string) {
	gKeyNameAppId = name
}

func SetKeyNameSign(name string) {
	gKeyNameSign = name
}

////

type DefaultKeyName struct {
	keyNameTimestamp string
	keyNameNonceStr  string
	keyNameAppId     string
	keyNameSign      string
}

func newDefaultKeyName() *DefaultKeyName {
	return &DefaultKeyName{
		keyNameTimestamp: gKeyNameTimestamp,
		keyNameNonceStr:  gKeyNameNonceStr,
		keyNameAppId:     gKeyNameAppId,
		keyNameSign:      gKeyNameSign,
	}
}

func (d *DefaultKeyName) SetKeyNameTimestamp(name string) {
	d.keyNameTimestamp = name
}

func (d *DefaultKeyName) SetKeyNameNonceStr(name string) {
	d.keyNameNonceStr = name
}

func (d *DefaultKeyName) SetKeyNameAppId(name string) {
	d.keyNameAppId = name
}

func (d *DefaultKeyName) SetKeyNameSign(name string) {
	d.keyNameSign = name
}
