package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	student := model.Student{Name: "Jarbinhas a IA do Bem", RG: "123456789", CPF: "12345678901"}
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

func TestFetchStudentBy(t *testing.T) {
	database.ConnectDatabase()
	CreateStudentMock()
	defer DeleteStudentMock()
	r := SetupRouteTest()
	r.GET("/student/:id", controller.FetchStudentByID)
	path := "/student/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	var stutentMock model.Student
	json.Unmarshal(res.Body.Bytes(), &stutentMock)
	fmt.Println(stutentMock.Name)
	assert.Equal(t, "Jarbinhas a IA do Bem", stutentMock.Name, "O nome do estudante esperado era 'Jarbinhas a IA do Bem', mas veio %s", stutentMock.Name)
	assert.Equal(t, "123456789", stutentMock.RG, "O RG do estudante esperado era '123456789', mas veio %s", stutentMock.RG)
	assert.Equal(t, "12345678901", stutentMock.CPF, "O CPF do estudante esperado era '12345678901', mas veio %s", stutentMock.CPF)
}

func TestDeleteStudent(t *testing.T) {
	database.ConnectDatabase()
	CreateStudentMock()
	r := SetupRouteTest()
	r.DELETE("/delete-student/:id", controller.DeleteStudent)
	path := "/delete-student/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "O status code esperado era 200, mas veio %d", res.Code)
}

func TestUpdateStudent(t *testing.T) {
	database.ConnectDatabase()
	CreateStudentMock()
	defer DeleteStudentMock()
	r := SetupRouteTest()
	r.PATCH("/update-student/:id", controller.UpdateStudent)
	path := "/update-student/" + strconv.Itoa(ID)
	student := model.Student{Name: "Airton Marques", RG: "423456781", CPF: "52345678901"}
	studentJSON, _ := json.Marshal(student)
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(studentJSON))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	var updatedStudent model.Student
	json.Unmarshal(res.Body.Bytes(), &updatedStudent)
	assert.Equal(t, "Airton Marques", updatedStudent.Name, "O nome do estudante esperado era 'Airton Marques', mas veio %s", updatedStudent.Name)
	assert.Equal(t, "423456781", updatedStudent.RG, "O RG do estudante esperado era '423456781', mas veio %s", updatedStudent.RG)
	assert.Equal(t, "52345678901", updatedStudent.CPF, "O CPF do estudante esperado era '52345678901', mas veio %s", updatedStudent.CPF)
	assert.Equal(t, http.StatusOK, res.Code, "O status code esperado era 200, mas veio %d", res.Code)
}
