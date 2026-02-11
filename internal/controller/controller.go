package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/database"
	"github.com/vilar95/gin-api-rest/model"
)

func ShowAllStudents(c *gin.Context) {
	var students []model.Student
	database.DB.Find(&students)
	c.JSON(200, students)
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
	if err := model.ValidateStudentInfo(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}

	database.DB.Create(&student)
	c.JSON(http.StatusCreated, student)
}

func FetchStudentByID(c *gin.Context) {
	var student model.Student
	id := c.Param("id")

	if err := database.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Estudante não encontrado, verifique o ID e tente novamente."})
		return
	}

	c.JSON(http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	var student model.Student
	id := c.Param("id")

	if err := database.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Estudante não encontrado, verifique o ID e tente novamente."})
		return
	}

	database.DB.Delete(&student)
	c.JSON(http.StatusOK, gin.H{"Deleted": "Estudante deletado com sucesso."})
}

func UpdateStudent(c *gin.Context) {
	var student model.Student
	id := c.Param("id")

	if err := database.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Estudante não encontrado, verifique o ID e tente novamente."})
		return
	}

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := model.ValidateStudentInfo(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}

	database.DB.Save(&student)
	c.JSON(http.StatusOK, student)
}

func FetchStudentByCPF(c *gin.Context) {
	var student model.Student
	cpf := c.Param("cpf")

	if err := database.DB.Where("cpf = ?", cpf).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Estudante não encontrado, verifique o CPF e tente novamente."})
		return
	}

	c.JSON(http.StatusOK, student)
}
