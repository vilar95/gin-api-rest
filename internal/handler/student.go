package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/internal/model"
	"github.com/vilar95/gin-api-rest/internal/repository"
)

// StudentHandler agrupa os handlers HTTP relacionados a Student.
//
// Decisão: usar struct com dependência injetada (repository) em vez de funções
// soltas que acessam variáveis globais (database.DB).
// Vantagens:
// - Dependências explícitas e rastreáveis.
// - Facilita testes (pode injetar um repository de teste).
// - Segue o padrão idiomático de Go para handlers com estado.
type StudentHandler struct {
	repo *repository.StudentRepository
}

// NewStudentHandler cria uma nova instância do handler com suas dependências.
func NewStudentHandler(repo *repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

// List retorna todos os estudantes em JSON.
// GET /api/students
func (h *StudentHandler) List(c *gin.Context) {
	students, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar estudantes"})
		return
	}
	c.JSON(http.StatusOK, students)
}

// Create valida e insere um novo estudante.
// POST /api/students
func (h *StudentHandler) Create(c *gin.Context) {
	var student model.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decisão: chamar student.Validate() como método — mais idiomático que
	// a versão anterior ValidateStudentInfo(&student).
	if err := student.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar estudante"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// GetByID busca um estudante pelo ID.
// GET /api/students/:id
func (h *StudentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	student, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "estudante não encontrado"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Update atualiza os dados de um estudante existente.
// PATCH /api/students/:id
//
// Decisão: reutilizar input.Validate() em vez da validação manual duplicada
// que existia antes (if input.Name == "" || len(input.RG) != 9 ...).
// Isso garante consistência — as regras vivem em um só lugar (tags do model).
func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	student, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "estudante não encontrado"})
		return
	}

	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student.Name = input.Name
	student.RG = input.RG
	student.CPF = input.CPF

	if err := h.repo.Update(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar estudante"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Delete remove um estudante pelo ID.
// DELETE /api/students/:id
func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	student, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "estudante não encontrado"})
		return
	}

	if err := h.repo.Delete(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao deletar estudante"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estudante deletado com sucesso"})
}

// GetByCPF busca um estudante pelo CPF.
// GET /api/students/cpf/:cpf
func (h *StudentHandler) GetByCPF(c *gin.Context) {
	cpf := c.Param("cpf")

	student, err := h.repo.FindByCPF(cpf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "estudante não encontrado"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Greeting retorna uma saudação personalizada.
// GET /api/greeting/:name
func (h *StudentHandler) Greeting(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "E ai " + name + ", tudo bem?",
	})
}

// ShowIndex renderiza a página HTML com a lista de estudantes.
// GET /
func (h *StudentHandler) ShowIndex(c *gin.Context) {
	students, err := h.repo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"students": nil})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"students": students,
	})
}

// NotFound renderiza a página 404.
// Decisão: manter como função de pacote (não precisa do repository).
func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "not_found.html", nil)
}
