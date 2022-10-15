package utils

import (
	"errors"
	"testing"
)

func TestPrintStackTrace(t *testing.T) {
	t.Run("mock a error", func(t *testing.T) {
		PrintStackTrace("mock a error", errors.New("throw a error"))
	})

	t.Run("mock a recover", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				PrintStackTrace("mock a recover", r)
			}
		}()

		panic("throw a panic")
	})

}
