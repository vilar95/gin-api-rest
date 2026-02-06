package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/database"
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

func CreateStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&student)
	c.JSON(http.StatusCreated, student)
}
