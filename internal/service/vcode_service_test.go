package service

import (
	"testing"

	"github.com/go-eagle/eagle/internal/dao"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/conf"

	"github.com/spf13/pflag"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "eagle config file path.")
)

func Test_vcodeService_GenLoginVCode(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	*cfgFile = "../../../config/config.yaml"
	cfg, err := conf.Init(*cfgFile)
	if err != nil {
		panic(err)
	}
	s := New(cfg, dao.New(cfg, model.GetDB()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.VCode().GenLoginVCode(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenLoginVCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenLoginVCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
