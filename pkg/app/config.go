package app

import (
	"time"
)

var (
	// Conf global app var
	Conf *Config
)

// Config global config
// nolint
type Config struct {
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
	EnableTrace       bool
	EnablePprof       bool
	HTTP              ServerConfig
	GRPC              ServerConfig
}

// ServerConfig server config.
type ServerConfig struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
