package main

import "github.com/vilar95/gin-api-rest/internal/router"

func main() {
	router := router.SetupRouter()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	router.Run()
}
