package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetJwtSecret(t *testing.T) {
	// Backup and restore global Conf
	originalConf := Conf
	defer func() {
		Conf = originalConf
	}()

	t.Run("prioritize JWT_SECRET over EAGLE_JWT_SECRET and Conf", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt_secret_from_primary_env_var_32chars")
		t.Setenv("EAGLE_JWT_SECRET", "eagle_secret_from_secondary_env_32chars")
		Conf = &Config{JwtSecret: "config_secret_from_yaml_file_32chars"}

		secret := GetJwtSecret()
		assert.Equal(t, "jwt_secret_from_primary_env_var_32chars", secret)
	})

	t.Run("use EAGLE_JWT_SECRET when JWT_SECRET is unset", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		t.Setenv("EAGLE_JWT_SECRET", "eagle_secret_from_secondary_env_32chars")
		Conf = &Config{JwtSecret: "config_secret_from_yaml_file_32chars"}

		secret := GetJwtSecret()
		assert.Equal(t, "eagle_secret_from_secondary_env_32chars", secret)
	})

	t.Run("fallback to Conf.JwtSecret when environment variables are unset", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		t.Setenv("EAGLE_JWT_SECRET", "")
		Conf = &Config{JwtSecret: "config_secret_from_yaml_file_32chars"}

		secret := GetJwtSecret()
		assert.Equal(t, "config_secret_from_yaml_file_32chars", secret)
	})

	t.Run("return empty string when Conf is nil and env vars are unset", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		t.Setenv("EAGLE_JWT_SECRET", "")
		Conf = nil

		secret := GetJwtSecret()
		assert.Equal(t, "", secret)
	})
}

func TestIsInsecureJwtSecret(t *testing.T) {
	testCases := []struct {
		name     string
		secret   string
		insecure bool
	}{
		{
			name:     "empty secret",
			secret:   "",
			insecure: true,
		},
		{
			name:     "default uppercase JWT_SECRET",
			secret:   DefaultJwtSecret,
			insecure: true,
		},
		{
			name:     "default lowercase jwt_secret",
			secret:   "jwt_secret",
			insecure: true,
		},
		{
			name:     "weak secret word",
			secret:   "secret",
			insecure: true,
		},
		{
			name:     "short secret less than 32 chars",
			secret:   "short_secret_key_12345",
			insecure: true,
		},
		{
			name:     "31 character secret",
			secret:   "1234567890123456789012345678901",
			insecure: true,
		},
		{
			name:     "valid 32 character secure secret",
			secret:   "12345678901234567890123456789012",
			insecure: false,
		},
		{
			name:     "valid strong secret exceeding 32 chars",
			secret:   "a_very_secure_and_long_jwt_secret_phrase_for_testing_purposes",
			insecure: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsInsecureJwtSecret(tc.secret)
			assert.Equal(t, tc.insecure, got)
		})
	}
}

func TestWarnIfInsecureJwtSecret(t *testing.T) {
	t.Run("warns on default secret", func(t *testing.T) {
		warned := WarnIfInsecureJwtSecret("JWT_SECRET")
		assert.True(t, warned)
	})

	t.Run("warns on short secret", func(t *testing.T) {
		warned := WarnIfInsecureJwtSecret("short_key")
		assert.True(t, warned)
	})

	t.Run("does not warn on secure secret", func(t *testing.T) {
		warned := WarnIfInsecureJwtSecret("secure_production_secret_key_that_has_at_least_32_bytes_len")
		assert.False(t, warned)
	})
}

func TestInitJwt(t *testing.T) {
	originalConf := Conf
	defer func() {
		Conf = originalConf
	}()

	t.Run("synchronizes Conf.JwtSecret and returns resolved secret", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "initialized_secret_key_32_characters_long!")
		Conf = &Config{JwtSecret: "old_secret"}

		secret := InitJwt()
		assert.Equal(t, "initialized_secret_key_32_characters_long!", secret)
		assert.Equal(t, "initialized_secret_key_32_characters_long!", Conf.JwtSecret)
	})

	t.Run("returns secret when Conf is nil", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "nil_conf_secret_key_32_characters_long!")
		Conf = nil

		secret := InitJwt()
		assert.Equal(t, "nil_conf_secret_key_32_characters_long!", secret)
	})
}

func TestSignAndParse(t *testing.T) {
	const secureKey = "production_grade_secret_key_with_sufficient_entropy_32"
	ctx := context.Background()

	t.Run("sign and parse successfully with explicit secret", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":  uint64(12345),
			"username": "testuser",
		}
		token, err := Sign(ctx, payload, secureKey, 3600)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)

		parsed, err := Parse(token, secureKey)
		assert.Nil(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, uint64(12345), parsed.UserID)
	})

	t.Run("sign and parse with environment variable fallback", func(t *testing.T) {
		t.Setenv("JWT_SECRET", secureKey)

		payload := map[string]interface{}{
			"user_id": uint64(67890),
		}
		// Passing empty secret triggers GetJwtSecret() fallback
		token, err := Sign(ctx, payload, "", 3600)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)

		parsed, err := Parse(token, "")
		assert.Nil(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, uint64(67890), parsed.UserID)
	})

	t.Run("parse with wrong secret fails", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id": uint64(111),
		}
		token, err := Sign(ctx, payload, secureKey, 3600)
		assert.Nil(t, err)

		_, err = Parse(token, "different_secret_key_that_is_also_32_chars_long")
		assert.NotNil(t, err)
	})

	t.Run("parse invalid token string fails", func(t *testing.T) {
		_, err := Parse("invalid.jwt.token", secureKey)
		assert.NotNil(t, err)
	})
}

func TestParseRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secureKey = "authorization_test_secret_key_at_least_32_chars"
	t.Setenv("JWT_SECRET", secureKey)

	t.Run("valid authorization header", func(t *testing.T) {
		token, err := Sign(context.Background(), map[string]interface{}{"user_id": uint64(999)}, secureKey, 3600)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		payload, err := ParseRequest(c)
		assert.Nil(t, err)
		assert.NotNil(t, payload)
		assert.Equal(t, uint64(999), payload.UserID)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		_, err := ParseRequest(c)
		assert.Equal(t, ErrMissingHeader, err)
	})

	t.Run("invalid token header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token-string")

		_, err := ParseRequest(c)
		assert.NotNil(t, err)
	})
}
