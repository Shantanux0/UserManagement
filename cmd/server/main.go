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
	logger.Init()
	defer func() {
		_ = logger.Log.Sync()
	}()

	logger.Log.Info("Starting User Management service...")

	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Unable to connect to database pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Log.Fatal("Failed to ping database pool", zap.Error(err))
	}
	logger.Log.Info("Successfully connected to PostgreSQL database")

	if err := db.RunMigrations(ctx, pool, logger.Log); err != nil {
		logger.Log.Fatal("Database migrations failed", zap.Error(err))
	}

	handler.InitValidator()

	userRepo := repository.NewPostgresUserRepository(pool)
	userService := service.NewDefaultUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger.Log)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logger.Log))

	routes.SetupRoutes(app, userHandler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Log.Info("Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			logger.Log.Error("Failed to shutdown server gracefully", zap.Error(err))
		}
	}()

	logger.Log.Info("Server listening", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Log.Fatal("Server failed to start or shut down with error", zap.Error(err))
	}
}
