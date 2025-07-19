package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/auth"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/melnik-dev/go_todo_jwt/pkg/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загрука конфигурации
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configs: %s", err)
	}

	// Инициализация логгера
	err = logger.InitLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err)
	}
	mainLogger := logger.GetLogger()

	mainLogger.WithFields(logrus.Fields{
		"app_name": cfg.App.Name,
		"env":      cfg.App.Env,
		"version":  cfg.App.Version,
	}).Info("Application starting")

	// Подключение к БД
	pgDB, err := db.InitPostgres(cfg, mainLogger)
	if err != nil {
		mainLogger.Fatalf("Error connecting to database: %s", err)
	}
	defer pgDB.Close()

	// Настройка Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	route := gin.New()

	// Middleware
	route.Use(gin.Recovery())
	route.Use(requestid.New())
	route.Use(middleware.Logger())

	route.GET("/ping", ping(mainLogger))

	// Repositories
	userRepo := user.NewRepository(pgDB, mainLogger)
	taskRepo := task.NewRepository(pgDB, mainLogger)

	// Services
	authService := auth.NewService(userRepo, mainLogger)
	taskService := task.NewService(taskRepo, mainLogger)

	// Handlers
	auth.NewHandler(route, &auth.HandlerDeps{
		AuthService: authService,
		Config:      cfg,
	})
	task.NewHandler(route, &task.HandlerDeps{
		TaskService: taskService,
		Config:      cfg,
	})

	// Настройка HTTP сервера
	serverAddress := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:           serverAddress,
		Handler:        route,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Запуск сервера
	go func() {
		mainLogger.WithField("address", serverAddress).Info("Server starting")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			mainLogger.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	mainLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		mainLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	mainLogger.Info("Server stopped gracefully")
}

func ping(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Ping endpoint called")

		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"timestamp": time.Now(),
		})
	}
}
