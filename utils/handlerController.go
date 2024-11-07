package utils

import "github.com/gin-gonic/gin"


func WrapHandler(controller func(*gin.Context) ApiResponse) gin.HandlerFunc{
	return func(c *gin.Context){
		response := controller(c)

		c.JSON(response.StatusCode, response.Body)
	}
}