package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Importa driver do PostgreSQL
)

var DB *sqlx.DB

func ConnectToDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL não definida no ambiente")
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	// Testa a conexão pra garantir
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Erro ao pingar banco: %v", err)
	}

	DB = db
}
