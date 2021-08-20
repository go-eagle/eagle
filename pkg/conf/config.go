package conf

import (
	"time"

	"github.com/go-eagle/eagle/pkg/storage/orm"
	"github.com/go-eagle/eagle/pkg/storage/sql"

	"github.com/go-eagle/eagle/pkg/email"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-eagle/eagle/pkg/storage/mongodb"
	"github.com/go-eagle/eagle/pkg/trace"
)

// Config global config
// include common and biz config
type Config struct {
	// common
	App    AppConfig
	Http   ServerConfig
	Grpc   ServerConfig
	Web    WebConfig
	Cookie CookieConfig
	QiNiu  QiNiuConfig

	// component config
	Logger  log.Config
	ORM     orm.Config
	MySQL   *sql.Config
	Redis   redis.Config
	Email   email.Config
	Trace   trace.Config
	MongoDB mongodb.Config
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

// ServerConfig server config
type ServerConfig struct {
	Network      string
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
