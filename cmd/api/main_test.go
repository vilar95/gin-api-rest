package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vilar95/gin-api-rest/database"
	"github.com/vilar95/gin-api-rest/internal/handler"
	"github.com/vilar95/gin-api-rest/internal/model"
	"github.com/vilar95/gin-api-rest/internal/repository"
	"gorm.io/gorm"
)

// Variáveis de teste — escopo do pacote para compartilhar entre funções auxiliares.
var (
	testID      int
	testDB      *gorm.DB
	testRepo    *repository.StudentRepository
	testHandler *handler.StudentHandler
)

// setupTest inicializa as dependências de teste (DB, repository, handler)
// e retorna um gin.Engine limpo para registrar rotas individuais.
//
// Decisão: inicializar DB apenas uma vez (lazy init) para evitar múltiplas
// conexões durante a execução dos testes.
func setupTest() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	if testDB == nil {
		var err error
		testDB, err = database.ConnectDatabase()
		if err != nil {
			panic("falha ao conectar ao banco de teste: " + err.Error())
		}
		testDB.AutoMigrate(&model.Student{})
		testRepo = repository.NewStudentRepository(testDB)
		testHandler = handler.NewStudentHandler(testRepo)
	}

	return gin.Default()
}

func createStudentMock() {
	student := model.Student{Name: "Jarbinhas a IA do Bem", RG: "123456789", CPF: "12345678901"}
	testDB.Create(&student)
	testID = int(student.ID)
}

func deleteStudentMock() {
	var student model.Student
	testDB.Where("id = ?", testID).First(&student)
	testDB.Delete(&student)
}

func TestGreeting(t *testing.T) {
	r := setupTest()
	r.GET("/api/greeting/:name", testHandler.Greeting)

	req, _ := http.NewRequest("GET", "/api/greeting/vilar", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	expected := `{"message":"E ai vilar, tudo bem?"}`
	assert.Equal(t, expected, res.Body.String())
}

func TestListStudents(t *testing.T) {
	r := setupTest()
	createStudentMock()
	defer deleteStudentMock()

	r.GET("/api/students", testHandler.List)

	req, _ := http.NewRequest("GET", "/api/students", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetStudentByCPF(t *testing.T) {
	r := setupTest()
	createStudentMock()
	defer deleteStudentMock()

	r.GET("/api/students/cpf/:cpf", testHandler.GetByCPF)

	// Decisão: usar o CPF exato que está no banco ("12345678901").
	// O teste anterior usava "123.456.789-10" (formatado), que não correspondia
	// ao valor armazenado — era efetivamente um teste de "não encontrado" acidental.
	req, _ := http.NewRequest("GET", "/api/students/cpf/12345678901", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetStudentByID(t *testing.T) {
	r := setupTest()
	createStudentMock()
	defer deleteStudentMock()

	r.GET("/api/students/:id", testHandler.GetByID)

	path := "/api/students/" + strconv.Itoa(testID)
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var student model.Student
	json.Unmarshal(res.Body.Bytes(), &student)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Jarbinhas a IA do Bem", student.Name)
	assert.Equal(t, "123456789", student.RG)
	assert.Equal(t, "12345678901", student.CPF)
}

func TestDeleteStudent(t *testing.T) {
	r := setupTest()
	createStudentMock()

	r.DELETE("/api/students/:id", testHandler.Delete)

	path := "/api/students/" + strconv.Itoa(testID)
	req, _ := http.NewRequest("DELETE", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestUpdateStudent(t *testing.T) {
	r := setupTest()
	createStudentMock()
	defer deleteStudentMock()

	r.PATCH("/api/students/:id", testHandler.Update)

	path := "/api/students/" + strconv.Itoa(testID)
	update := model.Student{Name: "Airton Marques", RG: "423456781", CPF: "52345678901"}
	body, _ := json.Marshal(update)

	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var updated model.Student
	json.Unmarshal(res.Body.Bytes(), &updated)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Airton Marques", updated.Name)
	assert.Equal(t, "423456781", updated.RG)
	assert.Equal(t, "52345678901", updated.CPF)
}
