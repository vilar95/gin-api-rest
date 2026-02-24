package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vilar95/gin-api-rest/internal/handler"
	"github.com/vilar95/gin-api-rest/internal/repository"
	"gorm.io/gorm"
)

// SetupRouter configura e retorna o engine do Gin com todas as rotas.
//
// Decisão: receber *gorm.DB como parâmetro em vez de usar variável global.
// Isso torna as dependências explícitas — ao ler a assinatura, fica claro
// que o router precisa de uma conexão com o banco.
//
// Aqui montamos a cadeia de injeção de dependência:
// DB → Repository → Handler → Rotas
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "templates/assets")

	// Injeção de dependência: criamos repository e handler com dependências explícitas
	studentRepo := repository.NewStudentRepository(db)
	studentHandler := handler.NewStudentHandler(studentRepo)

	// Decisão: agrupar rotas da API sob /api com nomes RESTful.
	// Antes: POST /create-student, DELETE /delete-student/:id (verbos na URL).
	// Agora: POST /api/students, DELETE /api/students/:id (substantivos + método HTTP).
	// Isso segue as convenções REST e evita conflito de rotas (o antigo GET /:name
	// capturava qualquer rota de primeiro nível).
	api := r.Group("/api")
	{
		api.GET("/greeting/:name", studentHandler.Greeting)
		api.GET("/students", studentHandler.List)
		api.POST("/students", studentHandler.Create)
		api.GET("/students/:id", studentHandler.GetByID)
		api.PATCH("/students/:id", studentHandler.Update)
		api.DELETE("/students/:id", studentHandler.Delete)
		api.GET("/students/cpf/:cpf", studentHandler.GetByCPF)
	}

	// Rotas de páginas HTML
	r.GET("/", studentHandler.ShowIndex)
	r.NoRoute(handler.NotFound)

	return r
}
