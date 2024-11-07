package controller

import (
	"boiler-go/entities"
	"boiler-go/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserUseCase entities.UserUsecase
}

func NewUserHandler(G * gin.RouterGroup, uu entities.UserUsecase) {
	handler := &UserHandler{
		UserUseCase: uu,
	}
	group := G.Group("/users")
	group.POST("/", utils.WrapHandler(handler.InsertOne))
	group.GET("/:id", utils.WrapHandler(handler.FindOne))
	group.PUT("/:id", utils.WrapHandler(handler.UpdateOne))
}

func isRequestValid(m *entities.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (userH *UserHandler) InsertOne(c *gin.Context) utils.ApiResponse{
	var (
		usr entities.User
	)
	if err := c.ShouldBindJSON(&usr); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err)
	}
	if ok, err := isRequestValid(&usr); !ok {
		return utils.ErrorResponse(http.StatusBadRequest, err)

	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := userH.UserUseCase.InsertOne(ctx, &usr)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnprocessableEntity, err)
	}
	return utils.SuccessResponse(http.StatusOK, result)
}

func (userH *UserHandler) FindOne(c *gin.Context) utils.ApiResponse {
	var (
		id = c.Param("id")
	)
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := userH.UserUseCase.FindOne(ctx, id)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, )
	}

	return utils.SuccessResponse(http.StatusOK, result, )
}

func (userH *UserHandler) UpdateOne(c *gin.Context) utils.ApiResponse {
	id := c.Param("id")

	var (
		usr entities.User
	)

	if err := c.ShouldBindJSON(&usr); err != nil {
		return utils.ErrorResponse(http.StatusUnprocessableEntity, err, )
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := userH.UserUseCase.UpdateOne(ctx, &usr, id)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, )
	}

	return utils.SuccessResponse(http.StatusOK, result)
}
