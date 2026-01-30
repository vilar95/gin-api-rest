package controller

import "github.com/gin-gonic/gin"

func ShowAllStudents(c *gin.Context) {
	c.JSON(200, gin.H{
		"id":   1,
		"name": "John Doe",
	})
}

func Greeting(c *gin.Context) {
	name := c.Param("name")
	c.JSON(200, gin.H{
		"message": "Olá " + name,
	})
}