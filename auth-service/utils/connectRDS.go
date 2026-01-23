package utils

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

func ConnectRDS() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL não definido (env vazia)")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// força conexão agora (pra falhar cedo e com mensagem clara)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}