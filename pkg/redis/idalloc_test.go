package redis

import (
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	InitTestRedis()

	type args struct {
		conn *redis.Client
	}
	tests := []struct {
		name string
		args args
		want *IDAlloc
	}{
		{
			name: "test new id alloc",
			args: args{
				conn: RedisClient,
			},
			want: &IDAlloc{
				redisClient: RedisClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIDAlloc(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdAlloc_GetCurrentID(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ia := &IDAlloc{
				redisClient: tt.fields.redisClient,
			}
			got, err := ia.GetCurrentID("user_id")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdAlloc_GetKey(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ia := &IDAlloc{
				redisClient: tt.fields.redisClient,
			}
			if got := ia.GetKey("user_id"); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// nolint
func TestIdAlloc_GetNewID(t *testing.T) {
	type fields struct {
		key         string
		redisClient *redis.Client
	}
	type args struct {
		step int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ia := &IDAlloc{
				redisClient: tt.fields.redisClient,
			}
			got, err := ia.GetNewID("user_id", tt.args.step)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNewID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNewID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
