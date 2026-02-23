# Gin API REST — Projeto de Estudo em Go

API REST para cadastro de estudantes, desenvolvida em Go com o framework [Gin](https://github.com/gin-gonic/gin) e banco de dados PostgreSQL via [GORM](https://gorm.io/). O projeto serve como base de estudo para desenvolvimento de APIs idiomáticas em Go.

## O que este projeto ensina

- Estruturação de projetos Go seguindo convenções da comunidade (`cmd/`, `internal/`)
- Criação de APIs REST com o framework Gin
- Persistência com GORM e PostgreSQL
- Injeção de dependência simples (sem frameworks)
- Separação de responsabilidades: handler → repository → model
- Validação de dados com tags declarativas
- Testes de integração com `httptest`
- Graceful shutdown do servidor HTTP
- Configuração via variáveis de ambiente

---

## Visão Geral da Arquitetura

```
gin-api-rest/
├── cmd/
│   └── api/
│       ├── main.go              # Ponto de entrada da aplicação
│       └── main_test.go         # Testes de integração
├── database/
│   └── db.go                    # Conexão com PostgreSQL
├── internal/
│   ├── handler/
│   │   └── student.go           # Handlers HTTP (recebe request, retorna response)
│   ├── model/
│   │   └── student.go           # Entidade do domínio + validação
│   ├── repository/
│   │   └── student.go           # Acesso ao banco de dados (queries)
│   └── router/
│       └── router.go            # Configuração de rotas
├── templates/                   # Páginas HTML
│   ├── index.html
│   ├── not_found.html
│   └── assets/
├── docs/
│   └── adding-new-feature.md    # Guia para implementar novas features
├── .env.dev                 # Exemplo de variáveis de ambiente
├── .gitignore
├── docker-compose.yml           # PostgreSQL + pgAdmin
├── go.mod
└── README.md
```

### Fluxo de uma requisição

```
HTTP Request
    │
    ▼
 Router (router.go)         → Define qual handler atende a rota
    │
    ▼
 Handler (handler/student.go) → Parse do request, validação, resposta HTTP
    │
    ▼
 Repository (repository/student.go) → Executa queries no banco via GORM
    │
    ▼
 Model (model/student.go)   → Struct da entidade + regras de validação
    │
    ▼
 Database (database/db.go)  → Conexão com PostgreSQL
```

### Responsabilidade de cada camada

| Camada | Pacote | O que faz | O que NÃO faz |
|--------|--------|-----------|----------------|
| **Handler** | `internal/handler` | Parse HTTP, validação, serialização JSON | Queries SQL, regras de negócio complexas |
| **Repository** | `internal/repository` | CRUD no banco de dados | Lógica HTTP, validação de input |
| **Model** | `internal/model` | Define structs do domínio, regras de validação | Acesso ao banco, lógica HTTP |
| **Router** | `internal/router` | Mapeia URLs para handlers | Lógica de negócio |
| **Database** | `database` | Abre conexão com PostgreSQL | Queries específicas |

### Por que `internal/`?

Em Go, o diretório `internal/` tem significado especial: código dentro dele **não pode ser importado por outros módulos**. Isso protege a implementação interna do projeto. Apenas `cmd/` e outros pacotes dentro do mesmo módulo podem usá-lo.

---

## Endpoints da API

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/greeting/:name` | Saudação personalizada |
| GET | `/api/students` | Lista todos os estudantes |
| POST | `/api/students` | Cria um novo estudante |
| GET | `/api/students/:id` | Busca estudante por ID |
| PATCH | `/api/students/:id` | Atualiza estudante por ID |
| DELETE | `/api/students/:id` | Remove estudante por ID |
| GET | `/api/students/cpf/:cpf` | Busca estudante por CPF |
| GET | `/` | Página HTML com lista de estudantes |

### Modelo de dados

```json
{
  "name": "João Silva",
  "rg": "123456789",
  "cpf": "12345678901"
}
```

**Regras de validação:**
- `name`: obrigatório (não pode ser vazio)
- `rg`: exatamente 9 dígitos numéricos
- `cpf`: exatamente 11 dígitos numéricos

---

## Como Rodar o Projeto

### Pré-requisitos

- [Go 1.24+](https://go.dev/dl/)
- [Docker e Docker Compose](https://docs.docker.com/get-docker/) (para o banco de dados)

### 1. Clonar o repositório

```bash
git clone https://github.com/vilar95/gin-api-rest.git
cd gin-api-rest
```

### 2. Subir o banco de dados

```bash
docker compose up -d
```

Isso inicia:
- **PostgreSQL** na porta `5432` (user: `root`, password: `root`, database: `root`)
- **pgAdmin** na porta `54321` (email: configurado no docker-compose.yml)

### 3. Configurar variáveis de ambiente (opcional)

```bash
cp .env.dev .env
# Edite .env se necessário (os valores padrão funcionam com o docker-compose)
```

Se você não criar o `.env`, o projeto usa os valores padrão definidos em `database/db.go`.

### 4. Instalar dependências

```bash
go mod download
```

### 5. Rodar o servidor

```bash
go run cmd/api/main.go
```

O servidor inicia em `http://localhost:8080`.

### 6. Testar a API

```bash
# Criar estudante
curl -X POST http://localhost:8080/api/students \
  -H 'Content-Type: application/json' \
  -d '{"name": "João Silva", "rg": "123456789", "cpf": "12345678901"}'

# Listar todos
curl http://localhost:8080/api/students

# Buscar por ID
curl http://localhost:8080/api/students/1

# Buscar por CPF
curl http://localhost:8080/api/students/cpf/12345678901

# Atualizar
curl -X PATCH http://localhost:8080/api/students/1 \
  -H 'Content-Type: application/json' \
  -d '{"name": "João Atualizado", "rg": "987654321", "cpf": "10987654321"}'

# Deletar
curl -X DELETE http://localhost:8080/api/students/1

# Saudação
curl http://localhost:8080/api/greeting/mundo
```

### 7. Rodar os testes

```bash
# Requer o banco de dados rodando (docker compose up -d)
go test ./cmd/api/ -v
```

---

## Como Contribuir / Estudar Este Projeto

### Roteiro sugerido para leitura do código

1. **Comece pelo model** ([internal/model/student.go](internal/model/student.go))
   Entenda a entidade principal e como a validação funciona com tags.

2. **Veja o repository** ([internal/repository/student.go](internal/repository/student.go))
   Observe como o GORM traduz chamadas de método em queries SQL.

3. **Estude os handlers** ([internal/handler/student.go](internal/handler/student.go))
   Veja como cada endpoint recebe a request, chama o repository e retorna a response.

4. **Analise o router** ([internal/router/router.go](internal/router/router.go))
   Entenda como as rotas são organizadas e como a injeção de dependência acontece.

5. **Veja o main** ([cmd/api/main.go](cmd/api/main.go))
   Observe como tudo é conectado: banco → router → servidor HTTP.

6. **Leia os testes** ([cmd/api/main_test.go](cmd/api/main_test.go))
   Veja como criar testes de integração usando `httptest`.

### Guia para implementar novas features

Veja o guia passo a passo em [docs/adding-new-feature.md](docs/adding-new-feature.md).

---

## Decisões de Arquitetura

Registro das principais decisões tomadas e suas motivações:

### 1. Retornar `error` em vez de usar `panic`

**Antes:** `database.ConnectDatabase()` usava `panic()` se a conexão falhasse.
**Agora:** Retorna `(*gorm.DB, error)` — quem chama decide como tratar.

**Por quê:** Em Go, `panic` é reservado para situações irrecuperáveis (bugs no código). Erros esperados (banco offline, credenciais erradas) devem ser tratados com o padrão `if err != nil`. Isso é fundamental na filosofia Go.

### 2. Eliminar variáveis globais (`database.DB`)

**Antes:** `database.DB` era uma variável global acessada diretamente pelos handlers.
**Agora:** O `*gorm.DB` é passado como dependência via construtor (injeção de dependência).

**Por quê:** Variáveis globais criam acoplamento oculto — ao ler um handler, não fica claro que ele depende do banco. Com injeção de dependência, as dependências são explícitas na assinatura das funções e construtores.

### 3. Separar handler e repository

**Antes:** O `controller` fazia tudo — parse HTTP, validação e queries no banco.
**Agora:** `handler` cuida de HTTP; `repository` cuida de persistência.

**Por quê:** Misturar responsabilidades dificulta testes e manutenção. Se amanhã trocarmos o GORM por queries SQL puras, só o repository muda — os handlers continuam iguais.

### 4. Renomear `controller` para `handler`

**Antes:** `internal/controller/`
**Agora:** `internal/handler/`

**Por quê:** Em Go, o termo "handler" é mais idiomático (vem de `http.Handler`). "Controller" é mais comum em frameworks MVC como Rails ou Spring.

### 5. Mover model para `internal/model`

**Antes:** `model/` na raiz do projeto.
**Agora:** `internal/model/`

**Por quê:** O model é interno a este projeto — não precisa ser exportado. O diretório `internal/` em Go garante isso a nível de compilação.

### 6. Rotas RESTful com grupo `/api`

**Antes:** `POST /create-student`, `DELETE /delete-student/:id`, `GET /:name`
**Agora:** `POST /api/students`, `DELETE /api/students/:id`, `GET /api/greeting/:name`

**Por quê:** Em REST, URLs representam recursos (substantivos), não ações (verbos). O método HTTP já indica a ação. Além disso, `GET /:name` capturava qualquer rota de primeiro nível, causando conflitos.

### 7. Validação no model via método `Validate()`

**Antes:** Função solta `ValidateStudentInfo(&student)` + validação manual duplicada no `UpdateStudent`.
**Agora:** Método `student.Validate()` reutilizado em todos os handlers.

**Por quê:** Duplicação de regras de validação é fonte de bugs. Com um único ponto de validação, mudanças nas regras se propagam automaticamente.

### 8. Graceful shutdown

**Antes:** `router.Run()` sem tratamento de encerramento.
**Agora:** `http.Server` com `Shutdown()` que aguarda requisições em andamento.

**Por quê:** Se o servidor recebe um `Ctrl+C` durante uma requisição, o `Run()` simples corta a conexão imediatamente. O graceful shutdown dá tempo para as requisições em andamento terminarem.

### 9. Configuração via variáveis de ambiente

**Antes:** DSN do banco hardcoded no código.
**Agora:** Lido de variáveis de ambiente com fallback para valores padrão.

**Por quê:** Hardcoding credenciais é uma má prática — impede deploy em diferentes ambientes (dev, staging, prod). Variáveis de ambiente seguem o princípio [12-Factor App](https://12factor.net/config).

---

## Tecnologias

- **Go 1.24** — Linguagem de programação
- **Gin** — Framework HTTP leve e rápido
- **GORM** — ORM para Go
- **PostgreSQL** — Banco de dados relacional
- **validator.v2** — Validação declarativa via tags
- **testify** — Assertions para testes
- **Docker Compose** — Orquestração de containers



