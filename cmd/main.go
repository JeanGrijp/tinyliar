package main

import (
	"log"
	"net/http"

	"github.com/JeanGrijp/tinyliar/internal/config"
	"github.com/JeanGrijp/tinyliar/internal/handler"
	"github.com/JeanGrijp/tinyliar/internal/repository"
	"github.com/JeanGrijp/tinyliar/internal/routes"
)

func main() {
	// Conecta no banco

	config.ConnectToDatabase()
	db := config.DB
	if db == nil {
		log.Fatal("Banco de dados não inicializado")
	}

	// Cria o repositório
	linkRepo := repository.NewLinkRepository(config.DB)

	// Cria o handler com injeção de dependência
	linkHandler := &handler.LinkHandler{
		Repo: linkRepo,
	}

	// Configura as rotas com os handlers injetados
	r := routes.SetupRoutes(linkHandler)

	// Sobe o servidor
	log.Println("Servidor rodando em http://localhost:8888")
	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
