# Guia: Como Implementar uma Nova Feature

Este guia usa como exemplo a implementação de uma feature de **busca de estudantes por nome** (busca parcial). Siga os mesmos passos para qualquer nova feature.

---

## Passo 1: Definir a regra de negócio

Antes de escrever código, pergunte-se:

- **O que essa feature faz?** Busca estudantes cujo nome contenha um trecho informado.
- **Onde ela se encaixa?** É uma operação de leitura (`GET`) sobre o recurso `students`.
- **Qual a rota?** `GET /api/students/search?name=João`
- **Qual a resposta esperada?** Lista de estudantes que correspondem à busca.

---

## Passo 2: Ajustar o Model (se necessário)

Neste caso, não precisamos alterar o model — a busca usa campos já existentes. Mas se a feature exigir novos campos, edite `internal/model/student.go`:

```go
// Exemplo: adicionar campo Email
type Student struct {
    gorm.Model
    Name  string `json:"name" validate:"nonzero"`
    RG    string `json:"rg" validate:"len=9,regexp=^[0-9]*$"`
    CPF   string `json:"cpf" validate:"len=11,regexp=^[0-9]*$"`
    Email string `json:"email" validate:"nonzero"` // ← novo campo
}
```

> **Importante:** Após adicionar campos, o GORM faz auto-migrate automaticamente ao iniciar o servidor. Para produção, use ferramentas de migration como [golang-migrate](https://github.com/golang-migrate/migrate).

---

## Passo 3: Adicionar método no Repository

O repository é responsável por queries no banco. Adicione o método em `internal/repository/student.go`:

```go
// SearchByName busca estudantes cujo nome contenha o trecho informado.
// Usa ILIKE para busca case-insensitive (específico do PostgreSQL).
func (r *StudentRepository) SearchByName(name string) ([]model.Student, error) {
    var students []model.Student
    if err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&students).Error; err != nil {
        return nil, err
    }
    return students, nil
}
```

**Por que no repository?**
- O handler não deve saber detalhes de SQL/GORM.
- Se trocarmos o banco ou ORM, só o repository muda.
- Facilita testar o handler isoladamente (pode usar mock do repository).

---

## Passo 4: Adicionar o Handler

O handler cuida de HTTP: parse da request, chamada ao repository e montagem da response. Adicione em `internal/handler/student.go`:

```go
// SearchByName busca estudantes por nome (busca parcial).
// GET /api/students/search?name=João
func (h *StudentHandler) SearchByName(c *gin.Context) {
    name := c.Query("name")
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro 'name' é obrigatório"})
        return
    }

    students, err := h.repo.SearchByName(name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar estudantes"})
        return
    }

    c.JSON(http.StatusOK, students)
}
```

**Observe o padrão do handler:**
1. Ler dados da request (`c.Query`, `c.Param`, `c.ShouldBindJSON`)
2. Validar input
3. Chamar o repository
4. Retornar response com status HTTP adequado

---

## Passo 5: Registrar a Rota

Adicione a rota no grupo `/api` em `internal/router/router.go`:

```go
api := r.Group("/api")
{
    // ... rotas existentes ...
    api.GET("/students/search", studentHandler.SearchByName) // ← nova rota
}
```

**Atenção à ordem:** Rotas mais específicas devem vir antes de rotas com parâmetros. Por exemplo, `/students/search` deve ser registrada **antes** de `/students/:id`, senão o Gin interpreta "search" como um ID.

---

## Passo 6: Escrever Testes

Adicione um teste em `cmd/api/main_test.go`:

```go
func TestSearchStudentByName(t *testing.T) {
    r := setupTest()
    createStudentMock() // Cria "Jarbinhas a IA do Bem"
    defer deleteStudentMock()

    r.GET("/api/students/search", testHandler.SearchByName)

    // Busca parcial por "Jarbinhas"
    req, _ := http.NewRequest("GET", "/api/students/search?name=Jarbinhas", nil)
    res := httptest.NewRecorder()
    r.ServeHTTP(res, req)

    assert.Equal(t, http.StatusOK, res.Code)

    var students []model.Student
    json.Unmarshal(res.Body.Bytes(), &students)
    assert.Greater(t, len(students), 0, "deveria encontrar pelo menos 1 estudante")
    assert.Contains(t, students[0].Name, "Jarbinhas")
}

func TestSearchStudentByNameEmpty(t *testing.T) {
    r := setupTest()
    r.GET("/api/students/search", testHandler.SearchByName)

    // Sem parâmetro name → deve retornar 400
    req, _ := http.NewRequest("GET", "/api/students/search", nil)
    res := httptest.NewRecorder()
    r.ServeHTTP(res, req)

    assert.Equal(t, http.StatusBadRequest, res.Code)
}
```

**Rode os testes:**

```bash
go test ./cmd/api/ -v -run TestSearchStudent
```

---

## Resumo: Checklist para Nova Feature

- [ ] **Definir:** O que a feature faz? Qual rota? Qual response?
- [ ] **Model:** Precisa de novos campos? Adicionar em `internal/model/student.go`
- [ ] **Repository:** Criar método de acesso ao banco em `internal/repository/student.go`
- [ ] **Handler:** Criar handler HTTP em `internal/handler/student.go`
- [ ] **Router:** Registrar a rota em `internal/router/router.go`
- [ ] **Testes:** Escrever pelo menos um teste de sucesso e um de erro
- [ ] **Testar manualmente:** `curl` ou Postman para validar

---

## Dicas para Expandir o Projeto

Depois de dominar o fluxo básico, experimente:

1. **Adicionar uma nova entidade** (ex: `Course`) com relacionamento N:N com `Student`
2. **Implementar paginação** no `GET /api/students` usando query params `?page=1&limit=10`
3. **Adicionar autenticação** com JWT middleware
4. **Criar DTOs** (Data Transfer Objects) separados para input/output, desacoplando a API do model de banco
5. **Extrair interfaces** do repository para permitir mocks em testes unitários:
   ```go
   type StudentFinder interface {
       FindByID(id string) (model.Student, error)
       FindAll() ([]model.Student, error)
   }
   ```
6. **Adicionar middleware** de logging, CORS, rate limiting
7. **Usar migrations** com [golang-migrate](https://github.com/golang-migrate/migrate) em vez de AutoMigrate
