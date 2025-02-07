package engine

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_project_template/internal/conf"
	"golang.org/x/net/context"
	"time"
)

func newMongo(mongoBase *conf.Mongodb) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//construct url: mongodb://username:password@127.0.0.1:27017/dbname
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoBase.ConnectStr).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		panic(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	//collection := client.Database("eobd")
	return client
}
//
//func MongoDXJDisconnect() {
//	DB.MongoDXJ.Disconnect(context.TODO())
//}
