package config

import (
	"boiler-go/database/mongo"
	"context"
	"log"
	"time"
)

func InitMongoDatabase() mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.GetMongoClient(App.Config.GetString(`mongo.url`))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

    log.Println("Connected to MongoDB!")

	return client
}
