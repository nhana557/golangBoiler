package config

import (
	"boiler-go/database/mongo"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	App *Application
)

type Application struct {
	Config *viper.Viper
	MySql  *sql.DB
	Mongo  mongo.Client
	Redis	*redis.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = InitConfig()
	App.Mongo = InitMongoDatabase()
	App.Redis = InitRedis()
}
