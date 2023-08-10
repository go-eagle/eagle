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
	var config AppConfig

	t.Run("using yaml config", func(t *testing.T) {

		c := New("./testdata")
		err := c.Load("app", &config)
		assert.Nil(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "eagle", config.Name)
	})
}
