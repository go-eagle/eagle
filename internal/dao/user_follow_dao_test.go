package dao

import (
	"context"
	"flag"
	"testing"

	logger "github.com/1024casts/snake/pkg/log"

	"github.com/1024casts/snake/internal/model"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/pkg/testing/lich"
	"github.com/spf13/pflag"
)

var (
	d       *Dao
	cfgFile = pflag.StringP("config", "c", "", "snake config file path.")
	cfg     *conf.Config
)

func TestMain(m *testing.M) {
	pflag.Parse()

	*cfgFile = "../../config/config.yaml"

	flag.Set("f", "../../test/docker-compose.yaml")
	flag.Parse()

	cfg, err := conf.Init(*cfgFile)
	if err != nil {
		panic(err)
	}
	if err := lich.Setup(); err != nil {
		panic(err)
	}
	defer lich.Teardown()

	// init log
	logger.Init(&cfg.Logger)
	// init db
	model.Init(&cfg.MySQL)

	d = New(model.GetDB())
	if code := m.Run(); code != 0 {
		panic(code)
	}
}

func TestDao_GetFollowerUserList(t *testing.T) {
	followers, err := d.GetFollowerUserList(context.Background(), 1, 0, 1)
	if err != nil {
		t.Error(err)
	}
	if len(followers) > 0 {
		t.Log("follower num is: ", len(followers))
	}
}
