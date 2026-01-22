package utils

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

func ConnectRDS() *sql.DB {
	// Em produção, isso viria de variáveis de ambiente
	//connStr := "user=admin password=secret dbname=ecommerce_db sslmode=disable host=db port=5432"
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	
	// Setup simples da tabela (Migration simplificada)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email TEXT UNIQUE, password TEXT)`)
	return db
}