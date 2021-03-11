package mongodb

import (
	"context"
	"time"

	"github.com/1024casts/snake/pkg/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

// NewMongoDBConn Create new MongoDB client
func NewMongoDBConn(ctx context.Context, cfg *conf.Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(cfg.MongoDB.URI).
			SetAuth(options.Credential{
				Username: cfg.MongoDB.User,
				Password: cfg.MongoDB.Password,
			}).
			SetConnectTimeout(connectTimeout).
			SetMaxConnIdleTime(maxConnIdleTime).
			SetMinPoolSize(minPoolSize).
			SetMaxPoolSize(maxPoolSize))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
