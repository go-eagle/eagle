// 主要是配合jwt来生成用户登录token

package app

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-eagle/eagle/pkg/log"
)

const (
	// DefaultJwtSecret is the insecure default JWT secret used in configuration templates.
	DefaultJwtSecret = "JWT_SECRET"
	// MinJwtSecretLength is the recommended minimum length for JWT HMAC secrets (256 bits / 32 bytes).
	MinJwtSecretLength = 32
)

var (
	// ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
)

// Payload is the data of the JSON web token.
type Payload struct {
	UserID uint64
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

// IsInsecureJwtSecret checks whether the secret matches default/weak values
// or is shorter than MinJwtSecretLength (32 characters).
func IsInsecureJwtSecret(secret string) bool {
	if len(secret) < MinJwtSecretLength {
		return true
	}
	trimmed := strings.TrimSpace(secret)
	if strings.EqualFold(trimmed, DefaultJwtSecret) || strings.EqualFold(trimmed, "secret") {
		return true
	}
	return false
}

// WarnIfInsecureJwtSecret checks the provided secret and emits a warning log if it is insecure.
// Returns true if the secret is considered insecure.
func WarnIfInsecureJwtSecret(secret string) bool {
	if !IsInsecureJwtSecret(secret) {
		return false
	}

	msg := "[SECURITY WARNING] Insecure JWT secret detected. The JWT secret matches the default or is shorter than 32 characters. Set JWT_SECRET or EAGLE_JWT_SECRET environment variable in production."
	if l := log.GetLogger(); l != nil {
		l.Warn(msg)
	} else {
		stdlog.Println(msg)
	}
	return true
}

// GetJwtSecret resolves the JWT secret, prioritizing environment variables over configuration defaults.
// Priority:
// 1. JWT_SECRET environment variable
// 2. EAGLE_JWT_SECRET environment variable
// 3. Conf.JwtSecret from application configuration
func GetJwtSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	if secret := os.Getenv("EAGLE_JWT_SECRET"); secret != "" {
		return secret
	}
	if Conf != nil && Conf.JwtSecret != "" {
		return Conf.JwtSecret
	}
	return ""
}

// InitJwt resolves the JWT secret, synchronizes Conf.JwtSecret if Conf is configured,
// and emits a visible warning log if an insecure secret is detected.
func InitJwt() string {
	secret := GetJwtSecret()
	if Conf != nil && secret != "" {
		Conf.JwtSecret = secret
	}
	WarnIfInsecureJwtSecret(secret)
	return secret
}

// Parse validates the token with the specified secret,
// and returns the payloads if the token was valid.
func Parse(tokenString string, secret string) (*Payload, error) {
	if secret == "" {
		secret = GetJwtSecret()
	}
	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return nil, err
	}
	// Read the token if it's valid.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payloads := &Payload{}
		payloads.UserID = uint64(claims["user_id"].(float64))
		return payloads, nil
	}

	// Other errors.
	return nil, err
}

// ParseRequest gets the token from the header and
// pass it to the Parse function to parses the token.
func ParseRequest(c *gin.Context) (*Payload, error) {
	header := c.Request.Header.Get("Authorization")

	// Load the jwt secret from env or config
	secret := GetJwtSecret()

	if len(header) == 0 {
		return &Payload{}, ErrMissingHeader
	}

	var t string
	// Parse the header to get the token part.
	_, err := fmt.Sscanf(header, "Bearer %s", &t)
	if err != nil {
		fmt.Printf("fmt.Sscanf err: %+v", err)
	}
	return Parse(t, secret)
}

// Sign signs the payload with the specified secret.
// The token content.
// iss: （Issuer）签发者
// iat: （Issued At）签发时间，用Unix时间戳表示
// exp: （Expiration Time）过期时间，用Unix时间戳表示
// aud: （Audience）接收该JWT的一方
// sub: （Subject）该JWT的主题
// nbf: （Not Before）不要早于这个时间
// jti: （JWT ID）用于标识JWT的唯一ID
func Sign(ctx context.Context, payload map[string]interface{}, secret string, timeout int64) (tokenString string, err error) {
	if secret == "" {
		secret = GetJwtSecret()
	}
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["nbf"] = now
	claims["iat"] = now
	claims["exp"] = now + timeout

	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(secret))

	return
}
