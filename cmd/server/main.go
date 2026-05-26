package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/marioalvaro/erp-system/internal/auth"
	"github.com/marioalvaro/erp-system/internal/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file (this is fine in production if vars are set in ECS)")
	}

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("POSTGRES_DB"),
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	userHandler := &users.UserHandler{DB: db}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ERP System is up and running!"))
	})

	mux.HandleFunc("POST /api/v1/auth/register", userHandler.RegisterUser)
	mux.HandleFunc("POST /api/v1/auth/login", userHandler.LoginUser)

	mux.Handle("GET /api/v1/users/profile", auth.AuthMiddleware(http.HandlerFunc(userHandler.GetProfile)))

	mux.Handle("GET /api/v1/admin/dashboard", auth.AuthMiddleware(auth.RequireRole("Administrator")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Admin Dashboard!"))
	}))))

	fmt.Println("Starting ERP Server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
