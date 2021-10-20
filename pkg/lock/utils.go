package lock

import (
	"github.com/google/uuid"
)

// genToken 生成token
func genToken() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
