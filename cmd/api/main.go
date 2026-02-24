package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vilar95/gin-api-rest/database"
	"github.com/vilar95/gin-api-rest/internal/model"
	"github.com/vilar95/gin-api-rest/internal/router"
)

func main() {
	// Conectar ao banco de dados.
	// Decisão: tratar erro com log.Fatal em vez de panic.
	// log.Fatal é mais idiomático em main() — imprime a mensagem e chama os.Exit(1).
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}

	// Auto-migrate: cria/atualiza tabelas automaticamente.
	if err := db.AutoMigrate(&model.Student{}); err != nil {
		log.Fatal("Erro ao executar auto-migrate: ", err)
	}

	// Configurar rotas — passamos o db como dependência explícita.
	r := router.SetupRouter(db)

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Erro ao configurar trusted proxies: ", err)
	}

	// Configurar porta do servidor.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Decisão: usar http.Server + graceful shutdown em vez de r.Run().
	// Graceful shutdown aguarda requisições em andamento terminarem antes
	// de encerrar o processo, evitando respostas cortadas.
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Iniciar servidor em goroutine separada.
	go func() {
		log.Printf("Servidor rodando em http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Erro ao iniciar servidor: ", err)
		}
	}()

	// Aguardar sinal de interrupção (Ctrl+C ou kill).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando servidor...")

	// Dar até 5 segundos para requisições em andamento terminarem.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Erro ao desligar servidor: ", err)
	}

	log.Println("Servidor encerrado.")
}
