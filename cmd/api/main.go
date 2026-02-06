package main

import (
	"github.com/vilar95/gin-api-rest/database"
	"github.com/vilar95/gin-api-rest/internal/router"
	"github.com/vilar95/gin-api-rest/model"
)

func main() {
	database.ConnectDatabase()
	model.Students = []model.Student{
		{Name: "Eduardo Vilar", CPF: "123.456.789-00", RG: "12.345.678-9"},
		{Name: "Paula Azevedo", CPF: "987.654.321-00", RG: "98.765.432-1"},

	}
	router := router.SetupRouter()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	router.Run()
}
