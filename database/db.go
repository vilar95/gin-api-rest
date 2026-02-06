package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vilar95/gin-api-rest/model"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectDatabase() {
	dsn := "host=localhost user=root password=root dbname=root port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic("A conexão com o banco de dados falhou!")
	}
	DB.AutoMigrate(&model.Student{})
}
