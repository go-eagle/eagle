package base

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	r := NewRepo("https://github.com/1024casts/snake-layout.git")
	if err := r.Clone(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "/tmp/test_snake_repo", "github.com/1024casts/snake-layout", nil); err != nil {
		t.Fatal(err)
	}
}
