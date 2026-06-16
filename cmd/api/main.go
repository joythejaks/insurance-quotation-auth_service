package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jordisetiawan/insurance-auth-service/internal/config"
	"github.com/jordisetiawan/insurance-auth-service/internal/database"
	"github.com/jordisetiawan/insurance-auth-service/internal/handler"
	"github.com/jordisetiawan/insurance-auth-service/internal/repository"
	"github.com/jordisetiawan/insurance-auth-service/internal/router"
	"github.com/jordisetiawan/insurance-auth-service/internal/service"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"
)

// @title Insurance Auth Service API
// @version 1.0
// @description API Service for Authentication and Authorization.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	utils.InitLogger()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)                             // Inisialisasi RoleRepository
	authService := service.NewAuthService(userRepo, roleRepo, cfg.JWTSecret) // Lewatkan RoleRepository
	authHandler := handler.NewAuthHandler(authService)

	// Setup Gin Router
	r := router.SetupRouter(authHandler, cfg.JWTSecret)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: r,
	}

	// Graceful shutdown logic
	go func() {
		log.Printf("Auth Service starting on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
