package vcode

import "testing"

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &vcodeService{}
			got, err := srv.GenLoginVCode(tt.args.phone)
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
