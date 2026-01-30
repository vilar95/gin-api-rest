package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/internal/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/students", controller.ShowAllStudents)
	r.GET("/:name", controller.Greeting)
	return r
}
