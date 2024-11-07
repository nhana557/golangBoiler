package app

import (
	"boiler-go/config"
	"boiler-go/libs/auth/controller"
	"boiler-go/libs/auth/usecase"
	_jwt "boiler-go/libs/jwt/usecase"
	userHttp "boiler-go/libs/users/controller/http"
	userRepo "boiler-go/libs/users/repository"
	userUseCase "boiler-go/libs/users/usecase"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunApp(router *gin.Engine){
	var (
		ctx context.Context
		v1 *gin.RouterGroup = router.Group("/api/v1")
	)
	defer config.App.Mongo.Disconnect(ctx)
	if config.App.Config.GetString("GIN_NODE") == "release"{
		gin.SetMode(gin.ReleaseMode)
	}
	
	// gin.SetMode(gin.ReleaseMode)
	// router.Use(gin.Recovery())
	router.Use(cors.Default())

	timeoutContext := time.Duration(config.App.Config.GetInt("context.timeout")) * time.Second

	database := config.App.Mongo.Database(config.App.Config.GetString("mongo.dbName"))
	// cacheRedis := config.App.Redis


	
	
	userRepo := userRepo.NewMongoRepository(database)
	usrUsecase := userUseCase.NewUserUsecase(userRepo, timeoutContext)
	userHttp.NewUserHandler(v1, usrUsecase)

	jwt := _jwt.NewJwtUsecase(userRepo, timeoutContext, config.App.Config)
	userJwt := v1.Group("")
	jwt.SetJwtUser(userJwt)
	adminJwt := v1.Group("")
	jwt.SetJwtUser(adminJwt)
	generalJwt := v1.Group("")
	jwt.SetJwtUser(generalJwt)

	// login handler
	loginUsecase  := usecase.NewLoginUseCase(userRepo, timeoutContext)
	controller.NewLoginHandler(v1, loginUsecase, config.App.Config)

	appPort := fmt.Sprintf(":%s", config.App.Config.GetString("server.address"))
	log.Fatal(router.Run(appPort))
}