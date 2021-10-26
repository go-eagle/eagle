package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
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
	var appConfig AppConfig

	err := Load("../../config/app.yaml", &appConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(appConfig.Name)
}
