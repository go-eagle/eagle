package utils

import (
	"errors"
	"testing"
)

func TestPrintStackTrace(t *testing.T) {
	t.Run("mock a error", func(t *testing.T) {
		err := PrintStackTrace("mock a error", errors.New("throw a error"))
		t.Log(err)
	})

	t.Run("mock a recover", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err := PrintStackTrace("mock a recover", r)
				t.Log(err)
			}
		}()

		panic("throw a panic")
	})
}
