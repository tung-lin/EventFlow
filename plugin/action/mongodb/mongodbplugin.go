package mongodb

import (
	"EventFlow/common/interface/pluginbase"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbPlugin struct {
	messageCannel chan string
	pluginbase.ActionHandler
	Setting SettingConfig
}

func NewMongodbPlugin() *MongodbPlugin {

	messageCannel = make(chan string, 50)

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		for {
			message, more := <-messageCannel
			if more {
			} else {

				break
			}
		}
	}()
}

func InitialConnection() {

}

func (action *MongodbPlugin) FireAction(messageFromTrigger *string, parameters *map[string]interface{}) {

}
