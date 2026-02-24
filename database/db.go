package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase abre a conexão com o PostgreSQL e retorna a instância do GORM.
//
// Decisão: retornar (*gorm.DB, error) em vez de usar variável global + panic.
// - Variáveis globais dificultam testes e escondem dependências.
// - Retornar error segue o padrão idiomático de Go: quem chama decide como tratar.
// - O chamador (main.go) pode usar log.Fatal ou qualquer outra estratégia.
//
// As credenciais vêm de variáveis de ambiente com fallback para valores padrão,
// permitindo configuração flexível sem alterar código.
func ConnectDatabase() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "root")
	dbname := getEnv("DB_NAME", "root")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	return db, nil
}

// getEnv retorna o valor da variável de ambiente ou o fallback informado.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
