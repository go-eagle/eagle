package conf

import (
	"time"

	"github.com/1024casts/snake/pkg/email"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/redis"
	"github.com/1024casts/snake/pkg/storage/mongodb"
	"github.com/1024casts/snake/pkg/storage/orm"
	"github.com/1024casts/snake/pkg/tracing"
)

// Config global config
// include common and biz config
type Config struct {
	// common
	App     AppConfig
	Http    ServerConfig
	Grpc    ServerConfig
	Logger  log.Config
	MySQL   orm.Config
	Redis   redis.Config
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
	PprofPort         string
	URL               string
	JwtSecret         string
	JwtTimeout        int
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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
