package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vilar95/gin-api-rest/database"
	"github.com/vilar95/gin-api-rest/internal/controller"
	"github.com/vilar95/gin-api-rest/model"
)

var ID int

func SetupRouteTest() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routers := gin.Default()
	return routers
}

func CreateStudentMock() {
	student := model.Student{Name: "Nome do aluno teste", RG: "123456789", CPF: "12345678901"}
	database.DB.Create(&student)
	ID = int(student.ID)
}

func DeleteStudentMock() {
	var student model.Student
	database.DB.Where("id = ?", ID).First(&student)
	database.DB.Delete(&student)
}

func TestVerifyStatusCode(t *testing.T) {
	r := SetupRouteTest()
	r.GET("/:name", controller.Greeting)
	req, _ := http.NewRequest("GET", "/vilar", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	// Verifica se o status code é 200 e se a resposta é a esperada
	assert.Equal(t, http.StatusOK, res.Code, "O status code esperado era 200, mas veio %d", res.Code)
	responseMock := `{"API diz: ":"E ai vilar, tudo bem?"}`
	assert.Equal(t, responseMock, res.Body.String(), "A resposta esperada era %s, mas veio %s", responseMock, res.Body.String())
}

func TestVerifyGetAllStudents(t *testing.T) {
	database.ConnectDatabase()
	CreateStudentMock()
	defer DeleteStudentMock()
	r := SetupRouteTest()
	r.GET("/all-students", controller.ShowAllStudents)
	req, _ := http.NewRequest("GET", "/all-students", nil)
	// Armazena a resposta da requisição na variável
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "O status code esperado era 200, mas veio %d", res.Code)
}

func TestFetchStudentByCPF(t *testing.T) {
	database.ConnectDatabase()
	CreateStudentMock()
	defer DeleteStudentMock()
	r := SetupRouteTest()
	r.GET("/student/cpf/:cpf", controller.FetchStudentByCPF)
	req, _ := http.NewRequest("GET", "/student/cpf/123.456.789-10", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "O status code esperado era 200, mas veio %d", res.Code)
}
