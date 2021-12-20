package repository

import (
	"context"
	"flag"
	"testing"

	"github.com/spf13/pflag"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/testing/lich"
)

var (
	d      *repository
	cfgDir = pflag.StringP("config", "c", "", "eagle config file path.")
)

func TestMain(m *testing.M) {
	pflag.Parse()

	*cfgDir = "../../config/config"

	flag.Set("f", "../../test/docker-compose.yaml")
	flag.Parse()

	c := config.New(*cfgDir)
	var cfg app.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	if err := lich.Setup(); err != nil {
		panic(err)
	}
	defer lich.Teardown()

	// init log
	logger.Init()
	// init db
	model.Init()

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
