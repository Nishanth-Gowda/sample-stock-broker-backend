package main

import (
	"broker-backend/internal/infra/repository/mysql"
	httpHandler "broker-backend/internal/interface/http"
	"broker-backend/internal/usecase"
	"broker-backend/pkg/auth"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present.
	_ = godotenv.Load()

	// Gather DB env vars.
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "user"
	}
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "broker"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me"
	}

	// Connect to DB.
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	if err := ensureSchema(db.DB); err != nil {
		log.Fatalf("schema: %v", err)
	}

	// Wiring dependencies.
	userRepo := mysql.NewUserRepository(db)
	jwtManager := auth.NewJWTManager(secret)
	passwordHasher := auth.NewPasswordHasher(0)
	authUC := usecase.NewAuthUsecase(userRepo, jwtManager, passwordHasher)

	// Handlers.
	authHandler := httpHandler.NewAuthHandler(authUC)
	dataHandler := httpHandler.NewDataHandler()

	// Router.
	r := chi.NewRouter()

	r.Post("/signup", authHandler.SignUp)
	r.Post("/login", authHandler.Login)

	// Protected routes.
	jwtAuth := jwtManager.TokenAuth()
	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(jwtAuth))
		protected.Use(jwtauth.Authenticator)
		protected.Get("/holdings", dataHandler.Holdings)
		protected.Get("/orderbook", dataHandler.OrderBook)
		protected.Get("/positions", dataHandler.Positions)
	})

	addr := ":8080"
	fmt.Printf("server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// ensureSchema creates minimal users table when not exists.
func ensureSchema(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}
