package main

import (
	"context"
	"fmt"
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
	"go.uber.org/zap"
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
		utils.Log.Fatal("Failed to connect to database", zap.Error(err))
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
		utils.Log.Info("Auth Service starting", zap.String("port", cfg.AppPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Fatal("Listen failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		utils.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	utils.Log.Info("Server exiting")
}
