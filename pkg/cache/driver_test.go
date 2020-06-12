package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	asserts := assert.New(t)

	asserts.NoError(Set("test-key", "123", 10*time.Second))
}
