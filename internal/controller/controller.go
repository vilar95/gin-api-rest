package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/model"
)

func ShowAllStudents(c *gin.Context) {
	c.JSON(200, model.Students)
}

func Greeting(c *gin.Context) {
	name := c.Param("name")
	c.JSON(200, gin.H{
		"API diz: ": "E ai " + name + ", tudo bem?",
	})
}
