package app

import (
	"database/sql"
	"fmt"
	"log"
	"minjust-website/internal/config"
	"minjust-website/internal/repository/postgres"
	v1 "minjust-website/internal/transport/http/v1"
	"minjust-website/internal/usecase/admin_auth"
	"minjust-website/internal/usecase/employee"
	"minjust-website/internal/usecase/password_request"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	db, err := connectDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	employeeRepo := postgres.NewPostgresEmployeeRepository(db)
	passwordRepo := postgres.NewPostgresPasswordRequestRepository(db)
	authRepo := postgres.NewPostgresAuthRepository(db)

	employeeUC := employee.NewEmployeeUsecase(employeeRepo)
	passwordUC := password_request.NewPasswordRequestUsecase(passwordRepo)
	authUC := admin_auth.NewAuthUsecase(authRepo, employeeRepo, cfg.JWTSecret)

	authHandler := v1.NewAuthHandler(authUC)
	reqHandler := v1.NewRequestHandler(passwordUC)
	empHandler := v1.NewEmployeeHandler(employeeUC)

	router := gin.Default()

	v1.SetupRoutes(router, cfg.JWTSecret, reqHandler, empHandler, authHandler)

	log.Printf("server started on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}

func connectDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	for attempt := 1; attempt <= 10; attempt++ {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		log.Printf("waiting for database... attempt %d/10", attempt)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database is not available")
}
