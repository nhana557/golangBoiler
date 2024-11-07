package main

import (
	"boiler-go/app"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	app.RunApp(router)
}