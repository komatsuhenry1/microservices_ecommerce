package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// --- Domain ---
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // A senha virá em texto plano no JSON
}

// --- Config & DB ---
// Simula a conexão com AWS RDS
func connectRDS() *sql.DB {
	// Em produção, isso viria de variáveis de ambiente
	connStr := "user=admin password=secret dbname=ecommerce_db sslmode=disable host=db port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	
	// Setup simples da tabela (Migration simplificada)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email TEXT UNIQUE, password TEXT)`)
	return db
}

// --- Handlers ---
type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Clean Code: Hash da senha (Nunca salve senha pura!)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	
	_, err := h.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, string(hashed))
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input User
	json.NewDecoder(r.Body).Decode(&input)

	// Busca usuário no banco
	var storedPassword string
	var userID int
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", input.Email).Scan(&userID, &storedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compara Hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Gera Token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	
	tokenString, _ := token.SignedString([]byte("MINHA_CHAVE_SECRETA"))

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// --- Main ---
func main() {
	db := connectRDS()
	handler := &AuthHandler{DB: db}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", handler.Register)
	mux.HandleFunc("POST /login", handler.Login)

	log.Println("Auth Service running on port 8081")
	http.ListenAndServe(":8081", mux)
}