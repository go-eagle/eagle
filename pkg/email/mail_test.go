package email

import "testing"

func TestSend(t *testing.T) {
	type args struct {
		to      string
		subject string
		body    string
	}

	subject, body := NewActivationEmail("test", "snake.com")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test send mail", args{"go-snake@gmail.com", subject, body}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Send(tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
