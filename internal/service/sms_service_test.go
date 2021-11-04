package service

import (
	"testing"
)

func Test_smsService_Send(t *testing.T) {
	type args struct {
		phoneNumber string
		verifyCode  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SMS().SendSMS(tt.args.phoneNumber, tt.args.verifyCode); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
