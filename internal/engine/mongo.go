package engine

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_project_template/internal/conf"
	"golang.org/x/net/context"
	"time"
)

func newMongo(mongoBase *conf.Mongodb) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	poolSize := mongoBase.MaxPoolSize
	if poolSize == 0 {
		poolSize = 20
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoBase.ConnectStr).SetMaxPoolSize(poolSize))
	if err != nil {
		return nil, fmt.Errorf("mongo connect failed: %w", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}
	return client, nil
}
