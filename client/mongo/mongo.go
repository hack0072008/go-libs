package mongoClient

import (
	"context"
	"strings"

	"github.com/hack0072008/go-libs/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func Connect(mongodbStr string) error {
	var clientOptions options.ClientOptions
	adds := strings.Split(mongodbStr, ",")
	for _, addr := range adds {
		opt := options.Client().ApplyURI(addr)
		if clientOptions.Auth == nil && opt.Auth != nil {
			clientOptions.Auth = opt.Auth
		}
		if opt.Hosts != nil {
			clientOptions.Hosts = append(clientOptions.Hosts, opt.Hosts...)
		}
	}
	c, err := mongo.Connect(context.Background(), &clientOptions)
	if err != nil {
		log.Errorf("mongo[%s] connect error,%s", mongodbStr, err.Error())
		return err
	}
	if err = c.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Errorf("mongo[%s] ping error,%s", mongodbStr, err.Error())
		return err
	}
	client = c
	log.Info("successfully connect to " + mongodbStr)
	return nil
}

func Disconnect() {
	if client != nil {
		_ = client.Disconnect(context.TODO())
	}
}

func GetClient() *mongo.Client {
	return client
}
