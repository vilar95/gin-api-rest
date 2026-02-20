package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/internal/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "templates/assets")

	r.GET("/:name", controller.Greeting)
	r.GET("/all-students", controller.ShowAllStudents)
	r.POST("/create-student", controller.CreateStudent)
	r.GET("/student/:id", controller.FetchStudentByID)
	r.DELETE("/delete-student/:id", controller.DeleteStudent)
	r.PATCH("/update-student/:id", controller.UpdateStudent)
	r.GET("/student/cpf/:cpf", controller.FetchStudentByCPF)
	r.GET("/index", controller.ShowIndexPage)

	r.NoRoute(controller.NotFound)

	return r
}
