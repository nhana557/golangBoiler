package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper{
	config := viper.New()
	
	config.SetConfigFile(`./config.yaml`)	
	err := config.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic("Cant Find File config.yaml")
	}

	if config.GetString(`GIN_NODE`) == "debug" {
		log.Println("Service RUN on DEBUG mode")
	}

	return config
}