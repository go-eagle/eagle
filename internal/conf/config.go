package conf

import (
	"time"

	"github.com/1024casts/snake/pkg/mongodb"

	"github.com/1024casts/snake/pkg/net/tracing"

	"github.com/1024casts/snake/pkg/email"

	"github.com/1024casts/snake/pkg/log"

	redis2 "github.com/1024casts/snake/pkg/redis"

	"github.com/1024casts/snake/pkg/database/orm"
)

// Config global config
// include common and biz config
type Config struct {
	// common
	App     AppConfig
	Logger  log.Config
	MySQL   orm.Config
	Redis   redis2.Config
	Email   email.Config
	Web     WebConfig
	Cookie  CookieConfig
	QiNiu   QiNiuConfig
	Jaeger  tracing.Config
	MongoDB mongodb.Config

	// here can add biz conf

}

// AppConfig app config
type AppConfig struct {
	Name              string
	Version           string
	Mode              string
	Port              string
	PprofPort         string
	URL               string
	JwtSecret         string
	JwtTimeout        int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

// WebConfig web config
type WebConfig struct {
	Name   string
	Domain string
	Secret string
	Static string
}

// CookieConfig cookie config
type CookieConfig struct {
	Name     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Domain   string
	Secret   string
}

// QiNiuConfig qiniu config
type QiNiuConfig struct {
	AccessKey   string
	SecretKey   string
	CdnURL      string
	SignatureID string
	TemplateID  string
}
