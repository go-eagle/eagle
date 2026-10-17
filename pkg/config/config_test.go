package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	t.Run("using yaml config default", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		t.Setenv("EAGLE_JWT_SECRET", "")
		c := New("./testdata")
		var config AppConfig
		err := c.Load("app", &config)
		assert.Nil(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "eagle", config.Name)
		assert.Equal(t, "JWT_SECRET", config.JwtSecret)
	})

	t.Run("jwt secret env override with JWT_SECRET", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "custom_secret_via_jwt_secret_env_var_32ch")
		t.Setenv("EAGLE_JWT_SECRET", "")
		c := New("./testdata")
		var cfg AppConfig
		err := c.Load("app", &cfg)
		assert.Nil(t, err)
		assert.Equal(t, "custom_secret_via_jwt_secret_env_var_32ch", cfg.JwtSecret)
	})

	t.Run("jwt secret env override with EAGLE_JWT_SECRET", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		t.Setenv("EAGLE_JWT_SECRET", "custom_secret_via_eagle_jwt_secret_32ch")
		c := New("./testdata")
		var cfg AppConfig
		err := c.Load("app", &cfg)
		assert.Nil(t, err)
		assert.Equal(t, "custom_secret_via_eagle_jwt_secret_32ch", cfg.JwtSecret)
	})

	t.Run("jwt secret env priority JWT_SECRET over EAGLE_JWT_SECRET", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "primary_jwt_secret_value_32chars_long")
		t.Setenv("EAGLE_JWT_SECRET", "secondary_eagle_jwt_secret_32chars_long")
		c := New("./testdata")
		var cfg AppConfig
		err := c.Load("app", &cfg)
		assert.Nil(t, err)
		assert.Equal(t, "primary_jwt_secret_value_32chars_long", cfg.JwtSecret)
	})

	t.Run("using LoadYaml with env override", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "load_yaml_custom_secret_32_characters_ok")
		t.Setenv("EAGLE_JWT_SECRET", "")
		c := New("./testdata")
		var cfg AppConfig
		err := c.LoadYaml("app", &cfg)
		assert.Nil(t, err)
		assert.Equal(t, "load_yaml_custom_secret_32_characters_ok", cfg.JwtSecret)
	})
}
