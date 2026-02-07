## Gin API REST

Esta é uma API REST desenvolvida em Go utilizando o framework Gin, com integração ao banco de dados PostgreSQL via GORM.

### Funcionalidades

- CRUD de estudantes (Create, Read, Update, Delete)
- Busca de estudante por ID e CPF
- Mensagem de saudação personalizada

### Modelo

O modelo principal é o estudante:

```
type Student struct {
  Name string
  CPF  string
  RG   string
}
```

### Endpoints

| Método | Rota                      | Descrição                                 |
|--------|---------------------------|--------------------------------------------|
| GET    | /:name                    | Saudação personalizada                     |
| GET    | /all-students             | Lista todos os estudantes                  |
| POST   | /create-student           | Cria um novo estudante                     |
| GET    | /student/:id              | Busca estudante por ID                     |
| DELETE | /delete-student/:id       | Deleta estudante por ID                    |
| PATCH  | /update-student/:id       | Atualiza estudante por ID                  |
| GET    | /student/cpf/:cpf         | Busca estudante por CPF                    |

### Banco de Dados

O banco utilizado é PostgreSQL. A conexão é feita via GORM e os dados de acesso estão definidos em `database/db.go`.

### Como executar

1. Instale as dependências:
	- Go
	- Docker (opcional, para rodar o banco via docker-compose)
2. Configure o banco de dados PostgreSQL (ou use o docker-compose.yml)
3. Execute o projeto:
	```bash
	go run cmd/api/main.go
	```

### Exemplo de requisição

```bash
curl -X POST http://localhost:8080/create-student \
	  -H 'Content-Type: application/json' \
	  -d '{"name": "João", "cpf": "12345678900", "rg": "123456"}'
```


