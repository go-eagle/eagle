package email

import (
	"testing"

	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
)

func TestSend(t *testing.T) {
	// init config
	cfg := "../../conf/config.sample.yaml"
	if _, err := conf.Init(cfg); err != nil {
		panic(err)
	}

	// init log
	log.InitLog(conf.Conf)

	Init()

	type args struct {
		to      string
		subject string
		body    string
	}

	subject, body := NewResetPasswordHTMLEmail("test", "http://snake.com")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test send reset mail", args{"test100@test.com", subject, body}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Send(tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
