package main

import (
	handlers "broker-backend/internal/handlers"
	mysqlRepo "broker-backend/internal/infra/repository/mysql"
	"broker-backend/internal/services"
	"broker-backend/pkg/auth"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrateMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		if strings.Contains(err.Error(), "Unknown database") {
			rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPass, dbHost, dbPort)
			tmp, terr := sqlx.Connect("mysql", rootDSN)
			if terr != nil {
				log.Fatalf("db connect root: %v", terr)
			}
			_, cerr := tmp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
			tmp.Close()
			if cerr != nil {
				log.Fatalf("create database: %v", cerr)
			}
			db, err = sqlx.Connect("mysql", dsn)
		}
	}
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := runMigrations(db.DB); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// Wiring dependencies.
	userRepo := mysqlRepo.NewUserRepository(db)
	holdingRepo := mysqlRepo.NewHoldingRepository(db)
	positionRepo := mysqlRepo.NewPositionRepository(db)
	orderRepo := mysqlRepo.NewOrderBookRepository(db)
	jwtManager := auth.NewJWTManager(secret)
	passwordHasher := auth.NewPasswordHasher(0)
	authSvc := services.NewAuthUsecase(userRepo, jwtManager, passwordHasher)
	holdingSvc := services.NewHoldingService(holdingRepo)
	positionSvc := services.NewPositionService(positionRepo)
	orderSvc := services.NewOrderBookService(orderRepo)

	// Handlers.
	authHandler := handlers.NewAuthHandler(authSvc)
	holdingHandler := handlers.NewHoldingHandler(holdingSvc)
	positionHandler := handlers.NewPositionHandler(positionSvc)
	orderHandler := handlers.NewOrderBookHandler(orderSvc)

	// Router.
	r := chi.NewRouter()

	r.Post("/signup", authHandler.SignUp)
	r.Post("/login", authHandler.Login)

	// Protected routes.
	jwtAuth := jwtManager.TokenAuth()
	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(jwtAuth))
		protected.Use(jwtauth.Authenticator)
		protected.Post("/holdings/create", holdingHandler.Create)
		protected.Get("/holdings", holdingHandler.List)
		protected.Post("/positions/create", positionHandler.Create)
		protected.Get("/positions", positionHandler.List)
		protected.Post("/orderbook/create", orderHandler.Create)
		protected.Get("/orderbook", orderHandler.List)
	})

	addr := ":8080"
	fmt.Printf("server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// runMigrations applies all up migrations using golang-migrate.
func runMigrations(db *sql.DB) error {
	driver, err := migrateMysql.WithInstance(db, &migrateMysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
