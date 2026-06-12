package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-management/config"
	"user-management/db"
	"user-management/internal/handler"
	"user-management/internal/logger"
	"user-management/internal/middleware"
	"user-management/internal/repository"
	"user-management/internal/routes"
	"user-management/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Init()
	defer func() {
		_ = logger.Log.Sync()
	}()

	logger.Log.Info("Starting User Management service...")

	// Load configuration
	cfg := config.Load()

	// Establish connection pool to PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Unable to connect to database pool", zap.Error(err))
	}
	defer pool.Close()

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		logger.Log.Fatal("Failed to ping database pool", zap.Error(err))
	}
	logger.Log.Info("Successfully connected to PostgreSQL database")

	// Execute schema migrations
	if err := db.RunMigrations(ctx, pool, logger.Log); err != nil {
		logger.Log.Fatal("Database migrations failed", zap.Error(err))
	}

	// Initialize input validation and custom rules
	handler.InitValidator()

	// Setup layers
	userRepo := repository.NewPostgresUserRepository(pool)
	userService := service.NewDefaultUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger.Log)

	// Create Fiber Application
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
	})

	// Register Middleware
	app.Use(recover.New()) // Recover middleware recovers from panics anywhere in the stack
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logger.Log))

	// Register application routes
	routes.SetupRoutes(app, userHandler)

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Log.Info("Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			logger.Log.Error("Failed to shutdown server gracefully", zap.Error(err))
		}
	}()

	// Start listening
	logger.Log.Info("Server listening", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Log.Fatal("Server failed to start or shut down with error", zap.Error(err))
	}
}
